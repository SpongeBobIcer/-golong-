package article

import (
	"EnglishProject/cors"
	"EnglishProject/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Article 结构体表示文章的数据模型
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func fetchArticles(offset int, limit int) ([]Article, error) {
	rows, err := db.Db.Query("SELECT id, title FROM reading_articles LIMIT ?, ?", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []Article{}
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title)
		if err != nil {
			return nil, err
		}
		log.Println(article.ID, article.Title)
		articles = append(articles, article)
	}

	return articles, nil
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	cors.EnableCORS(w, r)
	pageStr := r.FormValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}
	offset := (page - 1) * 5

	articles, err := fetchArticles(offset, 5)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 返回 JSON 格式的文章列表
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}
func fetchArticleContent(articleID int) (*Article, error) {

	// 查询文章内容
	row := db.Db.QueryRow("SELECT id, title, article_content FROM reading_articles WHERE id = ?", articleID)

	// 将查询结果赋值给 Article 结构体
	var article Article
	err := row.Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// GetArticleContent 函数处理HTTP请求，获取文章内容，并将其以JSON格式返回给前端
func GetArticleContent(w http.ResponseWriter, r *http.Request) {
	// 启用CORS
	cors.EnableCORS(w, r)

	// 从URL参数中获取文章ID
	articleIDStr := r.FormValue("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 获取文章内容
	article, err := fetchArticleContent(articleID)
	if err != nil {
		log.Println("Error fetching article content:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 返回 JSON 格式的文章内容
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}
