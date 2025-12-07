package settings

import (
	"os"
	"encoding/json"
	"errors"
)


type ImageSettingJson struct {
	Path     string `json:"path"`
	HotSpotX int16  `json:"hotSpotX"`
	HotSpotY int16  `json:"hotSpotY"`
}

type SettingsJson struct {
	Images       []ImageSettingJson `json:"images"`
	FrameIndexes []uint32           `json:"frameIndexes"`
	Rates        []uint32           `json:"rates"`
	Output       string             `json:"output"`
	ResourceType string             `json:"resourceType"`
}

func LoadSettingsJson(path string) (SettingsJson, error) {
	var settings SettingsJson
	
	file, err := os.Open(path)
	if err != nil {
		return SettingsJson{}, err
	}
	defer file.Close()
	
	if err := json.NewDecoder(file).Decode(&settings); err != nil {
		return SettingsJson{}, err
	}
	if settings.Images == nil {
		return SettingsJson{}, errors.New("images are not set")
	}
	if settings.FrameIndexes == nil {
		return SettingsJson{}, errors.New("frame indexes are not set")
	}
	if settings.Rates == nil {
		return SettingsJson{}, errors.New("rates are not set")
	}
	if settings.Output == "" {
		return SettingsJson{}, errors.New("output is not set")
	}
	if settings.ResourceType != "" && settings.ResourceType != "cursor" && settings.ResourceType != "icon" {
		return SettingsJson{}, errors.New("invalid resource type")
	}
	return settings, nil
}