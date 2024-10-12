package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Scenes       map[string]*Scene
	Data         map[string]interface{}
	screenWidth  int32
	screenHeight int32
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

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {

		g.Draw()
		g.Update()

	}
	rl.CloseWindow()
}
