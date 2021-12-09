package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

// Name2URL 通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {

	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		// checkError(err)
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取 URL 路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
