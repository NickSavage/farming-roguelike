package main

import (
	"log"
	"math/rand"
	"nsavage/farming-roguelike/engine"

	"github.com/gen2brain/raylib-go/raylib"
)

const TILE_HEIGHT = 45
const TILE_WIDTH = 45

const TILE_ROWS = 25
const TILE_COLUMNS = 25

func CheckVecVisible(vec rl.Vector2) bool {
	// todo this is the sidebar, how can I do this better
	if vec.X < 200 {
		return false
	}
	return true
}

func (g *Game) GetSquareFromCoords(input engine.BoardCoord) *BoardSquare {
	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	return &grid[input.Row][input.Column]

}

func (g *Game) GetVecFromCoords(input engine.BoardCoord) rl.Vector2 {
	return rl.Vector2{
		X: float32(input.Row*TILE_WIDTH + int(g.SidebarWidth)),
		Y: float32(input.Column * TILE_HEIGHT),
	}
}

func generateCoordinates(numPairs, maxX, maxY int) []rl.Vector2 {
	coordinates := make([]rl.Vector2, numPairs)

	for i := 0; i < numPairs; i++ {
		coordinates[i] = rl.Vector2{
			X: float32(rand.Intn(maxX)),
			Y: float32(rand.Intn(maxY)),
		}
	}

	return coordinates
}
func (g *Game) InitPlaceRandomTrees(numTrees int) {
	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	tile := g.Data["TreeTile"].(Tile)

	coords := generateCoordinates(numTrees, TILE_ROWS/2, TILE_COLUMNS/2)

	for _, coord := range coords {
		x := int(coord.X * 2)
		y := int(coord.Y * 2)
		boardSquare := grid[x][y]
		boardSquare.Tile = tile
		boardSquare.TileType = "Tree"
		boardSquare.Width = 2
		boardSquare.Height = 2
		boardSquare.Occupied = true
		boardSquare.MultiSquare = true
		boardSquare.IsTree = true

		for i := range boardSquare.Width {
			for j := range boardSquare.Height {
				grid[x+i][y+j] = boardSquare
				grid[x+i][y+j].Skip = true
				grid[x+i][y+j].Occupied = true
			}
		}
		grid[x][y] = boardSquare
	}

}

func (g *Game) GetGrassSquare(x, y int) *BoardSquare {

	square := &BoardSquare{
		Tile:     g.Data["GrassTile"].(Tile),
		TileType: "Grass",
		Row:      x,
		Column:   y,
		Width:    1,
		Height:   1,
		Skip:     false,
		Occupied: false,
	}
	return square
}

func (g *Game) InitBoard() {

	log.Printf("init")
	scene := g.Scenes["Board"]
	scene.Camera = rl.Camera2D{}
	scene.Camera.Zoom = 1
	scene.Camera.Target = rl.Vector2{X: -float32(g.SidebarWidth), Y: 0}

	rows := TILE_ROWS
	cols := TILE_COLUMNS

	scene.Data["Rows"] = rows
	scene.Data["Columns"] = cols

	grid := make([][]BoardSquare, rows)
	for i := range grid {
		grid[i] = make([]BoardSquare, cols)
	}
	for i := 0; i < int(rows); i++ {
		for j := 0; j < int(cols); j++ {
			grid[i][j] = *g.GetGrassSquare(i, j)
		}
	}
	g.Scenes["Board"].Data["Grid"] = grid

	g.Scenes["Board"].Data["HoverVector"] = engine.BoardCoord{}
	g.Scenes["Board"].Data["HoverVectorCounter"] = 0

	g.Scenes["Board"].Data["PreviousMousePosition"] = rl.Vector2{X: 0, Y: 0}

	// scene.KeyBindingFunctions = make(map[string]func(engine.GameInterface))
}

func (g *Game) drawTiles() {

	grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)
	for i := range grid {
		for j := range grid[i] {
			if i == 0 && j == 0 {
				DrawTile(
					g.Data["DirtTile"].(Tile),
					float32(i*TILE_WIDTH+int(g.SidebarWidth)),
					float32(j*TILE_HEIGHT),
				)

			} else {
				DrawTile(
					g.Data["GrassTile"].(Tile),
					float32(i*TILE_WIDTH+int(g.SidebarWidth)),
					float32(j*TILE_HEIGHT),
				)

			}
		}
	}
}
func (g *Game) DrawTechnologySpaces() {

	// scene := g.Scenes["Board"]
	// mousePosition := rl.GetMousePosition()
	// var rect rl.Rectangle
	for _, space := range g.Run.TechnologySpaces {
		if !space.Active {
			continue
		}
		space.Render()
	}
}

