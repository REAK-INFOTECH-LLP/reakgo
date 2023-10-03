package router

import (
	"net/http"
	"reakgo/controllers"
	"strings"
)

func Routes(w http.ResponseWriter, r *http.Request) {

	// Trailing slash is a pain in the ass so we just drop it
	route := strings.Trim(r.URL.Path, "/")
	switch route {
	case "", "index":
		controllers.CheckACL(w, r, []string{"guest", "admin", "user"})
		controllers.BaseAPI(w, r)
	case "login":
		controllers.CheckACL(w, r, []string{"admin", "user"})
		controllers.Login(w, r)
	case "dashboard":
		controllers.CheckACL(w, r, []string{"admin", "user"})
		controllers.Dashboard(w, r)
	case "addSimpleForm":
		controllers.CheckACL(w, r, []string{"admin", "user"})
		controllers.AddForm(w, r)
	case "viewSimpleForm":
		controllers.CheckACL(w, r, []string{"admin", "user"})
		controllers.ViewForm(w, r)
	case "register-2fa":
		controllers.CheckACL(w, r, []string{"admin", "user"})
		controllers.RegisterTwoFa(w, r)
	case "verify-2fa":
		controllers.CheckACL(w, r, []string{"admin", "user"})
		controllers.VerifyTwoFa(w, r)
	}
}
