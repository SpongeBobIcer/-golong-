package login

import (
	cors "EnglishProject/cors"
	"EnglishProject/db"
	"EnglishProject/jwt"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// 在用户登录时进行表单字段的验证
func validateLoginFields(email, password string) error {
	// 在这里添加各个字段的验证逻辑，例如：
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	// 添加更多验证逻辑，如验证邮箱格式等

	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)

	// 解析multipart/form-data请求
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 设置合适的表单大小限制
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		log.Printf("Error parsing form data: %v", err)
		return
	}

	// 获取用户提交的数据
	email := r.FormValue("email")
	password := r.FormValue("password")

	// 验证表单字段
	err = validateLoginFields(email, password)
	if err != nil {
		http.Error(w, "非法输入: "+err.Error(), http.StatusBadRequest)
		log.Printf("Invalid input data: %v", err)
		return
	}

	// 使用预备语句查询数据库验证用户信息
	stmt, err := db.Db.Prepare("SELECT password FROM users WHERE email = ?")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var storedPassword string
	err = stmt.QueryRow(email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			log.Printf("Error find user: %v", err)
			return
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// 检查邮箱和密码是否匹配
	if password != storedPassword {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	var username string
	var userID int
	err = db.Db.QueryRow("SELECT username, id FROM users WHERE email = ?", email).Scan(&username, &userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		log.Printf("Error getting user information: %v", err)
		return
	}

	// 如果验证通过，生成令牌并返回给客户端
	// 建议使用JWT令牌来进行用户认证
	var user jwt.User
	user.UserID = userID
	user.Username = username
	userToken, err := jwt.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	// 将令牌内容作为 JSON 响应返回给前端
	response := map[string]string{"token": userToken}
	json.NewEncoder(w).Encode(response)
	// 返回登录成功的响应
	log.Printf("Login successful")
}
