package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/oklog/run"
	"github.com/prometheus/common/model"
)

type Config struct {
	GlobalConfig  GlobalConfig    `yaml:"global"`
	ScrapeConfigs []*ScrapeConfig `yaml:"scrape_configs,omitempty"`

	original string
}

type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval model.Duration `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout model.Duration `yaml:"scrape_timeout,omitempty"`
	// How frequently to evaluate rules by default.
	EvaluationInterval model.Duration `yaml:"evaluation_interval,omitempty"`
	// The labels to add to any timeseries that this Prometheus instance scrapes.
	ExternalLabels model.LabelSet `yaml:"external_labels,omitempty"`
}

type ScrapeConfig struct {
	// The job name to which the job label is set by default.
	JobName string `yaml:"job_name"`
}

type closeOnce struct {
	C     chan struct{}
	once  sync.Once
	Close func()
}

var logger = log.New(os.Stdout, "", log.LstdFlags)

func main() {

	reloadReady := &closeOnce{
		C: make(chan struct{}),
	}

	reloadReady.Close = func() {
		reloadReady.once.Do(func() {
			close(reloadReady.C)
		})
	}

	var g run.Group
	{
		hup := make(chan os.Signal, 1)
		signal.Notify(hup, syscall.SIGHUP)
		cancel := make(chan struct{})
		g.Add(func() error {
			<-reloadReady.C

			for {
				select {
				case <-hup:
					reloadConfig()
				case <-cancel:
					fmt.Println("been canceled")
					return nil
				}
			}

		}, func(err error) {
			// hey, this is my first time to see "struct{}{}"
			cancel <- struct{}{}
		})

	}
	{
		// Initial configuration loading
		cancel := make(chan struct{})
		g.Add(func() error {
			fmt.Println("Initial configuration loading")
			reloadConfig()
			reloadReady.Close()
			<-cancel
			return nil
		}, func(err error) {
			close(cancel)
		})
	}

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}

}

func reloadConfig() {
	scrapeManager.ApplyConfig
	fmt.Println("reloadConfig OK")
}
