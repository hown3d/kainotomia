package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	JobImage  string `env:"KAINOTOMIA_JOB_IMAGE"`
	Namespace string `env:"KAINOTOMIA_NAMESPACE"`
}

func Parse() Config {
	c := &Config{}
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}
	return *c
}
