package policies

import (
	"github.com/zhangtaohua/goblog/app/models/article"
	"github.com/zhangtaohua/goblog/pkg/auth"
)

// CanModifyArticle 是否允许修改话题
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.UserID
}
