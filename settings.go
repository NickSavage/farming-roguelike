package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const SETTINGS_PATH = "settings.json"

type Settings struct {
	ScreenWidth  int32 `json:"screenWidth"`
	ScreenHeight int32 `json:"screenHeight"`
}

func (g *Game) InitSettings() {
	log.Printf("init settings")
	if _, err := os.Open(SETTINGS_PATH); err != nil {

		err := g.CreateSettingsFirstLoad()
		if err != nil {
			log.Fatal(fmt.Sprintf("unable to init settings: %v", err))
		}
		fmt.Println("File does not exist or you don't have permission to access it")
	}

	err := g.LoadSettingsFromDisk()
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to load settings: %v", err))
	}

	fmt.Println("File exists")

}

func (g *Game) CreateSettingsFirstLoad() error {
	log.Printf("?")
	settings := &Settings{
		ScreenWidth:  800,
		ScreenHeight: 600,
	}

	g.screenWidth = settings.ScreenWidth
	g.screenHeight = settings.ScreenHeight

	return g.WriteSettingsToDisk()
}

func (g *Game) WriteSettingsToDisk() error {

	settings := &Settings{
		ScreenWidth:  g.screenWidth,
		ScreenHeight: g.screenHeight,
	}

	jsonSettings, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return os.WriteFile("settings.json", jsonSettings, 0644)

}

func (g *Game) LoadSettingsFromDisk() error {
	settingsJSON, err := os.ReadFile("settings.json")
	if err != nil {
		return err
	}
	var settings Settings
	err = json.Unmarshal(settingsJSON, &settings)
	if err != nil {
		return err
	}
	g.screenWidth = settings.ScreenWidth
	g.screenHeight = settings.ScreenHeight
	return nil

}

func SaveButtonOnClick(g *Game) {
	log.Printf("asds")
	g.ActivateScene("Board")

}

func DrawSettings(g *Game) {
	button := g.Button("Save", 100, 100, SaveButtonOnClick)
	g.DrawButton(button)
	if g.WasButtonClicked(&button) {
		button.OnClick(g)

	}
}

func UpdateSettings(g *Game) {
}
