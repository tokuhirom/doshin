package watcher

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"time"
)

type NetTarget struct {
	Network  string `yaml:"network"`
	Address  string `yaml:"address"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
}

type NetConfig struct {
	Targets []NetTarget `yaml:"targets"`
}

func startWorker(target NetTarget) {
	interval := ParseDuration(60*time.Second, target.Interval)
	timeout := ParseDuration(1*time.Second, target.Timeout)

	log.Info("Watching HTTP url: ", target.Address, " interval: ", interval)
	go func() {
		for {
			start := time.Now()
			conn, err := net.DialTimeout(target.Network, target.Address, timeout)
			end := time.Now()
			sleep := interval - end.Sub(start)

			if err == nil {
				conn.Close()
				log.Debug("OK: ", target.Address, " sleep seconds: ", sleep)
			} else {
				log.Info("Cannot get ", target.Address, " : ", err, " sleep seconds: ", sleep)
			}
			time.Sleep(sleep)
		}
	}()
}

func StartNetWatchers(conf NetConfig) {
	for _, target := range conf.Targets {
		startWorker(target)
	}
}
