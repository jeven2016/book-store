package common

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type StartupParams struct {
	EnableMongodb bool
	EnableRedis   bool
	EnableEtcd    bool
	Config        Config
	ShutdownHook  func() error
}

func Startup(params *StartupParams) *System {
	// 创建一个全局的App
	sys := &System{}
	sys.Config = params.Config

	ctx := context.Background()
	sc := params.Config.GetServerConfig()

	// log初始化
	logger := SetupLog(sc.ApplicationName, sc.LogSetting)
	sys.Log = logger

	if params.EnableRedis {
		// 初始化redis
		redisClient, err := NewRedis(sc.Redis)
		if err != nil {
			sys.Log.Error("failed to initialize for redis", zap.Error(err))
			shutdown(ctx, sys, params)
			return nil
		} else {
			sys.Log.Info("Connecting to redis successfully")
			sys.RedisClient = redisClient
		}
	}

	if params.EnableMongodb {
		// 初始化Mongodb
		mongoClient := &MongoClient{Config: sc.Mongo}
		if err := mongoClient.StartInit(); err != nil {
			sys.Log.Error("failed to connect mongodb", zap.Error(err))
			shutdown(ctx, sys, params)
			return nil
		} else {
			sys.Log.Info("Connecting to mongodb successfully")
			sys.MongoClient = mongoClient
		}
	}

	//init a routine pool
	pool, err := ants.NewPool(sc.TaskPoolSetting.Capacity)
	if err != nil {
		sys.Log.Error("unable to init a routine pool", zap.Error(err))
		shutdown(ctx, sys, params)
		return nil
	} else {
		sys.Log.Info("task pool initialized successfully")
	}
	sys.TaskPool = pool

	if params.EnableEtcd {
		//submit a task to register this service
		if err = sys.RegisterService(sc); err != nil {
			sys.Log.Error("failed to register service in etcd", zap.String("app", sc.ApplicationName), zap.Error(err))
			shutdown(ctx, sys, params)
			return nil
		} else {
			sys.Log.Info("service registered in etcd", zap.String("app", sc.ApplicationName))
		}
	}

	sys.Log.Info("server started successfully")
	exitChan := make(chan os.Signal)

	err = sys.TaskPool.Submit(func() {
		// kill (no param) default send syscanll.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
		signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)
		<-exitChan
		shutdown(ctx, sys, params)
	})
	if err != nil {
		sys.Log.Info("unable to submit a shutdown hook", zap.Error(err))
		return nil
	}
	return sys
}

func shutdown(ctx context.Context, sys *System, params *StartupParams) {
	sys.Log.Info("server is shutting down")

	if params.ShutdownHook != nil {
		if err := params.ShutdownHook(); err != nil {
			sys.Log.Warn("an error occurs while calling shutdown hook", zap.Error(err))
		}
	}

	if sys.RedisClient != nil {
		if err := sys.RedisClient.Client.Close(); err != nil {
			sys.Log.Warn("an error occurs while closing redis's connection", zap.Error(err))
		}
	}

	if sys.MongoClient != nil {
		if err := sys.MongoClient.Client.Disconnect(ctx); err != nil {
			sys.Log.Warn("an error occurs while closing mongodb's connection", zap.Error(err))
		}
	}

	//if sys.ServiceRegister != nil {
	//	if sys.ServiceRegister != nil {
	//		if err := sys.ServiceRegister.Cancel(ctx); err != nil {
	//			sys.Log.Warn("an error occurs while closing etcd's connection", zap.Error(err))
	//		}
	//	}
	//}
	if sys.TaskPool != nil {
		sys.TaskPool.Release()
	}
	sys.Log.Info("shutdown completed")
	os.Exit(0)
}
