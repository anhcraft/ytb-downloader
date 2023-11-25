package settings

import (
	"encoding/json"
	"log"
	"os"
)

const SETTINGS_FILE = "settings.json"

var settings *Settings

func Get() *Settings {
	return settings
}

func Load() {
	ctn, err := os.ReadFile(SETTINGS_FILE)
	if err != nil {
		log.Printf("error reading settings file: %v\n", err)
		settings = NewSettings()
		return
	}
	if err := json.Unmarshal(ctn, &settings); err != nil {
		log.Printf("error unmarshalling settings file: %v\n", err)
		settings = NewSettings()
	}
}

func Save() {
	data, err := json.Marshal(settings)
	if err != nil {
		log.Printf("error marshalling settings file: %v\n", err)
		return
	}
	err = os.WriteFile(SETTINGS_FILE, data, 0777)
	if err != nil {
		log.Printf("error writing settings file: %v\n", err)
	}
}
