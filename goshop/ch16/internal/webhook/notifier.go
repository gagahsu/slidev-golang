package webhook

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Notifier 用背景的 worker pool 送出通知，呼叫端不必等待對方回應。
type Notifier struct {
	url    string
	client *http.Client
	jobs   chan Event
	wg     sync.WaitGroup
}

// NewNotifier 啟動 workers 個 goroutine，從佇列取出事件送到 url。
func NewNotifier(url string, workers int) *Notifier {
	n := &Notifier{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
		jobs:   make(chan Event, 100), // 有緩衝的通道就是一個佇列
	}
	for id := range workers {
		n.wg.Go(func() { n.worker(id) })
	}
	return n
}

func (n *Notifier) worker(id int) {
	for e := range n.jobs { // 通道被關閉、而且取完之後，迴圈才會結束
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := Send(ctx, n.client, n.url, e)
		cancel()
		if err != nil {
			slog.Warn("通知失敗", "worker", id, "err", err)
			continue
		}
		slog.Debug("通知成功", "worker", id, "type", e.Type)
	}
}

// Notify 把事件放進佇列；佇列滿了就放棄，絕不讓呼叫端卡住。
func (n *Notifier) Notify(e Event) bool {
	select {
	case n.jobs <- e:
		return true
	default:
		slog.Warn("通知佇列已滿，放棄這次通知", "type", e.Type)
		return false
	}
}

// Close 停止接收新事件，並等待佇列中剩下的事件都送完。
func (n *Notifier) Close() {
	close(n.jobs)
	n.wg.Wait()
}
