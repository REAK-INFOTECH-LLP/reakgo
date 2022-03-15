package router

import (
	"net/http"
	"reakgo/controllers"
	//"reakgo/utility"
	"strings"
)

func Routes(w http.ResponseWriter, r *http.Request) {

	// Trailing slash is a pain in the ass so we just drop it
	route := strings.Trim(r.URL.Path, "/")
	switch route {
	case "", "index":
		controllers.BaseIndex(w, r)
	case "login":
		controllers.Login(w, r)
	case "dashboard":
		controllers.Dashboard(w, r)
	case "addForm":
		controllers.AddForm(w, r)
	case "viewForm":
		controllers.ViewForm(w, r)
	}
}
