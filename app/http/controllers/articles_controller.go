package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
	"strconv"
	"text/template"
	"unicode/utf8"

	"gorm.io/gorm"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// Show 文章详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {

	// 1、获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2、读取对应的文章数据
	article, err := article.Get(id)

	// 3、如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内存错误")
		}
	} else {
		// 4. 读取成功，显示文章
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, article)
	}

}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	// 1、获取结果集
	articles, err := article.GetAll()

	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 2、加载模板
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		// 3、渲染模板，将所有文章的数据传输进去
		err = tmpl.Execute(w, articles)
		logger.LogError(err)
	}
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func validatearticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	// err := r.ParseForm()
	// if err != nil {
	// 	// 解析错误，这里应该有错误处理
	// 	fmt.Fprint(w, "请提供正确的数据！")
	// 	return
	// }

	// title := r.PostForm.Get("title")

	// fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	// fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	// fmt.Fprintf(w, "title 的值为: %v", title)

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	// errors := make(map[string]string)

	// 验证表单提交数据
	errors := validatearticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		// fmt.Fprint(w, "验证通过!<br>")
		// fmt.Fprintf(w, "title 的值为：%v <br>", title)
		// fmt.Fprintf(w, "title 的长度为：%v <br>", utf8.RuneCountInString(title))
		// fmt.Fprintf(w, "body 的值为：%v <br>", body)
		// fmt.Fprintf(w, "body 的长度为：%v <br>", utf8.RuneCountInString(body))

		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功, ID 为"+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}

	} else {

		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}

	// fmt.Fprintf(w, "r.Form 中 title 的值为：%v <br>", r.FormValue("title"))
	// fmt.Fprintf(w, "r.PostForm 中 title 的值为：%v <br>", r.PostFormValue("title"))
	// fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
	// fmt.Fprintf(w, "r.PostForm 中 test 的值为：%v <br>", r.PostFormValue("test"))
}
