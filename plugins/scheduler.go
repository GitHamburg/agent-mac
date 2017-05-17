package plugins

import (
	"bytes"
	//"encoding/json"
	"../g"
	//"github.com/open-falcon/common/model"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (this *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				PluginRun(this.Plugin)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *PluginScheduler) Stop() {
	close(this.Quit)
}

func PluginRun(plugin *Plugin) {

	debug := g.Config().Debug

	timeout := plugin.Cycle*1000 - 500
	if debug {
		log.Println("plugin timeout:", timeout)
		log.Println("plugin dir path:", g.Config().Plugin.Dir)
		log.Println("plugin file path:", plugin.FilePath)
	}
	fpath := filepath.Join(g.Config().Plugin.Dir, plugin.FilePath)
	if debug {
		log.Println("plugin path:", fpath)
	}

	if !file.IsExist(fpath) {
		log.Println("no such plugin:", fpath)
		return
	}


	if debug {
		log.Println(fpath, "running...")
	}

	cmd := exec.Command(fpath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Start()

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)

	errStr := stderr.String()
	if errStr != "" {
		logFile := filepath.Join(g.Config().Plugin.LogDir, plugin.FilePath+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Printf("[ERROR] write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil && debug {
			log.Println("[INFO] timeout and kill process", fpath, "successfully")
		}

		if err != nil {
			log.Println("[ERROR] kill process", fpath, "occur error:", err)
		}

		return
	}

	if err != nil {
		log.Println("[ERROR] exec plugin", fpath, "fail. error:", err)
		return
	}

	//// exec successfully
	//data := stdout.Bytes()
	//if len(data) == 0 {
	//	if debug {
	//		log.Println("[DEBUG] stdout of", fpath, "is blank")
	//	}
	//	return
	//}
	//
	//var metrics []*model.MetricValue
	//err = json.Unmarshal(data, &metrics)
	//if err != nil {
	//	log.Printf("[ERROR] json.Unmarshal stdout of %s fail. error:%s stdout: \n%s\n", fpath, err, stdout.String())
	//	return
	//}
	//
	//g.SendToTransfer(metrics)

	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
