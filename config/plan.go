package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Plan struct {
	Name      string    `yaml:"name"`
	Target    Target    `yaml:"target"`
	Scheduler Scheduler `yaml:"scheduler"`
	S3        *S3       `yaml:"s3"`
	SMTP      *SMTP     `yaml:"smtp"`
	Slack     *Slack    `yaml:"slack"`
}

type Target struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
}

type Scheduler struct {
	Cron      string `yaml:"cron"`
	Retention int    `yaml:"retention"`
	Timeout   int    `yaml:"timeout"`
}

type S3 struct {
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"accessKey"`
	API       string `yaml:"api"`
	SecretKey string `yaml:"secretKey"`
	URL       string `yaml:"url"`
}

type SMTP struct {
	Server   string   `yaml:"server"`
	Port     string   `yaml:"port"`
	Password string   `yaml:"password"`
	Username string   `yaml:"username"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

type Slack struct {
	URL      string `yaml:"url"`
	Channel  string `yaml:"channel"`
	Username string `yaml:"username"`
}

func LoadPlans(dir string) ([]Plan, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "yml") || strings.Contains(path, "yaml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "reading from %v failed", dir)
	}

	plans := make([]Plan, 0)

	for _, path := range files {
		var plan Plan
		if strings.Contains(path, "yml") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, errors.Wrapf(err, "reading %v failed", path)
			}

			if err := yaml.Unmarshal(data, &plan); err != nil {
				return nil, errors.Wrapf(err, "parsering %v failed", path)
			}
			_, filename := filepath.Split(path)
			plan.Name = strings.TrimSuffix(filename, filepath.Ext(filename))
			plans = append(plans, plan)
		}
	}
	if len(plans) < 1 {
		return nil, errors.Errorf("No backup plans found in %v", dir)
	}

	return plans, nil
}
