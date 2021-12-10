package main

import (
	"goblog/bootstrap"
	"goblog/pkg/database"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router *mux.Router

// Article 对应一条文章数据
type Article struct {
	Title, Body string
	ID          int64
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1、设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2、继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1、除首页外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 2、将请求传递下去
		next.ServeHTTP(w, r)
	})
}

func main() {

	database.Initialize()

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("homeURL：", homeURL)
	// articleURL, _ := router.Get("articles.show").URL("id", "23")
	// fmt.Println("articleURL：", articleURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
