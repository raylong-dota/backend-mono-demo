package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

const (
	env_config_prefix = "ORBIT_"
)

type Config = config.Config

func New(path string) Config {
	c := config.New(
		config.WithSource(
			file.NewSource(path),
			// 环境变量 ORBIT_<CONFIG_NAME>, 会填充配置文件中的 ${CONFIG_NAME} 占位符
			// 比如， ORBIT_LOGGER_LEVEL，会填充配置文件中 ${LOGGER_LEVEL:info} 占位符
			// :info， 表示，如果没有填充值时候的默认值
			env.NewSource(env_config_prefix),
		),
	)
	return c
}
