package router

import (
	"net/http"
	"reakgo/controllers"
	"reakgo/utility"
	"strings"
)

func Routes(w http.ResponseWriter, r *http.Request) {

	// Trailing slash is a pain in the ass so we just drop it
	route := strings.Trim(r.URL.Path, "/")
	switch route {
	case "", "index":
		//utility.CheckACL(w, r, 0)
		controllers.BaseAPI(w, r)
	case "login":
		utility.CheckACL(w, r, 0)
		controllers.Login(w, r)
	case "dashboard":
		utility.CheckACL(w, r, 1)
		controllers.Dashboard(w, r)
	case "addSimpleForm":
		utility.CheckACL(w, r, 0)
		controllers.AddForm(w, r)
	case "viewSimpleForm":
		utility.CheckACL(w, r, 0)
		controllers.ViewForm(w, r)
	case "register-2fa":
		utility.CheckACL(w, r, 1)
		controllers.RegisterTwoFa(w, r)
	case "verify-2fa":
		utility.CheckACL(w, r, 1)
		controllers.VerifyTwoFa(w, r)
	}
}
