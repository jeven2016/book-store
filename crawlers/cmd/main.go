package main

import (
	"context"
	"crawlers/pkg/api"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
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

			//global context
			ctx, cancelFunc := context.WithCancel(context.Background())
			sys := common.Startup(&common.StartupParams{
				EnableEtcd:    false,
				EnableMongodb: true,
				EnableRedis:   true,
				Config:        cfg,
				ShutdownHook: func() error {
					cancelFunc()
					return nil
				},
			})
			if sys != nil {
				common.SetSystem(sys)
				website.RegisterProcessors()
				dao.InitDao()

				if err := website.RegisterStream(ctx); err != nil {
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
