package watcher

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/tokuhirom/doshin/sender"
	"net"
	"sync"
	"time"
)

type NetTarget struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
	Timeout string `yaml:"timeout"`
}

type NetConfig struct {
	Targets []NetTarget `yaml:"targets"`
}

func startWorker(target NetTarget, ch chan *sender.AlertManagerResult) {
	timeout := ParseDuration(1*time.Second, target.Timeout)

	log.Info("Watching HTTP url: ", target.Address)
	conn, err := net.DialTimeout(target.Network, target.Address, timeout)

	if err == nil {
		conn.Close()
		log.Debug("OK: ", target.Address)
	} else {
		log.Info("Cannot get ", target.Address, " : ", err)
		ch <- sender.NewResult(
			fmt.Sprintf("Cannot get %v(%v): %v", target.Address, target.Network, err),
			"An error occurred in Net watcher agent")
	}
}

func StartNetWatchers(conf NetConfig, ch chan *sender.AlertManagerResult, wg sync.WaitGroup) {
	for _, target := range conf.Targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startWorker(target, ch)
		}()
	}
}
