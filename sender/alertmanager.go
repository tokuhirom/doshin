package sender

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
)

// https://github.com/prometheus/alertmanager/blob/master/api.go

type AlertManagerResult struct {
	Summary     string
	Description string
}

// Schema for configuration file
type AlertManagerConfig struct {
	Labels map[string]string `yaml:"labels"`
	Url    string            `yaml:"url"`
}

type AlertManagerRequest struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	GeneratorURL string            `json:"generatorURL"`
}

func NewResult(summary, description string) *AlertManagerResult {
	return &AlertManagerResult{summary, description}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func StartAlertManagerWorker(config AlertManagerConfig) chan *AlertManagerResult {
	ch := make(chan *AlertManagerResult)

	go func() {
		log.Info("Start alert manager sender thread: ", config.Url)
		for {
			result := <-ch
			log.Debug("Sending result to alert manager : ", result)

			generatorURL := "http://example.com/" + randStringBytes(16)
			req := &AlertManagerRequest{
				config.Labels,
				map[string]string{
					"summary":     result.Summary,
					"description": result.Description},
				generatorURL}
			reqs := make([]*AlertManagerRequest, 1)
			reqs[0] = req
			json, err := json.Marshal(reqs)
			log.Info(req.GeneratorURL)
			if err != nil {
				log.Fatal(err)
			}

			response, err := http.Post(config.Url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(json)))
			if err == nil {
				defer response.Body.Close()

				if 200 <= response.StatusCode && response.StatusCode < 300 {
					log.Info("Sent alert to alertmanager: ", config.Url, " ", response.Status, " ", result.Summary)
				} else {
					body, err := ioutil.ReadAll(response.Body)
					if err != nil {
						log.Warn("Cannot send alert to alertmanager. And cannot read response body! ", config.Url, " ", response.Status, err)
					} else {
						log.Warn("Cannot send alert to alertmanager. ", config.Url, " ", response.Status, string(body))
					}
				}
			} else {
				log.Warn("Cannot send alert to alertmanager: ", config.Url, " ", err)
			}
		}
	}()

	return ch
}
