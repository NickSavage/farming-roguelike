package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

type Technology struct {
	Name        string
	Description string
	Tile        BoardSquare
}

type Person struct {
}

type Run struct {
	Technology []Technology
	People     []Person
	Money      int32
}

func (g *Game) InitRun() {

	g.Run = &Run{
		Money:      100,
		Technology: make([]Technology, 0),
		People:     make([]Person, 1),
	}
	g.Run.Technology = append(g.Run.Technology, g.CreateChickenCoopTech())

}

func (g *Game) EndRound() {
	g.Run.Money += 100
}

func (g *Game) PlaceTech(tech *Technology, x, y float32) {

	row := int((x + TILE_WIDTH/2) / TILE_WIDTH)
	col := int((y + TILE_HEIGHT/2) / TILE_HEIGHT)

	// Store calculated row and column in tech
	tech.Tile.Row = row
	tech.Tile.Column = col

	log.Printf("tech %v", len(g.Run.Technology))
	g.Run.Technology = append(g.Run.Technology, *tech)
	log.Printf("tech afte %v", len(g.Run.Technology))

}

func (g *Game) CreateChickenCoopTech() Technology {
	result := Technology{
		Name: "Chicken Coop",
		Tile: BoardSquare{},
	}
	result.Tile = BoardSquare{
		Tile:     g.Data["ChickenCoopTile"].(Tile),
		TileType: "Technology",
		Row:      10,
		Column:   10,
		Width:    2,
		Height:   2,
		Occupied: true,
	}
	return result
}

func (g *Game) drawRunTech(tech Technology, x, y float32) {

	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  tech.Tile.Tile.TileFrame.Width,
		Height: tech.Tile.Tile.TileFrame.Height,
	}
	DrawTile(tech.Tile.Tile, x, y)

	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, rect) {
		toolTipRect := rl.Rectangle{
			X:      x + 30,
			Y:      y + 30,
			Width:  200,
			Height: 100,
		}
		rl.DrawRectangleRec(toolTipRect, rl.White)
		rl.DrawRectangleLinesEx(toolTipRect, 1, rl.Black)
		rl.DrawText("tooltip", int32(x+35), int32(y+35), 20, rl.Black)
	}

}

func (g *Game) DrawTechnologyWindow() {
	rl.DrawRectangle(200, 50, 900, 500, rl.White)

	rl.DrawText("Technology", 205, 55, 30, rl.Black)

	for _, tech := range g.Run.Technology {
		g.drawRunTech(tech, 210, 90)

	}

}
