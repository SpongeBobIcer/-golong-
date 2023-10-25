package userdata

import (
	"EnglishProject/cors"
	"EnglishProject/db"
	"EnglishProject/jwt"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type User struct {
	UserID      int    // 用户ID
	Username    string // 用户ID
	Email       string // 邮箱
	PhoneNumber string // 手机号码
}

func ShowUserDataHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		// 未提供令牌，进行处理或返回未授权的响应
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 去掉 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 调用 ParseJWT 函数来解析和验证 JWT 令牌
	claims, err := jwt.ParseJWT(tokenString)
	if err != nil {
		// 处理令牌无效或解析错误的情况
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}

	// 从 claims 中提取用户信息或其他数据
	userIDFloat64 := claims["userID"].(float64)
	userID := int(userIDFloat64)
	log.Printf("userID: %d", userID)
	// 从数据库或其他存储位置检索用户信息
	user, err := getUserInfoFromDatabase(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		log.Printf("Error retrieving user information: %v", err)
		return
	}

	// 将用户信息编码为JSON并发送给前端
	responseJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
		log.Printf("Error serializing JSON: %v", err)
		return
	}
	//log.Printf("json:%s", responseJSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
func ReadPasswordFromRequest(r *http.Request) (string, string, error) {
	var requestData struct {
		OldPassword        string `json:"oldPassword"`
		NewPassword        string `json:"newPassword"`
		ConfirmNewPassword string `json:"confirmNewPassword"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		return "", "", err
	}

	return requestData.OldPassword, requestData.NewPassword, nil
}

func getUserInfoFromDatabase(userID int) (*User, error) {
	// 在这里实现从数据库中检索用户信息的逻辑
	query := "SELECT id,username, email, phone FROM users WHERE id = ?"
	var username, email, phone string
	var userid int
	err := db.Db.QueryRow(query, userID).Scan(&userid, &username, &email, &phone)
	if err != nil {
		return nil, err
	}

	// 假设从数据库中获取用户信息
	user := &User{
		UserID:      userid,
		Username:    username,
		Email:       email,
		PhoneNumber: phone,
	}

	return user, nil
}
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		// 未提供令牌，进行处理或返回未授权的响应
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 去掉 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 调用 ParseJWT 函数来解析和验证 JWT 令牌
	claims, err := jwt.ParseJWT(tokenString)
	if err != nil {
		// 处理令牌无效或解析错误的情况
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}

	// 从 claims 中提取用户信息或其他数据
	userIDFloat64 := claims["userID"].(float64)
	userID := int(userIDFloat64)
	log.Printf("userID: %d", userID)
	// 获取原密码和新密码
	oldPassword, newPassword, err := ReadPasswordFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to read passwords from the request", http.StatusBadRequest)
		return
	}

	// 查询数据库以获取用户的原密码
	var dbOldPassword string
	qerr := db.Db.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&dbOldPassword)
	if qerr != nil {
		http.Error(w, "Failed to fetch the user's password", http.StatusInternalServerError)
		return
	}

	// 比较提供的原密码和数据库中的原密码
	if oldPassword != dbOldPassword {
		http.Error(w, "Original password does not match", http.StatusUnauthorized)
		log.Printf("输入的 password: %s,数据库的 password: %s", oldPassword, dbOldPassword)
		return
	}

	// 如果原密码匹配，更新数据库中的密码
	_, updateErr := db.Db.Exec("UPDATE users SET password = ? WHERE id = ?", newPassword, userID)
	if updateErr != nil {
		http.Error(w, "Failed to update the password", http.StatusInternalServerError)
		return
	}

	// 密码已成功更改
	// 可以返回成功响应或其他操作
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password has been changed successfully"))
}
