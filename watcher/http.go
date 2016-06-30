package watcher

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/tokuhirom/doshin/sender"
	"net/http"
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
	interval := ParseDuration(60*time.Second, target.Interval)
	timeout := ParseDuration(3*time.Second, target.Timeout)
	client := &http.Client{
		Timeout: timeout,
	}

	log.Info("Watching HTTP url: ", target.Url, " interval: ", interval, " timeout: ", timeout)
	go func() {
		for {
			start := time.Now()
			_, err := client.Get(target.Url)
			end := time.Now()
			sleep := interval - end.Sub(start)

			if err == nil {
				log.Debug("OK: ", target.Url, " sleep seconds: ", sleep)
			} else {
				log.Info("Cannot get ", target.Url, " : ", err, " sleep seconds: ", sleep)
				ch <- sender.NewResult(
					fmt.Sprintf("Cannot get %v : %v", target.Url, err),
					"An error occurred in HTTP watcher agent")
			}
			time.Sleep(sleep)
		}
	}()
}

func StartHttpWatchers(conf HttpConfig, ch chan *sender.AlertManagerResult) {
	for _, target := range conf.Targets {
		startHttpWorker(target, ch)
	}
}
