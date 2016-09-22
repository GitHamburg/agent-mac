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

func IOStatsForPage() (L [][]string) {
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
