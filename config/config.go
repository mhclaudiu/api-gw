package config

import (
	"api-gw/functions"
	"api-gw/logging"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var appLog logging.Log

func (cfg *CFG) Load(path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	buffCfg := CFG{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&buffCfg)
	if err != nil {
		log.Panic(err)
	}

	cfg.Update(buffCfg)

	if len(cfg.APP.Name) < 4 {

		appLog.Add(logging.Entry{
			Event: fmt.Sprintf("Invalid App Name - too short - '%s'", buffCfg.APP.Name),
			Exit:  true,
		})
	}

	if len(cfg.API.Path) < 1 {

		appLog.Add(logging.Entry{
			Event: fmt.Sprintf("Invalid API Main Path - too short - '%s'", buffCfg.API.Path),
			Exit:  true,
		})
	}

	appLog.Add(logging.Entry{
		Event: fmt.Sprintf("Succesfully loaded config file '%s':\n%+v\n", path, functions.Dump(cfg)),
	})

}

func (cfg *CFG) Update(ss CFG) {
	*cfg = ss
}
