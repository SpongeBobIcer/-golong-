package register

import (
	"EnglishProject/cors"
	"EnglishProject/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func checkIfUserExists(username, email, phone string) (string, error) {
	var existsField string

	// 查询数据库检查哪个字段已经存在
	err := db.Db.QueryRow("SELECT CASE WHEN username = ? THEN 'username' "+
		"WHEN email = ? THEN 'email' "+
		"WHEN phone = ? THEN 'phone' ELSE 'not_found' END "+
		"FROM users WHERE username = ? OR email = ? OR phone = ?", username, email, phone, username, email, phone).Scan(&existsField)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checkIfUserExists: %v", err)
		return "", err
	}

	return existsField, nil
}

// 在用户注册时进行表单字段的验证
func validateRegistrationFields(username, email, password, phone string) error {
	// 在这里添加各个字段的验证逻辑，例如：
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	// 添加更多验证逻辑，如验证邮箱格式、密码强度、电话号码格式等

	return nil
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)

	// 解析multipart/form-data请求
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 设置合适的表单大小限制
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		log.Printf("Error parsing form data: %v", err)
		return
	}

	// 获取表单字段值
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	phone := r.FormValue("phone")

	// 验证表单字段
	err = validateRegistrationFields(username, email, password, phone)
	if err != nil {
		http.Error(w, "非法输入: "+err.Error(), http.StatusBadRequest)
		log.Printf("Invalid input data: %v", err)
		return
	}

	// 检查用户名、邮箱和电话号码是否已经存在
	existsField, err := checkIfUserExists(username, email, phone)
	if err != nil {
		http.Error(w, "注册时发生数据库错误", http.StatusInternalServerError)
		return
	}
	if existsField != "" {
		fmt.Fprintf(w, "注册失败: %s 已经被注册了", existsField)
		return
	}

	// 在数据库中插入用户信息
	// 这里你可以添加密码哈希等更多安全性措施
	_, err = db.Db.Exec("INSERT INTO users (username, email, password, phone) VALUES (?, ?, ?, ?)", username, email, password, phone)
	if err != nil {
		http.Error(w, "Failed to insert user data into the database: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting user data into the database: %v", err)
		return
	}
	// 返回注册成功的响应

	fmt.Fprintln(w, "注册成功")
}
