package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const TILE_HEIGHT = 45
const TILE_WIDTH = 45

type BoardSquare struct {
	Tile     Tile
	TileType string
}

func (g *Game) InitBoard() {
	scene := g.Scenes["Board"]
	scene.Camera = rl.Camera2D{}
	scene.Camera.Zoom = 1.0

	rows := g.screenWidth / TILE_HEIGHT
	cols := g.screenHeight / TILE_WIDTH

	grid := make([][]BoardSquare, rows)
	for i := range grid {
		grid[i] = make([]BoardSquare, cols)
	}
	for i := 0; i < int(rows); i++ {
		for j := 0; j < int(cols); j++ {
			square := BoardSquare{Tile: g.Data["DirtTile"].(Tile), TileType: "Dirt"}
			grid[i][j] = square
		}
	}
	g.Scenes["Board"].Data["Grid"] = grid

}

func (g *Game) drawTiles() {

	grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)
	for i := range grid {
		for j := range grid[i] {
			DrawTile(grid[i][j].Tile, float32(i*TILE_HEIGHT), float32(j*TILE_WIDTH))
		}
	}

}

func (g *Game) drawGrid() {
	var spacing float32
	spacing = 45
	var x float32
	x = 0
	//	y := 0
	startVec := rl.Vector2{X: 0, Y: 0}
	endVec := rl.Vector2{X: 0, Y: float32(g.screenHeight)}
	for {
		startVec.X = x
		endVec.X = x
		rl.DrawLineEx(startVec, endVec, 2, rl.Black)

		x += spacing
		if x >= float32(g.screenWidth) {
			break
		}

	}
	var y float32
	y = 0
	startVec = rl.Vector2{X: 0, Y: 0}
	endVec = rl.Vector2{X: float32(g.screenWidth), Y: 0}
	for {
		startVec.Y = y
		endVec.Y = y
		rl.DrawLineEx(startVec, endVec, 2, rl.Black)

		y += spacing
		if y >= float32(g.screenHeight) {
			break
		}

	}

}

func DrawBoard(g *Game) {
	rl.DrawText("hello baord world", 0, 0, 10, rl.Black)
	rl.BeginMode2D(g.Scenes["Board"].Camera)
	g.drawTiles()
	g.drawGrid()
	rl.EndMode2D()
}

func UpdateBoard(g *Game) {
	scene := g.Scenes["Board"]
	// Camera zoom controls
	scene.Camera.Zoom += float32(rl.GetMouseWheelMove()) * 0.05

	if scene.Camera.Zoom > 3.0 {
		scene.Camera.Zoom = 3.0
	} else if scene.Camera.Zoom < 0.1 {
		scene.Camera.Zoom = 0.1
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePosition := rl.GetMousePosition()
		x := int(mousePosition.X / float32(TILE_WIDTH))
		y := int(mousePosition.Y / float32(TILE_HEIGHT))

		grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)
		if grid[x][y].TileType == "Grass" {
			grid[x][y].Tile = g.Data["DirtTile"].(Tile)
			grid[x][y].TileType = "Dirt"

		} else {
			grid[x][y].Tile = g.Data["GrassTile"].(Tile)
			grid[x][y].TileType = "Grass"

		}
	}

}
