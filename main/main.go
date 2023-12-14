package main

import (
	"EnglishProject/article"
	"EnglishProject/cors"
	"EnglishProject/db"
	"EnglishProject/login"
	"EnglishProject/register"
	"EnglishProject/userdata"
	"EnglishProject/wordbook"
	"EnglishProject/wordrecognize"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	db.Db, err = sql.Open("mysql", "root:Hello,world!@tcp(localhost:3306)/englishlearning")
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer db.Db.Close()

	// 测试数据库连接是否成功
	err = db.Db.Ping()
	if err != nil {
		fmt.Printf("Error pinging the database: %v\n", err)
		return
	}

	// 如果成功连接到数据库，继续进行其他操作
	fmt.Println("Connected to the database successfully")

	// 初始化HTTP路由
	http.HandleFunc("/register", register.RegisterHandler)                  // 注册路由
	http.HandleFunc("/login", login.LoginHandler)                           // 登录路由
	http.HandleFunc("/addToEasyWord", wordrecognize.AddToEasywordHandler)   // 简单词路由
	http.HandleFunc("/addToErrorWord", wordrecognize.AddToErrorWordHandler) // 错误词路由
	http.HandleFunc("/getRandomWord", wordrecognize.GetRandomWordHandler)   // 单词路由
	http.HandleFunc("/showWordList", wordbook.ShowWordListHandler)          // 词单路由
	http.HandleFunc("/deleteWord", wordbook.DeletWordFromListHandler)       //删除单词路由
	http.HandleFunc("/showUserData", userdata.ShowUserDataHandler)          //用户信息路由
	http.HandleFunc("/changePassword", userdata.ChangePasswordHandler)      //修改密码路由
	http.HandleFunc("/showArticles", article.GetArticles)                   //获取文章路由
	http.HandleFunc("/showArticleContent", article.GetArticleContent)       //获取文章内容路由
	// 启用CORS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cors.EnableCORS(w, r)
	})
	// 启动HTTP服务器
	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}