func (g *Game) drawGrid() {
	var spacing float32
	spacing = 45
	var x float32
	columns := g.Scenes["Board"].Data["Columns"].(int)
	rows := g.Scenes["Board"].Data["Rows"].(int)

	maxX := float32(columns)*spacing + float32(g.SidebarWidth)
	//maxY := float32(g.screenHeight - 100)
	maxY := float32(columns) * spacing
	x = float32(g.SidebarWidth)
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
	startVec = rl.Vector2{X: float32(g.SidebarWidth), Y: 0}
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
	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      float32(g.SidebarWidth),
		Y:      0,
		Width:  maxX,
		Height: maxY,
	},
		5,
		rl.Black,
	)

}

func (g *Game) DrawMessages() {
	scene := g.Scenes["Board"]
	var results []engine.Message
	for _, message := range scene.Messages {
		rl.DrawText(message.Text, int32(message.Vec.X+50), int32(message.Vec.Y), 25, rl.Black)
		message.Timer -= 1
		message.Vec.Y -= 1
		if message.Timer != 0 {
			results = append(results, message)
		}
	}
	scene.Messages = results

}

// main draw function

func DrawBoard(gi engine.GameInterface) {
	g := gi.(*Game)

	//	scene := g.Scenes["Board"]
	//	rl.BeginMode2D(g.Scenes["Board"].Camera)
	g.drawTiles()
	g.drawGrid()
	//	rl.EndMode2D()

	//	g.HandleHover()
	DrawHUD(g)
	g.DrawMessages()
	// g.DrawContextMenu(g.Scenes["Board"])
}

func (g *Game) HandleHover() {

	scene := g.Scenes["Board"]
	mousePosition := rl.GetMousePosition()

	oldVec := scene.Data["HoverVector"].(engine.BoardCoord)
	coords := g.GetBoardCoordAtPoint(mousePosition)
	if coords.Row < 0 || coords.Row > TILE_ROWS-1 {
		return
	}
	if coords.Column < 0 || coords.Column > TILE_COLUMNS-1 {
		return
	}
	//	square := g.GetSquareFromCoords(coords)
	//	log.Printf("square %v,%v - %v", coords.Row, coords.Column, square)
	if oldVec.Row == coords.Row && oldVec.Column == coords.Column {
		counter := scene.Data["HoverVectorCounter"].(int)
		if counter == 0 {

			square := g.GetSquareFromCoords(coords)
			if !square.IsTechnologySpace {
				return
			}

			space := square.TechnologySpace
			if !space.IsFilled {
				return
			}
			// if space.Technology != nil {
			// 	g.DrawTechHoverWindow(
			// 		space.Technology,
			// 		mousePosition.X,
			// 		mousePosition.Y,
			// 	)

			// }

		} else if scene.RenderMenu {
			counter = 10

		} else {
			counter = counter - 1
			scene.Data["HoverVectorCounter"] = counter
		}

	} else {
		scene.Data["HoverVector"] = coords
		scene.Data["HoverVectorCounter"] = 10
	}

}

func UpdateBoard(gi engine.GameInterface) {
	g := gi.(*Game)
	// scene := g.Scenes["Board"]
	UpdateHUD(g)

	// handle selecting
	// mousePosition := rl.GetMousePosition()
	// old := scene.Data["PreviousMousePosition"].(rl.Vector2)
	// if old.X != mousePosition.X && old.Y != mousePosition.Y {
	// 	g.MouseMode = true
	// } else {
	// 	g.MouseMode = false
	// }

	// if g.MouseMode {
	// 	g.Data["PreviousMousePosition"] = mousePosition
	// 	for _, component := range scene.Components {
	// 		if rl.CheckCollisionPointRec(mousePosition, component.Rect()) {
	// 			component.Select()
	// 		} else {
	// 			component.Unselect()

	// 		}
	// 	}
	// }

}

// cursor

func MoveCursorLeft(gi engine.GameInterface)  {}
func MoveCursorRight(gi engine.GameInterface) {}
func MoveCursorUp(gi engine.GameInterface)    {}
func MoveCursorDown(gi engine.GameInterface)  {}
