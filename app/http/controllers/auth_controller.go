package controllers

import (
	"encoding/json"
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/view"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 1、初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// 2、表单规则
	errs := requests.ValidateRegistrationForm(_user)

	// 3、执行
	if len(errs) > 0 {
		// 3.1 有错误发生，打印数据
		data, _ := json.MarshalIndent(errs, "", "  ")
		fmt.Fprint(w, string(data))
		// 3、表单不通过 —— 重新显示表单
		// view.RenderSimple(w, view.D{
		// 	"Errors": errs,
		// 	"User":   _user,
		// }, "auth.register")
	} else {
		// 4、验证成功，创建数据
		_user.Create()

		if _user.ID > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		}

	}
}
