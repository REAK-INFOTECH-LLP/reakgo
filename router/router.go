package router

import (
	"net/http"
	"reakgo/controllers"
	"reakgo/utility"

	//"reakgo/utility"
	"strings"
)

func Routes(w http.ResponseWriter, r *http.Request) {
	// Trailing slash is a pain in the ass so we just drop it

	route := strings.Trim(r.URL.Path, "/")
	switch route {
	case "", "index":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.BaseIndex(w, r)
	case "login":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.Login(w, r)
	case "dashboard":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.Dashboard(w, r)
	case "addSimpleForm":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.AddForm(w, r)
	case "viewSimpleForm":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.ViewForm(w, r)
	case "register-2fa":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.RegisterTwoFa(w, r)
	case "verify-2fa":
		utility.CheckACL(w, r, []string{"admin", "user"})
		controllers.VerifyTwoFa(w, r)
	}
}
