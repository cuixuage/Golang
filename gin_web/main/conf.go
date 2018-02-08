package main

import (
"os"
"fmt"
"github.com/viper"
"path/filepath"
"strings"
"code.byted.org/gin/ginex"
"code.byted.org/gopkg/logs"
"time"
)

type AppConfig struct {
	InnerTest bool
	EnableMetrics bool
	BackendTimeout time.Duration
	BackendConnTimeout time.Duration
	SuspectImeiThreshold int64
	BlockImeiThreshold int64
	BlockImeiStatDuration time.Duration
	*viper.Viper
}

var appConf *AppConfig

func parseConf() {
	v := viper.New()
	v.SetEnvPrefix("GIN")

	confFile := filepath.Join(ginex.ConfDir(), strings.Replace(ginex.PSM(), ".", "_", -1)+".yaml")
	v.SetConfigFile(confFile)
	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load app config: %s, %s\n", confFile, err)
		os.Exit(-1)
	}
	mode := "Develop"
	if ginex.Product() {
		mode = "Product"
	}

	vv := v.Sub(mode)
	if vv == nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config sub module: %s\n", mode)
		os.Exit(-1)
	} else {
		vv.SetDefault("InnerTest", true)
		vv.SetDefault("EnableMetrics", false)
		vv.SetDefault("BackendTimeout", 500 * time.Millisecond)
		vv.SetDefault("BackendConnTimeout", 50 * time.Millisecond)
		vv.SetDefault("SuspectImeiThreshold", 100)
		vv.SetDefault("BlockImeiThreshold", 200)
		vv.SetDefault("BlockImeiStatDuration", 1 * time.Hour)
	}

	appConf = &AppConfig{}
	if err := vv.Unmarshal(appConf); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal app config: %s\n", err)
		os.Exit(-1)
	}
	appConf.Viper = vv
	logs.Info("parse appConf %+v", appConf)
}
