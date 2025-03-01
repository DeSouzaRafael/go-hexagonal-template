package util

import "github.com/DeSouzaRafael/go-hexagonal-template/internal/config"

func CurrentExecutionEnvironmentProduction() bool {
	return config.AppConfig.Environment == "prd"
}
