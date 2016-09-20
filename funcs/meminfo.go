package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"../tools/mem"
)

func MemMetrics() []*model.MetricValue {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return nil
	}

	memFree := m.Free + m.Buffers + m.Cached
	memUsed := m.Total - memFree

	pmemFree := 0.0
	pmemUsed := 0.0
	if m.Total != 0 {
		pmemFree = float64(memFree) * 100.0 / float64(m.Total)
		pmemUsed = float64(memUsed) * 100.0 / float64(m.Total)
	}

	//pswapFree := 0.0
	//pswapUsed := 0.0
	//if m.SwapTotal != 0 {
	//	pswapFree = float64(m.SwapFree) * 100.0 / float64(m.SwapTotal)
	//	pswapUsed = float64(m.SwapUsed) * 100.0 / float64(m.SwapTotal)
	//}

	return []*model.MetricValue{
		GaugeValue("mem.memtotal", m.Total),
		GaugeValue("mem.memused", memUsed),
		GaugeValue("mem.memfree", memFree),
		//GaugeValue("mem.swaptotal", m.SwapTotal),
		//GaugeValue("mem.swapused", m.SwapUsed),
		//GaugeValue("mem.swapfree", m.SwapFree),
		GaugeValue("mem.memfree.percent", pmemFree),
		GaugeValue("mem.memused.percent", pmemUsed),
		//GaugeValue("mem.swapfree.percent", pswapFree),
		//GaugeValue("mem.swapused.percent", pswapUsed),
	}

}
