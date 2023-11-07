package common

var (
	internalCfg *ServerConfig
	system      *System
)

func SetConfig(cfg *ServerConfig) {
	internalCfg = cfg
}

func GetConfig() *ServerConfig {
	return internalCfg
}

func GetSystem() *System {
	return system
}

func SetSystem(sys *System) {
	system = sys
}
