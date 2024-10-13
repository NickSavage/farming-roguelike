package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Scenes       map[string]*Scene
	Data         map[string]interface{}
	screenWidth  int32
	screenHeight int32
	Run          *Run
}

func (g *Game) LoadAssets() {
	image := rl.LoadImage("assets/grass.png")
	rl.ImageResize(image, 264, 168)

	g.Data["GrassTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      120,
			Width:  45,
			Height: 45,
		},
		Color: rl.White,
	}
	image = rl.LoadImage("assets/dirt.png")
	rl.ImageResize(image, 264, 168)

	g.Data["DirtTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      22,
			Y:      120,
			Width:  45,
			Height: 45,
		},
		Color: rl.White,
	}
	image = rl.LoadImage("assets/tech/chicken_coop.png")
	rl.ImageResize(image, 45, 45)

	g.Data["ChickenCoopTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  45,
			Height: 45,
		},
		Color: rl.White,
	}
}

func (g *Game) LoadScenes() {
	g.Scenes["Board"] = &Scene{
		Active:      true,
		AutoDisable: true,
		DrawScene:   DrawBoard,
		UpdateScene: UpdateBoard,
		Data:        make(map[string]interface{}),
	}
	g.InitBoard()

	g.Scenes["HUD"] = &Scene{
		Active:      true,
		AutoDisable: false,
		DrawScene:   DrawHUD,
		UpdateScene: UpdateHUD,
		Buttons:     make([]Button, 1),
	}
	g.InitHUD()

}

func main() {
	g := Game{
		Scenes:       map[string]*Scene{},
		Data:         make(map[string]interface{}),
		screenWidth:  int32(1280),
		screenHeight: int32(800),
	}

	rl.InitWindow(g.screenWidth, g.screenHeight, "Farming Roguelike")

	g.LoadAssets()
	g.LoadScenes()

	g.InitRun()

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {

		g.Draw()
		g.Update()

	}
	rl.CloseWindow()
}
