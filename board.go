package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
	"math/rand"
)

const TILE_HEIGHT = 45
const TILE_WIDTH = 45

const TILE_ROWS = 30
const TILE_COLUMNS = 30

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

func (g *Game) InitPlaceTech() {

	scene := g.Scenes["Board"]
	scene.Data["PlaceTechCancelButton"] = ShopButton{
		Width:           200,
		Height:          40,
		Title:           "Cancel Placement",
		OnClick:         OnClickCancelTechPlacement,
		BackgroundColor: rl.SkyBlue,
	}

}
func (g *Game) InitBoard() {
	scene := g.Scenes["Board"]
	scene.Camera = rl.Camera2D{}
	scene.Camera.Zoom = 1.0
	scene.Camera.Target = rl.Vector2{X: 0, Y: 0}

	rows := 30
	cols := 30

	scene.Data["Rows"] = rows
	scene.Data["Columns"] = cols

	grid := make([][]BoardSquare, rows)
	for i := range grid {
		grid[i] = make([]BoardSquare, cols)
	}
	for i := 0; i < int(rows); i++ {
		for j := 0; j < int(cols); j++ {
			square := BoardSquare{
				Tile:     g.Data["GrassTile"].(Tile),
				TileType: "Grass",
				Row:      i,
				Column:   j,
				Width:    1,
				Height:   1,
				Skip:     false,
				Occupied: false,
			}
			grid[i][j] = square
		}
	}
	g.Scenes["Board"].Data["Grid"] = grid

	g.Scenes["Board"].Data["HoverVector"] = rl.Vector2{}
	g.Scenes["Board"].Data["HoverVectorCounter"] = 0
	g.InitPlaceRandomTrees(40)
	g.InitPlaceTech()
	g.InitDrawTechnology()

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

	for i := range grid {
		for j := range grid[i] {
			tile := grid[i][j]
			if tile.Skip {
				// if these match, it is the top left of a multicell tile
				// so we don't want to skip
				if !(tile.Row == i && tile.Column == j) {
					continue
				}
			}
			if tile.HoverActive && (tile.IsTechnology || tile.IsTree) {
				tile.Tile.Color = rl.Green
			} else {
				tile.Tile.Color = rl.White
			}
			DrawTile(
				g.Data["GrassTile"].(Tile),
				float32(i*TILE_HEIGHT),
				float32(j*TILE_WIDTH),
			)

			DrawTile(
				tile.Tile,
				float32(i*TILE_HEIGHT),
				float32(j*TILE_WIDTH),
			)
		}
	}

}

func (g *Game) InitDrawTechnology() {
	for _, tech := range g.Run.Technology {
		g.DrawTechnology(tech)
	}

}
func (g *Game) RedrawTechnology() {
	for _, tech := range g.Run.Technology {
		if !tech.Redraw {
			continue
		}
		g.DrawTechnology(tech)
	}
}

func (g *Game) DrawTechnology(tech *Technology) {
	grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)

	tile := tech.Tile
	for x := range tile.Width {
		for y := range tile.Height {
			tile.Occupied = true
			tile.Technology = tech
			if tile.MultiSquare {
				tile.Skip = true
			}
			grid[tile.Row+x][tile.Column+y] = tile
		}
	}
	grid[tile.Row][tile.Column] = tile

}

