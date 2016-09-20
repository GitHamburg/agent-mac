package http

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
	"../tools/host"
	"../tools/load"
)

func configSystemRoutes() {

	http.HandleFunc("/system/date", func(w http.ResponseWriter, req *http.Request) {
		RenderDataJson(w, time.Now().Format("2006-01-02 15:04:05"))
	})

	http.HandleFunc("/page/system/uptime", func(w http.ResponseWriter, req *http.Request) {
		days, hours, mins, err := host.SystemUptime()
		AutoRender(w, fmt.Sprintf("%d days %d hours %d minutes", days, hours, mins), err)
	})

	http.HandleFunc("/proc/system/uptime", func(w http.ResponseWriter, req *http.Request) {
		days, hours, mins, err := host.SystemUptime()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		RenderDataJson(w, map[string]interface{}{
			"days":  days,
			"hours": hours,
			"mins":  mins,
		})
	})

	http.HandleFunc("/page/system/loadavg", func(w http.ResponseWriter, req *http.Request) {
		cpuNum := runtime.NumCPU()
		load, err := load.Avg()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		ret := [3][2]interface{}{
			[2]interface{}{load.Load1, int64(load.Load1 * 100.0 / float64(cpuNum))},
			[2]interface{}{load.Load5, int64(load.Load5 * 100.0 / float64(cpuNum))},
			[2]interface{}{load.Load15, int64(load.Load15 * 100.0 / float64(cpuNum))},
		}
		RenderDataJson(w, ret)
	})

	http.HandleFunc("/proc/system/loadavg", func(w http.ResponseWriter, req *http.Request) {
		data, err := load.Avg()
		AutoRender(w, data, err)
	})

}
