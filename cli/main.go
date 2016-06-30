package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/tokuhirom/doshin/sender"
	"github.com/tokuhirom/doshin/watcher"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"time"
)

type Alert struct {
	Http watcher.HttpConfig `yaml:"http"`
	Net  watcher.NetConfig  `yaml:"net"`
}

type Config struct {
	Watch        Alert                     `yaml:"watch"`
	AlertManager sender.AlertManagerConfig `yaml:"alert_manager"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

	ch := sender.StartAlertManagerWorker(config.AlertManager)

	// read it
	watcher.StartHttpWatchers(config.Watch.Http, ch)
	watcher.StartNetWatchers(config.Watch.Net)

	// Sleep forever
	select {}
}