func (g *Game) RemoveTechnology(square *BoardSquare) {

	grid := g.Scenes["Board"].Data["Grid"].([][]BoardSquare)

	tech := square.Technology
	for x := range square.Width {
		for y := range square.Height {
			square.Occupied = false
			square.Skip = false
			square.Technology = tech
			if square.MultiSquare {
				square.Skip = false
			}
			square.Tile = g.Data["GrassTile"].(Tile)
			square.TileType = "Grass"
			grid[square.Row+x][square.Column+y] = *square
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

func (g *Game) CheckTilesOccupied(newBoardSquare BoardSquare, mouseX, mouseY float32) bool {
	scene := g.Scenes["Board"]

	row := int((mouseX + TILE_WIDTH/2) / TILE_WIDTH)
	col := int((mouseY + TILE_HEIGHT/2) / TILE_HEIGHT)

	for x := range newBoardSquare.Width {
		for y := range newBoardSquare.Height {
			testX := row + x - 1
			if testX < 0 {
				testX = 0
			}
			if testX >= TILE_COLUMNS {
				testX = TILE_COLUMNS - 1
			}
			testY := col + y - 1
			if testY < 0 {
				testY = 0
			}
			if testY >= TILE_ROWS {
				testY = TILE_ROWS - 1
			}
			if scene.Data["Grid"].([][]BoardSquare)[testX][testY].Occupied {
				return true
			}
		}
	}
	return false

}

func OnClickCancelTechPlacement(g *Game) {
	scene := g.Scenes["Board"]
	scene.Data["PlaceTech"] = false

}
func (g *Game) DrawPlaceTech() {
	scene := g.Scenes["Board"]
	if scene.Data["PlaceTech"] == nil || !scene.Data["PlaceTech"].(bool) {
		return
	}

	if g.ScreenSkip {
		g.ScreenSkip = false
		return
	}
	chosenTech := scene.Data["PlaceChosenTech"].(*Technology)
	mousePosition := rl.GetMousePosition()

	cancelButton := scene.Data["PlaceTechCancelButton"].(ShopButton)
	g.DrawShopButton(cancelButton, 200, 50)

	if rl.CheckCollisionPointRec(mousePosition, rl.Rectangle{
		X:      200,
		Y:      50,
		Width:  float32(cancelButton.Width),
		Height: float32(cancelButton.Height),
	}) {
		// don't display placement if you're over the cancel button
		return
	}
	if g.CheckTilesOccupied(chosenTech.Tile, mousePosition.X, mousePosition.Y) {
		occupiedTile := chosenTech.Tile.Tile
		occupiedTile.Color = rl.Red
		DrawTile(
			occupiedTile,
			float32(mousePosition.X)-(chosenTech.Tile.Tile.TileFrame.Width/2),
			float32(mousePosition.Y)-(chosenTech.Tile.Tile.TileFrame.Height/2),
		)

	} else {
		DrawTile(
			chosenTech.Tile.Tile,
			float32(mousePosition.X)-(chosenTech.Tile.Tile.TileFrame.Width/2),
			float32(mousePosition.Y)-(chosenTech.Tile.Tile.TileFrame.Height/2),
		)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			scene.Data["PlaceTech"] = false
			g.PlaceTech(
				chosenTech,
				float32(mousePosition.X)-(chosenTech.Tile.Tile.TileFrame.Width/2),
				float32(mousePosition.Y)-(chosenTech.Tile.Tile.TileFrame.Height/2),
			)
		}

	}

}

func DrawBoard(g *Game) {

	//	scene := g.Scenes["Board"]
	rl.BeginMode2D(g.Scenes["Board"].Camera)
	g.drawTiles()
	g.DrawPlaceTech()
	g.drawGrid()
	rl.EndMode2D()

	g.HandleHover()
	DrawHUD(g)
	g.RedrawTechnology()
	g.DrawClickMenu()
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
				} else if tile.TileType == "Dirt" {
					tile.Tile = g.Data["GrassTile"].(Tile)
					tile.TileType = "Grass"

				}

			}
		}

	}

}

func DrawGenericMenu(g *Game) {
	scene := g.Scenes["Board"]
	//	square := scene.Menu.BoardSquare

	rl.DrawRectangleRec(scene.Menu.Rectangle, rl.White)
	rl.DrawRectangleLinesEx(scene.Menu.Rectangle, 2, rl.Black)
}

func (g *Game) HandleLeftClick() {

	scene := g.Scenes["Board"]

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePosition := rl.GetMousePosition()

		grid := scene.Data["Grid"].([][]BoardSquare)
		x := int((mousePosition.X + scene.Camera.Target.X) / scene.Camera.Zoom / float32(TILE_WIDTH))
		y := int((mousePosition.Y + scene.Camera.Target.Y) / scene.Camera.Zoom / float32(TILE_HEIGHT))
		if scene.RenderMenu {
			// TODO build a rect of the menu and check if the click is within
			if scene.Menu.BoardSquare.Row == grid[x][y].Row ||
				scene.Menu.BoardSquare.Column == grid[x][y].Column {
				return
			}
		}
		if !(grid[x][y].IsTechnology || grid[x][y].IsTree) {

			scene.RenderMenu = false
			return
		}
		menu := &BoardRightClickMenu{
			Rectangle: rl.Rectangle{
				X:      mousePosition.X,
				Y:      mousePosition.Y,
				Height: 100,
				Width:  100,
			},
			BoardSquare: &grid[x][y],
			Items:       make([]BoardMenuItem, 0),
		}

		if grid[x][y].IsTree {
			menu.Items = TreeMenuItems()
		}
		g.ScreenSkip = true
		scene.RenderMenu = true
		scene.Menu = menu

		return
	}
	if scene.RenderMenu {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

			mousePosition := rl.GetMousePosition()
			if !rl.CheckCollisionPointRec(mousePosition, scene.Menu.Rectangle) {
				scene.RenderMenu = false
			}

		}
	}

}

