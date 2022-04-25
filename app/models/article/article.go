package article

import (
	"strconv"

	"github.com/zhangtaohua/goblog/app/models"
	"github.com/zhangtaohua/goblog/app/models/user"
	"github.com/zhangtaohua/goblog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title    string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body     string `gorm:"type:longtext;not null;" valid:"body"`
	UserID   uint64 `gorm:"not null;index"`
	User     user.User
	Category uint64 `gorm:"not null;default:4; index"`
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

// CreatedAtDate 创建日期
func (article Article) CreatedAtDate() string {
	return article.CreatedAt.Format("2006-01-02")
}
