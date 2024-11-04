package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const SETTINGS_PATH = "settings.json"

type Settings struct {
	ScreenWidth  int32            `json:"screenWidth"`
	ScreenHeight int32            `json:"screenHeight"`
	KeyBindings  []KeyBindingJSON `json:"keyBindings"`
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

func (g *Game) InitSettingsMenu() {

	scene := g.Scenes["Settings"]

	scene.KeyBindingFunctions = make(map[string]func(*Game))
	scene.KeyBindingFunctions["CloseMenu"] = CloseMenu

	g.LoadSceneShortcuts("Settings")
	log.Printf("settings shorcuts %v", scene.KeyBindings)
}

func (g *Game) CreateSettingsFirstLoad() error {
	log.Printf("?")
	bindings := g.LoadInitialBindings()
	g.KeyBindingJSONs = bindings
	settings := &Settings{
		ScreenWidth:  800,
		ScreenHeight: 600,
		KeyBindings:  bindings,
	}

	g.screenWidth = settings.ScreenWidth
	g.screenHeight = settings.ScreenHeight

	return g.WriteSettingsToDisk()
}

func (g *Game) LoadInitialBindings() []KeyBindingJSON {
	var initialBindings []KeyBindingJSON

	file, err := os.Open("./assets/bindings.json")
	if err != nil {
		fmt.Println(err)
		return initialBindings
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)

	err = jsonDecoder.Decode(&initialBindings)
	if err != nil {
		fmt.Println(err)
		return initialBindings
	}
	return initialBindings
}

func (g *Game) WriteSettingsToDisk() error {

	settings := &Settings{
		ScreenWidth:  g.screenWidth,
		ScreenHeight: g.screenHeight,
		KeyBindings:  g.KeyBindingJSONs,
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
	g.KeyBindingJSONs = settings.KeyBindings
	return nil

}

func SaveButtonOnClick(g *Game) {
	log.Printf("asds")
	CloseMenu(g)

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

func CloseMenu(g *Game) {
	returnScene := g.Scenes["Settings"].Data["Return"].(string)
	if returnScene == "" {
		g.ActivateScene("GameMenu")
	} else {
		g.ActivateScene(returnScene)
		if returnScene == "Board" {
			g.Scenes["HUD"].Active = true
		}
	}
}
