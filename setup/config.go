package setup

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

// GetConfig creates a config for the specified deployment
func GetConfig(deployment string, config *any) error {
	fileName := fmt.Sprintf("%s-config.json", deployment)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not find config file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Err(err).Msg("could not close config during defer")
		}
	}(file)
	decoder := json.NewDecoder(file)
	return decoder.Decode(config)
}
