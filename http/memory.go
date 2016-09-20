package http

import (
	"net/http"
	"../tools/mem"
)

func configMemoryRoutes() {
	http.HandleFunc("/page/memory", func(w http.ResponseWriter, r *http.Request) {
		m, err := mem.VirtualMemory()

		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		memFree := m.Free + m.Buffers + m.Cached
		memUsed := m.Total - memFree
		var t uint64 = 1024 * 1024
		RenderDataJson(w, []interface{}{m.Total / t, memUsed / t, memFree / t})
	})

	http.HandleFunc("/proc/memory", func(w http.ResponseWriter, r *http.Request) {
		m, err := mem.VirtualMemory()

		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		memFree := m.Free + m.Buffers + m.Cached
		memUsed := m.Total - memFree

		RenderDataJson(w, map[string]interface{}{
			"total": m.Total,
			"free":  memFree,
			"used":  memUsed,
		})
	})
}
