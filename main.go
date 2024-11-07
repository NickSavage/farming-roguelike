package main

import (
	"log"
	"nsavage/farming-roguelike/engine"
	"os"

	"github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) LoadAssets() {
	image := rl.LoadImage("assets/grass.png")
	rl.ImageResize(image, 264, 168)
	g.Data["GrassTile"] = LoadImage(image, 0, 120, 1, 1)

	image = rl.LoadImage("assets/dirt.png")

	rl.ImageResize(image, 264, 168)
	g.Data["DirtTile"] = LoadImage(image, 22, 120, 1, 1)

	image = rl.LoadImage("assets/tech/chicken_coop.png")

	rl.ImageResize(image, 90, 90)
	g.Data["ChickenCoopTile"] = LoadImage(image, 0, 0, 2, 2)
	rl.ImageResize(image, 45, 45)
	g.Data["ChickenCoopShopTile"] = LoadImage(image, 0, 0, 1, 1)

	image = rl.LoadImage("assets/trees.png")
	rl.ImageResize(image, 432, 240)

	g.Data["TreeTile"] = LoadImage(image, 54, 0, 2, 2)

	image = rl.LoadImage("assets/plants.png")
	rl.ImageResize(image, 225, 675)
	g.Data["WheatTile"] = LoadImage(image, 0, 357, 1, 1)
	g.Data["PotatoTile"] = LoadImage(image, 0, 92, 1, 1)

	image = rl.LoadImage("assets/tech/workstation.png")
	rl.ImageResize(image, 45, 45)
	g.Data["WorkstationTile"] = LoadImage(image, 0, 0, 1, 1)

	image = rl.LoadImage("assets/tech/chicken_egg_warmer.png")
	//	rl.ImageResize(image, 45, 45)
	g.Data["ChickenEggWarmerTile"] = LoadImage(image, 0, 0, 2, 2)
	rl.ImageResize(image, 45, 45)
	g.Data["ChickenEggWarmerShopTile"] = LoadImage(image, 0, 0, 1, 1)

	image = rl.LoadImage("assets/icons/wheat.png")
	g.Data["WheatIcon"] = LoadImage(image, 0, 0, 2, 2)
	image = rl.LoadImage("assets/icons/potato.png")
	g.Data["PotatoIcon"] = LoadImage(image, 0, 0, 2, 2)
	image = rl.LoadImage("assets/icons/carrots.png")
	g.Data["CarrotShopIcon"] = LoadImage(image, 0, 0, 2, 2)
	rl.ImageResize(image, 45, 45)
	g.Data["CarrotIcon"] = LoadImage(image, 0, 0, 1, 1)

	image = rl.LoadImage("assets/icons/apples.png")
	g.Data["AppleShopIcon"] = LoadImage(image, 0, 0, 2, 2)
	rl.ImageResize(image, 45, 45)
	g.Data["AppleIcon"] = LoadImage(image, 0, 0, 1, 1)

	image = rl.LoadImage("assets/icons/cell_tower.png")
	g.Data["CellTowerTile"] = LoadImage(image, 0, 0, 2, 2)

}

func (g *Game) LoadScenes() {
	g.Scenes["GameMenu"] = &engine.Scene{
		Active:      true,
		AutoDisable: true,
		DrawScene:   DrawGameMenu,
		UpdateScene: UpdateGameMenu,
		Data:        make(map[string]interface{}),
		KeyBindings: make(map[string]*engine.KeyBinding),
	}
	g.Scenes["Board"] = &engine.Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawBoard,
		UpdateScene: UpdateBoard,
		Data:        make(map[string]interface{}),
		KeyBindings: make(map[string]*engine.KeyBinding),
		Components:  make([]engine.UIComponent, 0),
	}

	// g.Scenes["HUD"] = &Scene{
	// 	Active:      false,
	// 	AutoDisable: true,
	// 	DrawScene:   DrawHUD,
	// 	UpdateScene: UpdateHUD,
	// 	Data:        make(map[string]interface{}),
	// 	Buttons:     make([]Button, 1),
	// 	KeyBindings: make(map[string]*KeyBinding),
	// }
	g.Scenes["Settings"] = &engine.Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawSettings,
		UpdateScene: UpdateSettings,
		Data:        make(map[string]interface{}),
		KeyBindings: make(map[string]*engine.KeyBinding),
	}
	g.InitHUD()

}

func main() {
	log.Printf("hello world")
	g := Game{
		Scenes:  map[string]*engine.Scene{},
		Data:    make(map[string]interface{}),
		Counter: 0,
	}
	g.InitSettings()

	g.Data["Message"] = ""

	file, err := os.Open("./save.json")
	if err != nil {
		g.ExistingSave = false
	} else {
		g.ExistingSave = true
	}
	file.Close()

	rl.InitWindow(g.screenWidth, g.screenHeight, "Farming Roguelike")

	g.LoadAssets()
	g.LoadScenes()

	g.InitTechnology()
	g.InitBoard()
	g.InitShopWindow()
	g.InitGameMenu()
	g.InitSettingsMenu()

	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	// dev shit

	g.InitRun(false)
	g.ActivateScene("Board")

	for !rl.WindowShouldClose() {
		g.Counter += 1
		if g.Counter == 60 {
			g.Seconds += 1
			g.Counter = 0
		}

		g.Draw()
		g.Update()

	}
	rl.CloseWindow()
}
