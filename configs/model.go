package configs

type ModelConfig struct {
	BaseConfig
	tPath string
}

func (ModelConfig) FromContext(c *cli.Context) ModelConfig {
	return ModelConfig{
		BaseConfig: BaseConfig{}.FromContext(c),
		tPath:      c.String(tPath),
	}
}
