package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"goshop/internal/checkout"
	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

// options 是命令列旗標解析後的設定
type options struct {
	dir       string // 資料夾位置
	dsn       string // MySQL 連線字串；空白代表使用記憶體＋gob 快照
	list      bool   // 列出商品
	importCSV string // 要匯入的商品 CSV
	exportCSV string // 要匯出的訂單 CSV
	buy       string // 下單內容
	coupon    string // 折價券代碼
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
	flag.Parse()

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

	var err error
	switch {
	case opt.importCSV != "":
		err = importProducts(ctx, st, opt.importCSV)
	case opt.exportCSV != "":
		err = exportOrders(ctx, st, opt.exportCSV)
	case opt.buy != "":
		err = placeOrder(ctx, st, opt.dir, opt.buy, opt.coupon)
	case opt.list:
		err = listProducts(ctx, st)
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

func listProducts(ctx context.Context, st store.Store) error {
	products, err := st.Products(ctx)
	if err != nil {
		return err
	}
	for _, p := range products {
		fmt.Printf("%s  %-10v 庫存 %3d  %s\n", p.SKU, p.Price, p.Stock, p.Name)
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

func placeOrder(ctx context.Context, st store.Store, dir, buy, coupon string) error {
	items, err := parseItems(buy)
	if err != nil {
		return err
	}
	now := time.Now()
	svc := &checkout.Service{
		Store: st,
		Rules: []checkout.Discount{
			checkout.PercentOff(10), checkout.Threshold(2000, 300)},
		Coupons: map[string]checkout.Coupon{
			"WEEK15": {Code: "WEEK15", PercentOff: 15,
				Start: now.AddDate(0, 0, -7), End: now.AddDate(0, 0, 7)},
		},
	}
	o, err := svc.Checkout(ctx, checkout.Cart{Items: items, Coupon: coupon})
	if err != nil {
		return err
	}
	if err := svc.Pay(ctx, o.ID, payment.CashOnDelivery{}); err != nil {
		return err
	}
	o, _ = st.Order(ctx, o.ID)
	fmt.Printf("訂單 #%d 成立：%v（%s），預計 %s 出貨\n", o.ID, o.Total,
		o.PaidBy, o.ShipBy.Format(time.DateOnly))
	return appendLog(dir, o)
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
