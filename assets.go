package main

import rl "github.com/gen2brain/raylib-go/raylib"

func LoadImage(image *rl.Image, x, y float32, widthTiles, heightTiles int) Tile {
	return Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      x,
			Y:      y,
			Width:  float32(widthTiles * TILE_WIDTH),
			Height: float32(heightTiles * TILE_HEIGHT),
		},
		Color: rl.White,
	}

}
