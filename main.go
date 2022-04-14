package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
			fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
	} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
					"<p>如有疑惑，请联系我们。</p>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
	"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

// func main() {
// 	http.HandleFunc("/", defaultHandler)
// 	http.HandleFunc("/about", aboutHandler)
// 	http.ListenAndServe(":3002", nil)
// }

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)

	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		switch r.Method {
		case "GET": 
			fmt.Fprint(w, "GET\n")
		case "POST": 
			fmt.Fprint(w, "POST\n")
		}
		fmt.Fprint(w, "article ID:" +  id)	
	})

	http.ListenAndServe(":3002", router)
}

// Content-Type
// text/html
// text/plain
// text/css
// text/javascript
// application/json
// application/xml
// image/png