func (g *Game) DrawClickMenu() {
	scene := g.Scenes["Board"]
	if !scene.RenderMenu {
		return
	}

	mousePosition := rl.GetMousePosition()

	var color rl.Color

	x := scene.Menu.Rectangle.X
	y := scene.Menu.Rectangle.Y

	for _, item := range scene.Menu.Items {
		rec := item.Rectangle
		rec.X = x
		rec.Y = y
		if rl.CheckCollisionPointRec(mousePosition, rec) {
			color = rl.Gray
		} else {
			color = rl.White
		}
		rl.DrawRectangleRec(rec, color)
		rl.DrawText(item.Text, int32(rec.X+5), int32(rec.Y+5), 15, rl.Black)

		y = rec.Y + rec.Height

		if g.ScreenSkip {
			if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
				g.ScreenSkip = false
			}
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !g.ScreenSkip {
			mousePosition := rl.GetMousePosition()
			if !rl.CheckCollisionPointRec(mousePosition, item.Rectangle) {
				item.OnClick(g)
				scene.RenderMenu = false
			}

		}
	}

}

func (g *Game) disableTechHoverHighlight(vec rl.Vector2) {

	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	square := &grid[int(vec.X)][int(vec.Y)]

	square.HoverActive = false
	// if !square.IsTechnology || !square.MultiSquare {
	// 	return
	// }
	if square.Height <= 1 || square.Width <= 1 {
		return
	}
	startX := square.Row
	startY := square.Column
	for x := range square.Width {
		for y := range square.Height {
			grid[startX+x][startY+y].HoverActive = false
		}
	}
}

func (g *Game) enableTechHoverHighlight(vec rl.Vector2) {
	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	square := &grid[int(vec.X)][int(vec.Y)]

	square.HoverActive = true
	if square.Height <= 1 || square.Width <= 1 {
		return
	}
	// if !square.IsTechnology || !square.MultiSquare {
	// 	return
	// }
	startX := square.Row
	startY := square.Column
	for x := range square.Width {
		for y := range square.Height {
			grid[startX+x][startY+y].HoverActive = true
		}
	}
}

func (g *Game) HandleHover() {

	scene := g.Scenes["Board"]
	//	grid := scene.Data["Grid"].([][]BoardSquare)
	mousePosition := rl.GetMousePosition()

	oldVec := scene.Data["HoverVector"].(rl.Vector2)
	g.disableTechHoverHighlight(oldVec)

	newVec := rl.Vector2{
		X: (mousePosition.X + scene.Camera.Target.X) / scene.Camera.Zoom / float32(TILE_WIDTH),
		Y: (mousePosition.Y + scene.Camera.Target.Y) / scene.Camera.Zoom / float32(TILE_HEIGHT),
	}

	if newVec.X < 0 || newVec.X > TILE_ROWS-1 {
		return
	}
	if newVec.Y < 0 || newVec.Y > TILE_COLUMNS-1 {
		return
	}
	g.enableTechHoverHighlight(newVec)
	if oldVec.X == newVec.X && oldVec.Y == newVec.Y {
		counter := scene.Data["HoverVectorCounter"].(int)
		if counter == 0 {
			square := scene.Data["Grid"].([][]BoardSquare)[int(newVec.X)][int(newVec.Y)]

			if square.Technology != nil {
				g.DrawTechHoverWindow(
					*square.Technology,
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
		scene.Data["HoverVector"] = newVec
		scene.Data["HoverVectorCounter"] = 10
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
		if scene.Camera.Target.X > TILE_COLUMNS*TILE_WIDTH-float32(g.screenWidth-50) {
			scene.Camera.Target.X = TILE_COLUMNS*TILE_WIDTH - float32(g.screenWidth-50)
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		scene.Camera.Target.X -= 5
		if scene.Camera.Target.X < -200 {
			scene.Camera.Target.X = -200
		}
	}
	if rl.IsKeyDown(rl.KeyDown) {
		scene.Camera.Target.Y += 5
		if scene.Camera.Target.Y > TILE_ROWS*TILE_HEIGHT-float32(g.screenHeight-200) {
			scene.Camera.Target.Y = TILE_ROWS*TILE_HEIGHT - float32(g.screenHeight-200)
		}
	}
	if rl.IsKeyDown(rl.KeyUp) {
		scene.Camera.Target.Y -= 5
		if scene.Camera.Target.Y < -300 {
			scene.Camera.Target.Y = -300
		}
	}
	g.SelectTiles()
	g.HandleLeftClick()

	// mousePosition := rl.GetMousePosition()
}
