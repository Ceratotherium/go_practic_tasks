//go:build task_template

package main

type Config struct {
	PodsCount    int
	StartTimeout time.Duration
}

type Loader interface {
	Load() (Config, error)
}

type ConfigUpdater struct{}

func New(loader Loader, updateInterval time.Duration) (*ConfigUpdater, func()) {
	return nil, func() {}
}

func (c *ConfigUpdater) GetConfig() Config {
	return Config{}
}
