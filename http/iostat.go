package http

import (
	"net/http"
)

func configIoStatRoutes() {
	http.HandleFunc("/page/diskio", func(w http.ResponseWriter, r *http.Request) {
		//RenderDataJson(w, funcs.IOStatsForPage())
	})
}
