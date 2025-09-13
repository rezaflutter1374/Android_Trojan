package websocket

import "github.com/gorilla/websocket"

// Client مدلی برای هر دستگاه آلوده (bot) است که از طریق WebSocket متصل می‌شود.
// این ساختار اطلاعات مربوط به شناسه و اتصال فعال آن را نگهداری می‌کند.
type Client struct {
	ID   string          // یک شناسه منحصر به فرد برای کلاینت (مثلاً UUID یا شناسه دستگاه)
	Conn *websocket.Conn // اشاره‌گر به اتصال WebSocket فعال
}