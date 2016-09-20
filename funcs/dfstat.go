package funcs

import (
	"fmt"
	"github.com/open-falcon/common/model"
	"log"
	"../tools/disk"
)

func DeviceMetrics() (L []*model.MetricValue) {

	var diskTotal uint64 = 0
	var diskUsed uint64 = 0

	var du, err = disk.Partitions(false)
	if err != nil {
		log.Println(err)
	}
	for _, ds := range du {

		var dusage, err = disk.Usage(ds.Mountpoint)
		if err != nil {
			log.Println(err)
		}

		diskTotal += dusage.Total
		diskUsed += dusage.Used

		tags := fmt.Sprintf("mount=%s,fstype=%s", dusage.Path, dusage.Fstype)
		L = append(L, GaugeValue("df.bytes.total", dusage.Total, tags))
		L = append(L, GaugeValue("df.bytes.used", dusage.Used, tags))
		L = append(L, GaugeValue("df.bytes.free", dusage.Free, tags))
		L = append(L, GaugeValue("df.bytes.used.percent", dusage.UsedPercent, tags))
		L = append(L, GaugeValue("df.bytes.free.percent", 100-dusage.UsedPercent, tags))
		L = append(L, GaugeValue("df.inodes.total", dusage.InodesTotal, tags))
		L = append(L, GaugeValue("df.inodes.used", dusage.InodesUsed, tags))
		L = append(L, GaugeValue("df.inodes.free", dusage.InodesFree, tags))
		L = append(L, GaugeValue("df.inodes.used.percent", dusage.InodesUsedPercent, tags))
		L = append(L, GaugeValue("df.inodes.free.percent", 100-dusage.InodesUsedPercent, tags))
	}

	if len(L) > 0 && diskTotal > 0 {
		L = append(L, GaugeValue("df.statistics.total", float64(diskTotal)))
		L = append(L, GaugeValue("df.statistics.used", float64(diskUsed)))
		L = append(L, GaugeValue("df.statistics.used.percent", float64(diskUsed)*100.0/float64(diskTotal)))
	}

	return
}
