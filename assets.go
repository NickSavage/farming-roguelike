package main

import rl "github.com/gen2brain/raylib-go/raylib"

func LoadImage(image *rl.Image, x, y, width, height float32) Tile {
	return Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Color: rl.White,
	}

}
