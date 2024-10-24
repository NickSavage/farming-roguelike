package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

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
	rl.ImageResize(image, 90, 90)

	g.Data["ChickenCoopTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  90,
			Height: 90,
		},
		Color: rl.White,
	}
	rl.ImageResize(image, 45, 45)
	g.Data["ChickenCoopShopTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  45,
			Height: 45,
		},
		Color: rl.White,
	}

	image = rl.LoadImage("assets/trees.png")
	rl.ImageResize(image, 432, 240)

	g.Data["TreeTile"] = Tile{

		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      54,
			Y:      0,
			Width:  90,
			Height: 90,
		},
		Color: rl.White,
	}

	image = rl.LoadImage("assets/plants.png")
	rl.ImageResize(image, 225, 675)
	g.Data["WheatTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      357,
			Width:  45,
			Height: 45,
		},
		Color: rl.White,
	}
	g.Data["PotatoTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      92,
			Width:  45,
			Height: 45,
		},
		Color: rl.White,
	}

	image = rl.LoadImage("assets/tech/workstation.png")
	rl.ImageResize(image, 45, 45)
	g.Data["WorkstationTile"] = Tile{
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

	g.Scenes["HUD"] = &Scene{
		Active:      true,
		AutoDisable: false,
		DrawScene:   DrawHUD,
		UpdateScene: UpdateHUD,
		Data:        make(map[string]interface{}),
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
		Counter:      0,
	}

	g.Data["Message"] = ""

	rl.InitWindow(g.screenWidth, g.screenHeight, "Farming Roguelike")

	g.LoadAssets()
	g.LoadScenes()

	g.InitTechnology()
	g.InitRun()
	g.InitShopWindow()

	g.InitBoard()
	rl.SetTargetFPS(60)

	//	g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["Prices"])

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
