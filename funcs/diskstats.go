package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"sync"
	"../tools/disk"
	"../tools/host"
	"strings"
	"regexp"
	"strconv"
	"fmt"
)

//var (
//	diskStatsMap = make(map[string][2]*disk.Statfs_t)
//	dsLock       = new(sync.RWMutex)
//)


var (
	diskStatsMap = make(map[string][2]*disk.IOCountersStat)
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

func IOStatsForPage2() (L [][]string) {
	dsLock.RLock()
	defer dsLock.RUnlock()
	var content, err = host.Cmdexec("iostat -d ")
	if err != nil {
		log.Println(err)
		return
	}
	reg := regexp.MustCompile(`[ ]+`)
	content= reg.ReplaceAllString(content, " ")
	cmdOuts := strings.Split(content,"\n")
	if len(cmdOuts) <1 {
		return
	}
	disks := strings.Split(cmdOuts[0]," ")
	for _, ds := range disks {
		deviceName :=ds
		if strings.Contains(deviceName,"disk") {
			diskName := strings.Replace(deviceName,"/dev/","",-1)
			cmdOut,err :=host.CmdexecIostat(diskName)
			if err != nil{
				continue
			}
			if len(cmdOut) >1 {
				i, err := strconv.ParseFloat(cmdOut[2], 64)
				if err != nil {
					continue
				}
				item := []string{
					ds,
					"/",
					"/",
					"/",
					"/",
					cmdOut[0],
					strconv.FormatFloat(i*1024,'G',10,64),
					"/",                                             // avgrq-sz: delta(rsect+wsect)/delta(rio+wio)
					"/", // avgqu-sz: delta(aveq)/s/1000
					"/",                                                // await: delta(ruse+wuse)/delta(rio+wio)
					"/",                                                // svctm: delta(use)/delta(rio+wio)
					"/",                                  // %util: delta(use)/s/1000 * 100%
				}
				L = append(L, item)
			}

		}
	}

	return
}




func UpdateDiskStats() error {
	dsList, err := disk.IOCounters()
	if err != nil {
		return err
	}

	dsLock.Lock()
	defer dsLock.Unlock()
	ret := make([]*disk.IOCountersStat, 0)
	for _, ds := range dsList {
		item := &disk.IOCountersStat{}
		*item = ds
		ret = append(ret, item)
	}
	for i := 0; i < len(ret); i++ {
		device := ret[i].Name
		diskStatsMap[device] = [2]*disk.IOCountersStat{ret[i], diskStatsMap[device][0]}
	}
	return nil
}

func IOReadRequests(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].ReadCount - arr[1].ReadCount
}

func IOReadMerged(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].MergedReadCount - arr[1].MergedReadCount
}

func IOReadSectors(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].ReadBytes - arr[1].ReadBytes
}

func IOMsecRead(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].ReadTime - arr[1].ReadTime
}

func IOWriteRequests(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].WriteCount - arr[1].WriteCount
}

func IOWriteMerged(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].MergedWriteCount - arr[1].MergedWriteCount
}

func IOWriteSectors(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].WriteBytes - arr[1].WriteBytes
}

func IOMsecWrite(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].WriteTime - arr[1].WriteTime
}

func IOMsecTotal(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].IopsInProgress - arr[1].IopsInProgress
}

func IOMsecWeightedTotal(arr [2]*disk.IOCountersStat) uint64 {
	return arr[0].WeightedIO - arr[1].WeightedIO
}

//func TS(arr [2]*disk.IOCountersStat) uint64 {
//	return uint64(arr[0].TS.Sub(arr[1].TS).Nanoseconds() / 1000000)
//}

func IODelta(device string, f func([2]*disk.IOCountersStat) uint64) uint64 {
	val, ok := diskStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

func IOStatsForPage() (L [][]string) {
	UpdateDiskStats()
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device, _ := range diskStatsMap {
		//if !ShouldHandleDevice(device) {
		//	continue
		//}

		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)

		delta_rsec := IODelta(device, IOReadSectors)
		delta_wsec := IODelta(device, IOWriteSectors)

		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		use := IODelta(device, IOMsecTotal)
		n_io := rio + wio
		avgrq_sz := 0.0
		await := 0.0
		svctm := 0.0
		if n_io != 0 {
			avgrq_sz = float64(delta_rsec+delta_wsec) / float64(n_io)
			await = float64(ruse+wuse) / float64(n_io)
			svctm = float64(use) / float64(n_io)
		}

		item := []string{
			device,
			fmt.Sprintf("%d", IODelta(device, IOReadMerged)),
			fmt.Sprintf("%d", IODelta(device, IOWriteMerged)),
			fmt.Sprintf("%d", rio),
			fmt.Sprintf("%d", wio),
			fmt.Sprintf("%.2f", float64(delta_rsec)/2.0),
			fmt.Sprintf("%.2f", float64(delta_wsec)/2.0),
			fmt.Sprintf("%.2f", avgrq_sz),                                             // avgrq-sz: delta(rsect+wsect)/delta(rio+wio)
			fmt.Sprintf("%.2f", float64(IODelta(device, IOMsecWeightedTotal))/1000.0), // avgqu-sz: delta(aveq)/s/1000
			fmt.Sprintf("%.2f", await),                                                // await: delta(ruse+wuse)/delta(rio+wio)
			fmt.Sprintf("%.2f", svctm),                                                // svctm: delta(use)/delta(rio+wio)
			fmt.Sprintf("%.2f%%", float64(use)/10.0),                                  // %util: delta(use)/s/1000 * 100%
		}
		L = append(L, item)
	}

	return
}

func ShouldHandleDevice(device string) bool {
	normal := len(device) == 3 && (strings.HasPrefix(device, "sd") || strings.HasPrefix(device, "vd"))
	aws := len(device) >= 4 && strings.HasPrefix(device, "xvd")
	return normal || aws
}
