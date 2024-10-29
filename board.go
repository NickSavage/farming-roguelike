package main

import (
	"errors"
	//	"fmt"
	"log"
	"math/rand"

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

func (g *Game) GetSquareFromCoords(input BoardCoord) *BoardSquare {
	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	return &grid[input.Row][input.Column]

}

func (g *Game) GetOpenSpace(tech *Technology) (*TechnologySpace, error) {

	for _, space := range g.Run.TechnologySpaces {
		if space.IsFilled {
			continue
		}
		if space.TechnologyType != tech.TechnologyType {
			continue
		}
		return space, nil
	}
	return &TechnologySpace{}, errors.New("no empty space")
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

	g.Scenes["Board"].Data["HoverVector"] = BoardCoord{}
	g.Scenes["Board"].Data["HoverVectorCounter"] = 0
	//	g.InitPlaceRandomTrees(215)
	//	g.InitPlaceTech()

}

func (g *Game) drawTiles() {

	grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)
	for i := range grid {
		for j := range grid[i] {
			DrawTile(
				g.Data["GrassTile"].(Tile),
				float32(i*TILE_HEIGHT),
				float32(j*TILE_WIDTH),
			)
		}
	}
}
func (g *Game) DrawTechnologySpaces() {

	scene := g.Scenes["Board"]
	mousePosition := rl.GetMousePosition()
	var rect rl.Rectangle
	for _, space := range g.Run.TechnologySpaces {
		x := float32(space.Row * TILE_WIDTH)
		y := float32(space.Column * TILE_HEIGHT)
		width := float32(space.Width * TILE_WIDTH)
		height := float32(space.Height * TILE_HEIGHT)
		rect = rl.NewRectangle(x, y, width, height)
		rl.DrawRectangleRec(rect, rl.Blue)
		if !space.IsFilled {
			continue
		}
		mousePosition = rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, rect) {
			space.Technology.Tile.Color = rl.Green
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				result := g.HandleClickTech(space.Technology)
				message := Message{
					Text:  result,
					Vec:   rl.Vector2{X: x, Y: y},
					Timer: 30,
				}
				scene.Messages = append(scene.Messages, message)
			}

		} else {
			if space.Technology.ReadyToHarvest {
				space.Technology.Tile.Color = rl.Blue
			} else if space.Technology.ReadyToTouch {
				space.Technology.Tile.Color = rl.Red
			} else {
				space.Technology.Tile.Color = rl.White
			}
		}

		if space.Technology.TileFillSpace {
			for i := range space.Width {
				for j := range space.Height {
					DrawTile(
						space.Technology.Tile,
						float32(float32(x)+float32(i*TILE_WIDTH)),
						float32(float32(y)+float32(j*TILE_WIDTH)),
					)
				}
			}

		} else {
			DrawTile(space.Technology.Tile, float32(x), float32(y))

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
	//maxY := float32(g.screenHeight - 100)
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
	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      0,
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
	var results []Message
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

func DrawBoard(g *Game) {

	//	scene := g.Scenes["Board"]
	//	rl.BeginMode2D(g.Scenes["Board"].Camera)
	g.drawTiles()
	g.DrawTechnologySpaces()
	g.drawGrid()
	//	rl.EndMode2D()

	g.DrawMessages()
	g.HandleHover()
	DrawHUD(g)
	// g.DrawContextMenu(g.Scenes["Board"])
}

func (g *Game) HandleHover() {

	scene := g.Scenes["Board"]
	mousePosition := rl.GetMousePosition()

	oldVec := scene.Data["HoverVector"].(BoardCoord)
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
			if space.Technology != nil {
				g.DrawTechHoverWindow(
					*space.Technology,
					mousePosition.X,
					mousePosition.Y,
				)

			}

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

func UpdateBoard(g *Game) {
}
