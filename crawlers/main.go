package main

import (
	"context"
	"crawlers/pkg/api"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/website"
	_ "embed"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"net/http"
)

//go:embed config/internal_conf.yaml
var configFile string

const softwareVersion = "0.1"
const flagName = "config"

var extraConfigFile *string

// @title  crawler文档
// @version 0.2
// @description  crawler接口参考文档
// @termsOfService only for internal use
// @BasePath /api/v1/
// @query.collection.format multi
func main() {
	//printBanner()
	run()
}

func run() {
	var server *http.Server
	var rootCmd = &cobra.Command{
		Version: softwareVersion,
		Use:     "crawlers",
		Short:   "crawlers",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &common.ServerConfig{}
			if err := common.LoadConfig([]byte(configFile), cfg, extraConfigFile); err != nil {
				common.PrintCmdErr(err)
				return
			}
			common.SetConfig(cfg)

			//global context
			ctx, cancelFunc := context.WithCancel(context.Background())
			sys := systemInit(cfg, cancelFunc, server, ctx)
			if sys != nil {
				website.RegisterProcessors()
				dao.InitDao(ctx)

				if err := website.RegisterStream(ctx); err != nil {
					zap.L().Error("failed to register streams", zap.Error(err))
					return
				}

				engine := api.RegisterEndpoints()

				// run as a web server
				bindAddr := fmt.Sprintf("%v:%v", cfg.Http.Address, cfg.Http.Port)
				zap.L().Sugar().Info("server listens on ", bindAddr)
				server = &http.Server{Addr: bindAddr, Handler: engine}

				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					zap.L().Error("unable to start web server", zap.Error(err))
				}
			}
		},
	}

	// 配置文件的绝对路径
	extraConfigFile = rootCmd.Flags().StringP(flagName, "c", "", "the absolute path of yaml config file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func systemInit(cfg *common.ServerConfig, cancelFunc context.CancelFunc, server *http.Server, ctx context.Context) *common.System {
	sys := common.Startup(&common.StartupParams{
		EnableEtcd:    false,
		EnableMongodb: true,
		EnableRedis:   true,
		Config:        cfg,
		PreShutdown: func() error {
			cancelFunc()
			return nil
		},
		PostShutdown: func() error {
			if server != nil {
				if err := server.Shutdown(ctx); err != nil {
					zap.L().Error("unable to shut web server down", zap.Error(err))
					return err
				}
			}
			zap.S().Info("web server shuts down")
			return nil
		},
	})
	if sys != nil {
		common.SetSystem(sys)
	}
	return sys
}
