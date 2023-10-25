package cors

import "net/http"

// 启用CORS处理
func EnableCORS(w http.ResponseWriter, r *http.Request) {
	// 设置允许的源
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")

	// 设置允许的HTTP方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// 设置允许的自定义标头
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // 添加 "Authorization" 标头

	// 允许使用凭证（例如Cookie）
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// 如果请求方法是OPTIONS，则表示预检请求，直接返回200状态码
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
