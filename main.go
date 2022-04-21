package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zhangtaohua/goblog/bootstrap"
	"github.com/zhangtaohua/goblog/pkg/database"
	"github.com/zhangtaohua/goblog/pkg/logger"
)

var router *mux.Router
var db *sql.DB

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

type Article struct {
	Title, Body string
	ID          int64
}

// Delete 方法用以从数据库中删除单条记录
func (a Article) Delete() (rowsAffected int64, err error) {
	// Exec() 一般用在 CREATE/UPDATE/DELETE方法中
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	// √ 删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		fmt.Println("SQL 调用成功 ！！ id" + strconv.FormatInt(a.ID, 10))
		return n, nil
	}

	return 0, nil
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 执行查询语句，返回一个结果集
	// Query 是为了从数据库中读取多条数据
	// QueryRow() 是读取单条的数据
	// 单一参数 的纯文体模式， 和多个参数的 Prepare 模式
	// 纯文本模式只地发送一次SQL的请求， 而Prepare 模式 会发送两次
	// Rows 对象，是Query()返回的结果集， 包含数据 和 SQL 连接
	//
	rows, err := db.Query("SELECT * from articles")
	logger.LogError(err)
	// 需要在检测err 以后再调用defer，否则会让运行时 panic
	// 一个建议是如果 在循环中执行了Query() 并获取了Rows 结果集，不要在里面使用defer ,而是直接调用 rows.Close()
	defer rows.Close()

	var articles []Article
	//2. 循环读取结果
	for rows.Next() {
		var article Article
		// 2.1 扫描每一行的结果并赋值到一个 article 对象中
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		logger.LogError(err)
		// 2.2 将 article 追加到 articles 的这个数组中
		articles = append(articles, article)
	}

	// 2.3 检测遍历时是否发生错误
	err = rows.Err()
	logger.LogError(err)

	// 3. 加载模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)

	// 4. 渲染模板，将所有文章的数据传输进去
	err = tmpl.Execute(w, articles)
	logger.LogError(err)
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1、设置头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 SQL 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				fmt.Println("删除成功后开始调转！")
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()

	router = bootstrap.SetupRoute()
	// router := http.NewServeMux()
	// router := mux.NewRouter()
	// router := mux.NewRouter().StrictSlash(true) // cannot handle POST  not Use

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

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
