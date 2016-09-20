package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"../tools/load"
)

func LoadAvgMetrics() []*model.MetricValue {
	load, err := load.Avg()
	if err != nil {
		log.Println(err)
		return nil
	}

	return []*model.MetricValue{
		GaugeValue("load.1min", load.Load1),
		GaugeValue("load.5min", load.Load5),
		GaugeValue("load.15min", load.Load15),
	}

}
