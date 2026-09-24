// GoShop 是「Go 實戰開發」課程的專案：一個迷你電商後台。
//
// 用法：
//
//	goshop -list                        列出商品
//	goshop -import products.csv         從 CSV 匯入商品
//	goshop -buy SKU-001:2 -coupon WEEK15  下單
//	goshop -http :8080                  啟動 API 與後台網頁
//
// 沒有指定 -dsn 時，資料存在 -data 資料夾的 gob 快照中；
// 指定 -dsn（或環境變數 GOSHOP_DSN）則改用 MySQL。
package main
