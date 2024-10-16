package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

type Technology struct {
	Name        string
	Description string
	Tile        BoardSquare
	OnRoundEnd  func(*Game)
	OnBuild     func(*Game)
}

type Person struct {
}

type Run struct {
	Technology            []*Technology
	People                []Person
	Money                 float32
	Productivity          float32
	EndRoundMoney         float32
	RoundActions          int
	RoundActionsRemaining int
}

func (g *Game) InitRun() {

	g.Run = &Run{
		Money:                 100,
		Productivity:          1.0,
		RoundActions:          5,
		RoundActionsRemaining: 5,
		Technology:            make([]*Technology, 0),
		People:                make([]Person, 1),
	}
	g.Run.Technology = append(g.Run.Technology, g.CreateChickenCoopTech())
	g.Run.Technology = append(g.Run.Technology, g.CreateCornTech())

}
func OnClickEndRound(g *Game) {
	g.Run.EndRoundMoney += 100
	for _, tech := range g.Run.Technology {
		tech.OnRoundEnd(g)
	}
	g.Run.Money += g.Run.EndRoundMoney * g.Run.Productivity
	g.Run.EndRoundMoney = 0
}

func (g *Game) PlaceTech(tech *Technology, x, y float32) {

	row := int((x + TILE_WIDTH/2) / TILE_WIDTH)
	col := int((y + TILE_HEIGHT/2) / TILE_HEIGHT)

	// Store calculated row and column in tech
	tech.Tile.Row = row
	tech.Tile.Column = col

	log.Printf("tech %v", len(g.Run.Technology))
	g.Run.Technology = append(g.Run.Technology, tech)
	log.Printf("tech afte %v", len(g.Run.Technology))

}

func (g *Game) InitTechnology() {
	tech := make(map[string]Technology)
	chicken := Technology{
		Name:        "Chicken Coop",
		Tile:        BoardSquare{},
		Description: "asdasd",
		OnRoundEnd:  ChickenCoopRoundEnd,
	}
	tech["ChickenCoop"] = chicken

	tech["CornField"] = Technology{
		Name:        "Corn",
		Tile:        BoardSquare{},
		Description: "asdasd",
		OnRoundEnd:  CornFieldRoundEnd,
	}

	g.Data["Technology"] = tech
}

func (g *Game) CreateChickenCoopTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["ChickenCoop"]
	result.Tile = BoardSquare{
		Tile:     g.Data["ChickenCoopTile"].(Tile),
		TileType: "Technology",
		Row:      10,
		Column:   10,
		Width:    2,
		Height:   2,
		Occupied: true,
	}
	return &result
}

func (g *Game) CreateCornTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["CornField"]
	result.Tile = BoardSquare{
		Tile:     g.Data["CornTile"].(Tile),
		TileType: "Technology",
		Row:      8,
		Column:   8,
		Width:    5,
		Height:   10,
		Occupied: true,
	}
	return &result
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
		g.drawRunTech(*tech, 210, 90)

	}

}

func ChickenCoopOnBuild(g *Game) {
	g.Run.Productivity += 0.05

}

func ChickenCoopRoundEnd(g *Game) {
	g.Run.EndRoundMoney += 5
}

func CornFieldRoundEnd(g *Game) {
	g.Run.EndRoundMoney += 5
}
