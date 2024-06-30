package config

type Config interface {
	UpdateConfigFile(key string, value interface{})
	CommonConfig() Common
	MeteoblueConfig() Meteoblue
}
