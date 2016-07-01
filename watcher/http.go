package watcher

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/tokuhirom/doshin/sender"
	"net/http"
	"sync"
	"time"
)

type HttpTarget struct {
	Url      string `yaml:"url"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
}

type HttpConfig struct {
	Targets []HttpTarget `yaml:"targets"`
}

func startHttpWorker(target HttpTarget, ch chan *sender.AlertManagerResult) {
	timeout := ParseDuration(3*time.Second, target.Timeout)
	client := &http.Client{
		Timeout: timeout,
	}

	log.Info("Watching HTTP url: ", target.Url, " timeout: ", timeout)
	_, err := client.Get(target.Url)

	if err == nil {
		log.Debug("OK: ", target.Url)
	} else {
		log.Info("Cannot get ", target.Url, " : ", err)
		ch <- sender.NewResult(
			fmt.Sprintf("Cannot get %v : %v", target.Url, err),
			"An error occurred in HTTP watcher agent")
	}
}

func StartHttpWatchers(conf HttpConfig, ch chan *sender.AlertManagerResult, wg *sync.WaitGroup) {
	for _, target := range conf.Targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startHttpWorker(target, ch)
		}()
	}
}
