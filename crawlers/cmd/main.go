package main

import (
	"crawlers/pkg/api"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/stream"
	"crawlers/pkg/website"
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
)

//go:embed internal_conf.yaml
var configFile string

const softwareVersion = "0.0.1"
const flagName = "config"

func main() {
	//printBanner()
	run()
}

func run() {
	var rootCmd = &cobra.Command{
		Version: softwareVersion,
		Use:     "down",
		Short:   "crawlers",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &common.ServerConfig{}
			if err := common.LoadConfig([]byte(configFile), cfg); err != nil {
				common.PrintCmdErr(err)
				return
			}
			common.SetConfig(cfg)

			sys := common.Startup(&common.StartupParams{
				EnableEtcd:    false,
				EnableMongodb: true,
				EnableRedis:   true,
				Config:        cfg,
			})
			if sys != nil {
				common.SetSystem(sys)
				website.InitJobHandlers()
				dao.InitDao()

				if err := stream.RegisterStream(); err != nil {
					sys.Log.Error("failed to register streams", zap.Error(err))
					return
				}
				engine := api.RegisterEndpoints()

				// run as a web server
				bindAddr := fmt.Sprintf("%v:%v", cfg.Http.Address, cfg.Http.Port)
				sys.Log.Sugar().Info("server starts on ", bindAddr)
				if err := engine.Run(bindAddr); err != nil {
					sys.Log.Fatal("unable to start web server", zap.Error(err))
				}
			}
		},
	}

	// 配置文件的绝对路径
	rootCmd.Flags().StringP(flagName, "c", "", "the absolute path of yaml config file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func printBanner() {
	println(configFile)
}
