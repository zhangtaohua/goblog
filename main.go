package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

// func defaultHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	if r.URL.Path == "/" {
// 		fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
// 	} else {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
// 			"<p>如有疑惑，请联系我们。</p>")
// 	}
// }

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHadnler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "创建新的文章")
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1、设置头
		w.Header().Set("Content-Type", "text/html; charset=uft-8")
		// 2、继续处理请求
		next.ServeHTTP(w, r)
	})
}

// func main() {
// 	http.HandleFunc("/", defaultHandler)
// 	http.HandleFunc("/about", aboutHandler)
// 	http.ListenAndServe(":3002", nil)
// }

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1 除首页以外， 移除所有请求路径后面的斜杆
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// router := http.NewServeMux()
	router := mux.NewRouter()
	// router := mux.NewRouter().StrictSlash(true) // cannot handle POST  not Use

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	// 自定义 404 page
	router.NotFoundHandler = http.HandlerFunc(notFoundHadnler)

	// middleWare
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL:", homeURL)

	articleURL, whatFuck := router.Get("articles.show").URL("id", "25")
	fmt.Println("articleURL", articleURL, whatFuck)

	// http.ListenAndServe(":3002", router)
	http.ListenAndServe(":3002", removeTrailingSlash(router))
}

// Content-Type
// text/html
// text/plain
// text/css
// text/javascript
// application/json
// application/xml
// image/png
// go clean =modcache
// go clean -modcache
// go mod download
// go mod download
// go mod init
// go mod download
// go mod tidy
// go mod graph
// go mod edit
// go mod vendor
// go mod verify
// go mod verify
// go mod why
// GO111MODULE
// GOSUMDB
// GONOPROXY
// GOPROXY GONOPROXY GONOSUMDB GOPRIVATE GOPRIVATE GOPRIVATE GOENV
