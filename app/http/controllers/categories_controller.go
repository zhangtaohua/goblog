package controllers

import (
	"fmt"
	"net/http"

	"github.com/zhangtaohua/goblog/app/models/category"
	"github.com/zhangtaohua/goblog/app/requests"
	"github.com/zhangtaohua/goblog/pkg/flash"
	"github.com/zhangtaohua/goblog/pkg/route"
	"github.com/zhangtaohua/goblog/pkg/view"
)

type CategoriesController struct {
	BaseController
}

// Create 文章分类创建页面
func (*CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

// Store 保存文章分类
func (*CategoriesController) Store(w http.ResponseWriter, r *http.Request) {

	// 1. 初始化数据
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}

	// 2. 表单验证
	errors := requests.ValidateCategoryForm(_category)

	// 3. 检测错误
	if len(errors) == 0 {
		// 创建文章分类
		_category.Create()
		if _category.ID > 0 {
			// fmt.Fprint(w, "创建成功！")
			flash.Success("分类创建成功")
			indexURL := route.Name2URL("home")
			http.Redirect(w, r, indexURL, http.StatusFound)
			// indexURL := route.Name2URL("categories.show", "id", _category.GetStringID())
			// http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章分类失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Category": _category,
			"Errors":   errors,
		}, "categories.create")
	}
}
