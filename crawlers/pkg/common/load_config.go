package common

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"log"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func LoadConfig(internalCfg []byte, config Config) error {

	//load internal config
	if err := k.Load(rawbytes.Provider(internalCfg), yaml.Parser()); err != nil {
		return err
	}

	// load external configs
	for _, f := range ConfigFiles {
		if exists, err := IsFileExists(f); err != nil {
			log.Printf(f + " ignored: " + err.Error())
			continue
		} else if exists {
			if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
				PrintCmdErr(err)
				continue
			}
		}
	}

	if err := k.Unmarshal("", config); err != nil {
		return err
	}

	return nil
}
