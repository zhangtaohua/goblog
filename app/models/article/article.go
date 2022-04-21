package article

import (
	"strconv"

	"github.com/zhangtaohua/goblog/app/models"
	"github.com/zhangtaohua/goblog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title string
	Body  string
}

// 这是一个 Object 的方法
// func (article Article) Link() string {
// 	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
// 	if err != nil {
// 		logger.LogError(err)
// 		return ""
// 	}
// 	return showURL.String()
// }
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(article.ID, 10))
}
