package common

import (
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"sync"
)

var sysOnce = sync.Once{}
var lock = sync.Mutex{}

type System struct {
	// Log /**全局Logger*/
	Log *zap.Logger

	RedisClient *Redis

	MongoClient *MongoClient

	//ServiceRegister Register

	Config Config

	TaskPool *ants.Pool

	collectionMap map[string]*mongo.Collection
}

func (s *System) RegisterService(cfg *ServerConfig) error {
	//registerParam := &RegisterParam{
	//	ServiceName:           cfg.ApplicationName,
	//	ServiceHost:           cfg.Http.Address,
	//	ServicePort:           cfg.Http.Port,
	//	RefreshSeconds:        cfg.Registration.Etcd.RefreshSeconds,
	//	ConnectTimeoutSeconds: cfg.Registration.Etcd.ConnectTimeoutSeconds,
	//}
	//
	////register service to etcd
	//if register, err := NewRegister(cfg.Registration.Etcd.Endpoints, s.Log, registerParam, s.TaskPool); err != nil {
	//	return err
	//} else {
	//	s.ServiceRegister = register
	//	err = s.ServiceRegister.Register(context.Background())
	//
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

//func (s *System) GetServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
//	return s.ServiceRegister.ListServiceAddresses(ctx, serviceName)
//}

func (s *System) GetCollection(name string) *mongo.Collection {
	sysOnce.Do(func() {
		s.collectionMap = make(map[string]*mongo.Collection)
	})
	if collection, ok := s.collectionMap[name]; ok {
		return collection
	} else {
		lock.Lock()
		defer lock.Unlock()

		if collection, ok = s.collectionMap[name]; ok {
			return collection
		}
		s.collectionMap[name] = s.MongoClient.Db.Collection(name)
	}
	return s.collectionMap[name]
}
