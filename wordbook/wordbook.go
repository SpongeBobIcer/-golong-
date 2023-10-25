package wordbook

import (
	"EnglishProject/cors"
	"EnglishProject/db"
	"EnglishProject/jwt"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Word struct {
	WordID      int    `json:"wordID"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	ErrorCount  int    `json:"errorCount"`
}

func ShowWordListHandler(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)

	// 从 URL 查询参数中获取要获取的内容类型
	content := r.URL.Query().Get("content")
	log.Printf("content: %s", content)
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

	// 根据查询参数 content 来决定要获取的内容
	if content == "simpleWords" {
		// 查询数据库以获取简单词列表
		rows, err := db.Db.Query("SELECT word_id FROM easy_words WHERE user_id = ?", userID)
		if err != nil {
			http.Error(w, "Failed to fetch easy list", http.StatusInternalServerError)
			log.Printf("Error fetching word list: %v", err)
			return
		}
		defer rows.Close()

		var wordList []Word

		for rows.Next() {
			var wordID int
			err := rows.Scan(&wordID)
			if err != nil {
				http.Error(w, "Failed to scan word list", http.StatusInternalServerError)
				log.Printf("Error scanning word list: %v", err)
				return
			}

			// 查询 enwords 表以获取单词和翻译
			var word, translation string
			err = db.Db.QueryRow("SELECT word, translation FROM enwords WHERE id = ?", wordID).Scan(&word, &translation)
			if err != nil {
				http.Error(w, "Failed to fetch word details", http.StatusInternalServerError)
				log.Printf("Error fetching word details: %v", err)
				return
			}

			wordList = append(wordList, Word{
				WordID:      wordID,
				Word:        word,
				Translation: translation,
			})
		}

		// 将单词列表转换为JSON并返回给前端
		responseJSON, err := json.Marshal(wordList)
		if err != nil {
			http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
			log.Printf("Error serializing JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
		// ... (使用相同的代码来查询和返回简单词列表)
	} else if content == "errorWords" {
		// 查询数据库以获取错误词列表
		// 查询数据库以获取用户的单词列表
		rows, err := db.Db.Query("SELECT word_id FROM error_words WHERE user_id = ?", userID)
		if err != nil {
			http.Error(w, "Failed to fetch word list", http.StatusInternalServerError)
			log.Printf("Error fetching word list: %v", err)
			return
		}
		defer rows.Close()

		var wordList []Word

		for rows.Next() {
			var wordID int
			err := rows.Scan(&wordID)
			if err != nil {
				http.Error(w, "Failed to scan word list", http.StatusInternalServerError)
				log.Printf("Error scanning word list: %v", err)
				return
			}

			// 查询 enwords 表以获取单词和翻译
			var word, translation string
			err = db.Db.QueryRow("SELECT word, translation FROM enwords WHERE id = ?", wordID).Scan(&word, &translation)
			if err != nil {
				http.Error(w, "Failed to fetch word details", http.StatusInternalServerError)
				log.Printf("Error fetching word details: %v", err)
				return
			}

			wordList = append(wordList, Word{
				WordID:      wordID,
				Word:        word,
				Translation: translation,
			})
		}

		// 将单词列表转换为JSON并返回给前端
		responseJSON, err := json.Marshal(wordList)
		if err != nil {
			http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
			log.Printf("Error serializing JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else {
		// 未知的内容类型
		http.Error(w, "Unknown content type", http.StatusBadRequest)
	}
}
func DeletWordFromListHandler(w http.ResponseWriter, r *http.Request) {
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

	// 从 URL 查询参数中获取要删除的单词ID
	wordIDParam := r.URL.Query().Get("wordID")
	if wordIDParam == "" {
		http.Error(w, "Missing wordID parameter", http.StatusBadRequest)
		return
	}

	// 将 wordIDParam 转换为整数
	wordID, err := strconv.Atoi(wordIDParam)
	if err != nil {
		http.Error(w, "Invalid wordID parameter", http.StatusBadRequest)
		return
	}

	// 根据查询参数 content 来决定要删除的表
	content := r.URL.Query().Get("content")
	if content == "simpleWords" {
		// 删除简单词
		_, err = db.Db.Exec("DELETE FROM easy_words WHERE word_id = ? AND user_id = ?", wordID, userID)
	} else if content == "errorWords" {
		// 删除错误词
		_, err = db.Db.Exec("DELETE FROM error_words WHERE word_id = ? AND user_id = ?", wordID, userID)
	} else {
		http.Error(w, "Unknown content type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to delete word", http.StatusInternalServerError)
		log.Printf("Error deleting word: %v", err)
		return
	}

	// 返回成功的响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Word deleted successfully"))
}
