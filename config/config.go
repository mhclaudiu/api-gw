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

	if cfg.LOG.Enabled {

		cfg.CreateLogsDir()
	}
}

func (cfg *CFG) Update(ss CFG) {
	*cfg = ss
}

func (cfg *CFG) CreateLogsDir() {

	if _, err := os.Stat(cfg.LOG.Dir); !os.IsNotExist(err) {
		return
	}

	err := os.MkdirAll(cfg.LOG.Dir, 0755)
	if err != nil {

		appLog.Add(logging.Entry{
			Event: fmt.Sprintf("Error creating Logs Directory '%s' - %s", cfg.LOG.Dir, err.Error()),
			Code:  logging.CONST_CODE_ERROR,
			Exit:  true,
		})
	}

	appLog.Add(logging.Entry{
		Event: fmt.Sprintf("Logs Directory '%s' successfully created", cfg.LOG.Dir),
		Code:  logging.CONST_CODE_INFO,
	})
}
