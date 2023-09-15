package config

import (
	"github.com/lgdzz/vingo-utils/vingo/db/mysql"
	"github.com/lgdzz/vingo-utils/vingo/db/redis"
)

type Config struct {
	System   System            `yaml:"system" json:"system"`
	Database mysql.MysqlConfig `yaml:"database" json:"database"`
	Redis    redis.Option      `yaml:"redis" json:"redis"`
}
