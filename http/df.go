package http

import (
	"fmt"
	"github.com/toolkits/core"
	"net/http"
	"log"
	"../tools/disk"
)

func configDfRoutes() {
	http.HandleFunc("/page/df", func(w http.ResponseWriter, r *http.Request) {
		var diskTotal uint64 = 0
		var diskUsed uint64 = 0

		var mountPoints, err = disk.Partitions(false)
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		var ret [][]interface{} = make([][]interface{}, 0)
		for _, ds := range mountPoints {

			var dusage, err = disk.Usage(ds.Mountpoint)
			if err != nil {
				log.Println(err)
			}

			diskTotal += dusage.Total
			diskUsed += dusage.Used

			if err == nil {
				ret = append(ret,
					[]interface{}{
						dusage.Fstype,
						core.ReadableSize(float64(dusage.Total)),
						core.ReadableSize(float64(dusage.Used)),
						core.ReadableSize(float64(dusage.Free)),
						fmt.Sprintf("%.1f%%", dusage.UsedPercent),
						ds.Mountpoint,
						core.ReadableSize(float64(dusage.InodesTotal)),
						core.ReadableSize(float64(dusage.InodesUsed)),
						core.ReadableSize(float64(dusage.InodesFree)),
						fmt.Sprintf("%.1f%%", dusage.InodesUsedPercent),
						ds.Fstype,
					})
			}
		}

		RenderDataJson(w, ret)
	})
}
