package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/tokuhirom/doshin/sender"
	"github.com/tokuhirom/doshin/watcher"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

type Alert struct {
	Http watcher.HttpConfig `yaml:"http"`
	Net  watcher.NetConfig  `yaml:"net"`
}

type Config struct {
	Watch        Alert                     `yaml:"watch"`
	AlertManager sender.AlertManagerConfig `yaml:"alert_manager"`
	Interval     string                    `yaml:"interval"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func run(config Config) {
	ch := sender.StartAlertManagerWorker(config.AlertManager)

	interval := watcher.ParseDuration(60*time.Second, config.Interval)
	throttle := time.Tick(interval)

	for {
		// read it
		var wg sync.WaitGroup

		watcher.StartHttpWatchers(config.Watch.Http, ch, &wg)
		watcher.StartNetWatchers(config.Watch.Net, ch, wg)

		wg.Wait()

		<-throttle
	}
}

func main() {
	configFileName := flag.String("c", "config.yml", "path to config file")
	level := flag.String("log-level", "info", "logrus log level")
	flag.Parse()

	lvl, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)

	// read config yaml
	yamlFile, err := ioutil.ReadFile(*configFileName)
	if err != nil {
		log.Fatal(err)
	}

	// parse config file
	var config Config
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		log.Fatal(err)
	}

	run(config)
}
