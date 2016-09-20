package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"sync"
	"../tools/disk"
)

var (
	diskStatsMap = make(map[string][2]*disk.Statfs_t)
	dsLock       = new(sync.RWMutex)
)



func DiskIOMetrics() (L []*model.MetricValue) {

	dsList, err := disk.IOCounters()
	if err != nil {
		log.Println(err)
		return
	}

	for _, ds := range dsList {


		device := "device=" + ds.Name

		L = append(L, CounterValue("disk.io.read_requests", ds.ReadCount, device))
		L = append(L, CounterValue("disk.io.read_merged", ds.MergedReadCount, device))
		L = append(L, CounterValue("disk.io.read_bytes", ds.ReadBytes, device))
		L = append(L, CounterValue("disk.io.read_time", ds.ReadTime, device))
		L = append(L, CounterValue("disk.io.write_requests", ds.WriteCount, device))
		L = append(L, CounterValue("disk.io.write_merged", ds.MergedWriteCount, device))
		L = append(L, CounterValue("disk.io.write_bytes", ds.WriteBytes, device))
		L = append(L, CounterValue("disk.io.write_time", ds.WriteTime, device))
		L = append(L, CounterValue("disk.io.ios_in_progress", ds.IopsInProgress, device))
		L = append(L, CounterValue("disk.io.iotime", ds.IoTime, device))
	}
	return
}
