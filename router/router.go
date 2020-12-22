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
		utility.CheckACL(w, r, 0)
		controllers.BaseIndex(w, r)
	case "teams":
		controllers.Teams(w, r)
	case "manage":
		controllers.Manage(w, r)
	case "stats":
		controllers.Stats(w, r)
	case "test":
		utility.CheckACL(w, r, 2)
		controllers.Test(w, r)
	}
}
