package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

const TILE_HEIGHT = 45
const TILE_WIDTH = 45

type BoardSquare struct {
	Tile     Tile
	TileType string
	Row      int
	Column   int
}

func (g *Game) InitBoard() {
	scene := g.Scenes["Board"]
	scene.Camera = rl.Camera2D{}
	scene.Camera.Zoom = 1.0
	scene.Camera.Target = rl.Vector2{X: 100, Y: 100}

	rows := 100
	cols := 100

	scene.Data["Rows"] = rows
	scene.Data["Columns"] = cols

	grid := make([][]BoardSquare, rows)
	for i := range grid {
		grid[i] = make([]BoardSquare, cols)
	}
	for i := 0; i < int(rows); i++ {
		for j := 0; j < int(cols); j++ {
			square := BoardSquare{
				Tile:     g.Data["DirtTile"].(Tile),
				TileType: "Dirt",
				Row:      i,
				Column:   j,
			}
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
	columns := g.Scenes["Board"].Data["Columns"].(int)
	rows := g.Scenes["Board"].Data["Rows"].(int)

	maxX := float32(columns) * spacing
	maxY := float32(columns) * spacing
	x = 0
	//	y := 0
	startVec := rl.Vector2{X: 0, Y: 0}
	endVec := rl.Vector2{X: 0, Y: maxY}
	for {
		startVec.X = x
		endVec.X = x
		rl.DrawLineEx(startVec, endVec, 2, rl.Black)

		x += spacing
		columns -= 1
		if columns <= 0 {
			break
		}

	}
	var y float32
	y = 0
	startVec = rl.Vector2{X: 0, Y: 0}
	endVec = rl.Vector2{X: maxX, Y: 0}
	for {
		startVec.Y = y
		endVec.Y = y
		rl.DrawLineEx(startVec, endVec, 2, rl.Black)

		y += spacing
		rows -= 1
		if rows <= 0 {
			break
		}

	}

}

func DrawBoard(g *Game) {

	scene := g.Scenes["Board"]
	rl.BeginMode2D(g.Scenes["Board"].Camera)
	g.drawTiles()
	g.drawGrid()

	rl.EndMode2D()

	currentGesture := rl.GetGestureDetected()
	if currentGesture == rl.GestureNone {
		scene.Data["DragSelectionStart"] = nil
		scene.Data["DragSelectionStop"] = nil

	}
	log.Printf("gesture %v", currentGesture)
	if currentGesture == rl.GestureHold {
		scene.Data["DragSelectionStart"] = rl.GetMousePosition()
	}
	if currentGesture == rl.GestureDrag {
		log.Printf("???")
		current := rl.GetMousePosition()
		scene.Data["DragSelectionStop"] = rl.GetMousePosition()
		startVec := scene.Data["DragSelectionStart"].(rl.Vector2)
		width := current.X - startVec.X
		height := current.Y - startVec.Y
		log.Printf("draw (%v/%v) (%v/%v)", startVec.X, startVec.Y, width, height)
		rl.DrawRectangleLines(
			int32(startVec.X),
			int32(startVec.Y),
			int32(width),
			int32(height),
			rl.Black,
		)
	}
}

func (g *Game) SelectTiles() {

	scene := g.Scenes["Board"]

	currentGesture := rl.GetGestureDetected()
	if currentGesture == rl.GestureNone {
		if scene.Data["DragSelectionStart"] == nil {
			return
		}
		if scene.Data["DragSelectionStop"] == nil {
			return
		}
		start := scene.Data["DragSelectionStart"].(rl.Vector2)
		end := scene.Data["DragSelectionStop"].(rl.Vector2)
		xTiles := int((end.X-start.X)/TILE_WIDTH) + 1
		yTiles := int((end.Y-start.Y)/TILE_HEIGHT) + 1
		startX := int((start.X + scene.Camera.Target.X) / scene.Camera.Zoom / float32(TILE_WIDTH))
		startY := int((start.Y + scene.Camera.Target.Y) / scene.Camera.Zoom / float32(TILE_HEIGHT))

		for xOffset := range xTiles {
			for yOffset := range yTiles {
				grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)
				tile := &grid[startX+xOffset][startY+yOffset]
				if tile.TileType == "Grass" {
					tile.Tile = g.Data["DirtTile"].(Tile)
					tile.TileType = "Dirt"
				} else {
					tile.Tile = g.Data["GrassTile"].(Tile)
					tile.TileType = "Grass"

				}

			}
		}

	}

}

func UpdateBoard(g *Game) {
	scene := g.Scenes["Board"]

	//	Camera zoom controls
	scene.Camera.Zoom += float32(rl.GetMouseWheelMove()) * 0.05
	if scene.Camera.Zoom > 1.2 {
		scene.Camera.Zoom = 1.2
	} else if scene.Camera.Zoom < 0.8 {
		scene.Camera.Zoom = 0.8
	}
	if rl.IsKeyDown(rl.KeyRight) {
		scene.Camera.Target.X += 5
		if scene.Camera.Target.X > 1000 {
			scene.Camera.Target.X = 1000
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		scene.Camera.Target.X -= 5
		if scene.Camera.Target.X < 0 {
			scene.Camera.Target.X = 0
		}
	}
	g.SelectTiles()

	// mousePosition := rl.GetMousePosition()
}
