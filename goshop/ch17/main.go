package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"goshop/internal/checkout"
	"goshop/internal/payment"
	"goshop/internal/rates"
	"goshop/internal/shop"
	"goshop/internal/store"
	"goshop/internal/web"
	"goshop/internal/webhook"
)

// version 在編譯時用 -ldflags "-X main.version=v1.0.0" 注入
var version = "dev"

// options 是命令列旗標解析後的設定
type options struct {
	dir       string // 資料夾位置
	dsn       string // MySQL 連線字串；空白代表使用記憶體＋gob 快照
	list      bool   // 列出商品
	importCSV string // 要匯入的商品 CSV
	exportCSV string // 要匯出的訂單 CSV
	buy       string // 下單內容
	coupon    string // 折價券代碼
	currency  string // 列出商品時換算的外幣，例如 USD
	webhook   string // 訂單付款後要通知的網址
	addr      string // HTTP 伺服器的位址，例如 :8080
	version   bool   // 印出版本
}

func main() {
	var opt options
	flag.StringVar(&opt.dir, "data", "data", "資料夾位置")
	flag.StringVar(&opt.dsn, "dsn", os.Getenv("GOSHOP_DSN"),
		"MySQL 連線字串（預設讀取環境變數 GOSHOP_DSN）")
	flag.BoolVar(&opt.list, "list", false, "列出所有商品")
	flag.StringVar(&opt.importCSV, "import", "", "從 CSV 檔匯入或更新商品")
	flag.StringVar(&opt.exportCSV, "export", "", "把訂單匯出成 CSV 檔")
	flag.StringVar(&opt.buy, "buy", "", "下單，例如 SKU-001:2,SKU-003:1")
	flag.StringVar(&opt.coupon, "coupon", "", "折價券代碼")
	flag.StringVar(&opt.currency, "currency", "", "列出商品時換算成外幣，例如 USD、JPY")
	flag.StringVar(&opt.webhook, "webhook", os.Getenv("GOSHOP_WEBHOOK"),
		"訂單付款後要通知的網址")
	flag.StringVar(&opt.addr, "http", "", "啟動 HTTP 伺服器，例如 :8080")
	flag.BoolVar(&opt.version, "version", false, "印出版本後結束")
	flag.Parse()

	if opt.version {
		fmt.Println("GoShop", version)
		return
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err := run(opt); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

func run(opt options) error {
	ctx := context.Background()
	// 資料夾只有自己可以讀寫（rwx------）
	if err := os.MkdirAll(opt.dir, 0o700); err != nil {
		return err
	}

	var st store.Store
	if opt.dsn != "" {
		db, err := store.OpenMySQL(ctx, opt.dsn)
		if err != nil {
			return err
		}
		defer db.Close()
		st = db
	} else {
		mem, err := openStore(opt.dir)
		if err != nil {
			return err
		}
		st = mem
	}

	svc := newService(st)
	if opt.webhook != "" {
		n := webhook.NewNotifier(opt.webhook, 4)
		defer n.Close() // 結束前等通知送完
		svc.OnPaid = func(o shop.Order) {
			n.Notify(webhook.Event{Type: "order.paid", Data: o})
		}
	}

	var err error
	switch {
	case opt.addr != "":
		err = serve(svc, opt.addr)
	case opt.importCSV != "":
		err = importProducts(ctx, st, opt.importCSV)
	case opt.exportCSV != "":
		err = exportOrders(ctx, st, opt.exportCSV)
	case opt.buy != "":
		err = placeOrder(ctx, svc, opt)
	case opt.list:
		err = listProducts(ctx, st, opt.currency)
	default:
		flag.Usage()
		return nil
	}
	if err != nil {
		return err
	}
	// 只有記憶體版需要存快照；用型別斷言判斷實際的型別
	if mem, ok := st.(*store.Memory); ok {
		return saveStore(mem, opt.dir)
	}
	return nil
}

func listProducts(ctx context.Context, st store.Store, currency string) error {
	products, err := st.Products(ctx)
	if err != nil {
		return err
	}
	rate := 0.0
	if currency != "" {
		rate, err = rates.New(os.Getenv("GOSHOP_RATES_URL")).Rate(ctx, "TWD", currency)
		if err != nil {
			fmt.Println("無法換算外幣：", err) // 查不到匯率不影響列出商品
		}
	}
	for _, p := range products {
		fmt.Printf("%s  %-10v 庫存 %3d  %s", p.SKU, p.Price, p.Stock, p.Name)
		if rate > 0 {
			fmt.Printf("（約 %.2f %s）", float64(p.Price)*rate, currency)
		}
		fmt.Println()
	}
	return nil
}

func importProducts(ctx context.Context, st store.Store, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	products, err := store.ReadProductsCSV(f)
	for _, p := range products {
		if err := st.SaveProduct(ctx, p); err != nil {
			return err
		}
	}
	fmt.Printf("已匯入 %d 項商品\n", len(products))
	if err != nil {
		fmt.Println("以下資料列被略過：\n", err)
	}
	return nil
}

func exportOrders(ctx context.Context, st store.Store, name string) error {
	orders, err := st.Orders(ctx)
	if err != nil {
		return err
	}
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := store.WriteOrdersCSV(f, orders); err != nil {
		return err
	}
	fmt.Printf("已匯出 %d 張訂單到 %s\n", len(orders), name)
	return nil
}

// newService 建立結帳服務：全館 9 折、滿 2000 折 300，加上本週折價券
func newService(st store.Store) *checkout.Service {
	now := time.Now()
	return &checkout.Service{
		Store: st,
		Rules: []checkout.Discount{
			checkout.PercentOff(10), checkout.Threshold(2000, 300)},
		Coupons: map[string]checkout.Coupon{
			"WEEK15": {Code: "WEEK15", PercentOff: 15,
				Start: now.AddDate(0, 0, -7), End: now.AddDate(0, 0, 7)},
		},
	}
}

// serve 啟動 HTTP 伺服器，收到 Ctrl+C 或 SIGTERM 時優雅關閉
func serve(svc *checkout.Service, addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := &web.Server{Store: svc.Store, Checkout: svc}
	srv := &http.Server{
		Addr:              addr,
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	slog.Info("伺服器啟動", "addr", addr, "version", version)

	select {
	case err := <-errCh: // 例如連接埠已被佔用
		return err
	case <-ctx.Done():
	}
	slog.Info("收到結束訊號，關閉伺服器中…")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx) // 等進行中的請求處理完才關閉
}

func placeOrder(ctx context.Context, svc *checkout.Service, opt options) error {
	items, err := parseItems(opt.buy)
	if err != nil {
		return err
	}
	o, err := svc.Checkout(ctx, checkout.Cart{Items: items, Coupon: opt.coupon})
	if err != nil {
		return err
	}
	if err := svc.Pay(ctx, o.ID, payment.CashOnDelivery{}); err != nil {
		return err
	}
	o, _ = svc.Store.Order(ctx, o.ID)
	fmt.Printf("訂單 #%d 成立：%v（%s），預計 %s 出貨\n", o.ID, o.Total,
		o.PaidBy, o.ShipBy.Format(time.DateOnly))
	return appendLog(opt.dir, o)
}

// parseItems 把 "SKU-001:2,SKU-003:1" 解析成購物車品項
func parseItems(s string) ([]checkout.Item, error) {
	var items []checkout.Item
	for part := range strings.SplitSeq(s, ",") {
		sku, qtyText, ok := strings.Cut(part, ":")
		qty, err := strconv.Atoi(qtyText)
		if !ok || err != nil || qty <= 0 {
			return nil, fmt.Errorf("品項格式錯誤：%q（應為 SKU:數量）", part)
		}
		items = append(items, checkout.Item{SKU: sku, Qty: qty})
	}
	return items, nil
}

// appendLog 把訂單以 JSON Lines 格式附加到 orders.log 的最後面
func appendLog(dir string, o shop.Order) error {
	f, err := os.OpenFile(filepath.Join(dir, "orders.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(o) // Encode 會在結尾加上換行
}

// openStore 優先讀取上次存下的 gob 快照；
// 第一次執行時沒有快照，就從 products.json 匯入商品。
func openStore(dir string) (*store.Memory, error) {
	st := store.NewMemory()
	f, err := os.Open(filepath.Join(dir, "goshop.gob"))
	if err == nil {
		defer f.Close()
		return st, st.Load(f)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	f, err = os.Open(filepath.Join(dir, "products.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	products, err := store.DecodeProducts(f)
	if err != nil {
		return nil, err
	}
	return store.NewMemory(products...), nil
}

// saveStore 把資料存成 gob 快照；檔案權限 0600 表示只有自己能讀寫。
func saveStore(st *store.Memory, dir string) error {
	f, err := os.OpenFile(filepath.Join(dir, "goshop.gob"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return st.Save(f)
}
