package wordrecognize

import (
	"EnglishProject/cors"
	"EnglishProject/db"
	"EnglishProject/jwt"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Word struct {
	WordID      int    `json:"wordID"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
}
type Progress struct {
	DailyGoal     int `json:"dailyGoal"`
	DailyProgress int `json:"dailyProgress"`
	DailyError    int `json:"dailyError"`
}

func GetRandomWordHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)

	// 查询数据库以获取随机单词
	var word, translation string
	var wordID int
	err := db.Db.QueryRow("SELECT id, word, translation FROM enwords ORDER BY RAND() LIMIT 1").Scan(&wordID, &word, &translation)
	if err != nil {
		http.Error(w, "Failed to get a random word from the database", http.StatusInternalServerError)
		log.Printf("Error getting a random word: %v", err)
		return
	}

	// 构造 Word 结构并将其转换为 JSON
	randomWord := Word{WordID: wordID, Word: word, Translation: translation}
	responseJSON, err := json.Marshal(randomWord)
	if err != nil {
		http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
		log.Printf("Error serializing JSON: %v", err)
		return
	}

	// 打印 JSON 数据到服务器日志
	log.Printf("JSON Response: %s", responseJSON)

	// 设置响应头并写入 JSON 数据
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
func GetWordIDFromBody(r *http.Request) (int, error) {
	// 读取请求体的内容
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	// 创建一个结构体来解码 JSON
	var requestBody struct {
		WordID int `json:"wordID"`
	}

	// 解码 JSON
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return 0, err
	}

	// 提取 wordID 参数
	wordID := requestBody.WordID

	return wordID, nil
}

func AddToEasywordHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)
	// 从令牌获取用户信息
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		// 未提供令牌，进行处理或返回未授权的响应
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 去掉 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 调用 ParseJWT 函数来解析和验证 JWT 令牌
	log.Printf("tokenstring: %s", tokenString)
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
	// 获取 wordID 参数
	wordID, err := GetWordIDFromBody(r)
	log.Printf("wordID: %d", wordID)
	if err != nil {
		http.Error(w, "Failed to get wordID from request body", http.StatusBadRequest)
		log.Printf("Failed to get wordID from request body %s", err)
		return
	}
	if wordID == 0 {
		http.Error(w, "wordID parameter is missing", http.StatusBadRequest)
		return
	}

	// 查询数据库以检查是否已存在记录
	existingUserID, existingWordID := 0, 0
	qerr := db.Db.QueryRow("SELECT user_id, word_id FROM easy_words WHERE user_id = ? AND word_id = ?", userID, wordID).Scan(&existingUserID, &existingWordID)

	if qerr != nil && qerr != sql.ErrNoRows {
		// 查询数据库时出错
		log.Println("Error querying database:", qerr)
		return
	}

	if existingUserID == 0 && existingWordID == 0 {
		// 不存在匹配的记录，执行插入操作
		_, insertErr := db.Db.Exec("INSERT INTO easy_words (user_id, word_id) VALUES (?, ?)", userID, wordID)
		if insertErr != nil {
			// 插入出错
			log.Println("Error inserting record:", insertErr)
		} else {
			// 插入成功
			log.Println("New record inserted")
		}
	} else {
		// 存在匹配的记录
		log.Printf("Record already exists%d %d", existingUserID, existingWordID)
	}

}
func AddToErrorWordHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)

	// 从令牌获取用户信息
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

	// 获取 wordID 参数
	wordID, err := GetWordIDFromBody(r)
	log.Printf("wordID: %d", wordID)
	if err != nil {
		http.Error(w, "Failed to get wordID from request body", http.StatusBadRequest)
		log.Printf("Failed to get wordID from request body %s", err)
		return
	}
	if wordID == 0 {
		http.Error(w, "wordID parameter is missing", http.StatusBadRequest)
		return
	}

	// 查询数据库以检查是否已存在记录
	existingUserID, existingWordID, existingErrorCount := 0, 0, 0
	qerr := db.Db.QueryRow("SELECT user_id, word_id, error_count FROM error_words WHERE user_id = ? AND word_id = ?", userID, wordID).Scan(&existingUserID, &existingWordID, &existingErrorCount)

	if qerr != nil && qerr != sql.ErrNoRows {
		// 查询数据库时出错
		log.Println("Error querying database:", qerr)
		return
	}

	if existingUserID == 0 && existingWordID == 0 {
		// 不存在匹配的记录，执行插入操作并设置 error_count 为 1
		_, insertErr := db.Db.Exec("INSERT INTO error_words (user_id, word_id, error_count) VALUES (?, ?, 1)", userID, wordID)
		if insertErr != nil {
			// 插入出错
			log.Println("Error inserting record:", insertErr)
		} else {
			// 插入成功
			log.Println("New record inserted")
		}
	} else {
		// 存在匹配的记录，增加 error_count 值
		_, updateErr := db.Db.Exec("UPDATE error_words SET error_count = error_count + 1 WHERE user_id = ? AND word_id = ?", userID, wordID)
		if updateErr != nil {
			log.Println("Error updating record:", updateErr)
		} else {
			log.Println("Record updated")
		}
	}
}
func ShowProgessHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)
	// 从令牌获取用户信息
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
	getUserDataFromDatabase(userID)
}
func getUserDataFromDatabase(userID int) (*Progress, error) {
	// 查询数据库以根据 userID 检索用户数据。
	query := "SELECT daily_goal, daily_progress, total_progress FROM user_data WHERE user_id = ?"
	row := db.Db.QueryRow(query, userID)

	// 创建 UserData 结构以存储检索到的数据。
	userProgress := &Progress{}
	err := row.Scan(
		&userProgress.DailyGoal,
		&userProgress.DailyProgress,
	)
	if err != nil {
		return nil, err
	}

	return userProgress, nil

}
