package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
)

const YEARS int = 8

func (g *Game) InitRun() {

	g.Run = &Run{
		Money:                 100,
		Productivity:          1.0,
		CurrentRound:          1,
		RoundActions:          5,
		RoundActionsRemaining: 5,
		Technology:            make([]*Technology, 0),
		People:                make([]Person, 1),
		Events:                GenerateRandomEvents(),
		Products:              make(map[string]*Product),
	}
	g.Run.Technology = append(g.Run.Technology, g.CreateChickenCoopTech())
	g.Run.Technology = append(g.Run.Technology, g.CreateWheatTech())

}

func (r *Run) sellAllProducts() float32 {
	var result float32 = 0
	for _, product := range r.Products {
		result += +product.Quantity * product.Price
		log.Printf("selling %v %v = %v", product.Quantity, product.Name, result)

	}
	return result
}

func OnClickEndRound(g *Game) {
	g.Run.CurrentRound += 1
	g.Run.RoundActionsRemaining = g.Run.RoundActions
	for _, tech := range g.Run.Technology {
		tech.RoundHandler[tech.RoundHandlerIndex].OnRoundEnd(g, tech)
		g.Run.RoundActionsRemaining -= int(tech.RoundHandler[tech.RoundHandlerIndex].CostActions)
		g.Run.EndRoundMoney -= tech.RoundHandler[tech.RoundHandlerIndex].CostMoney
	}
	g.Run.EndRoundMoney += g.Run.sellAllProducts()
	g.Run.Money += g.Run.EndRoundMoney * g.Run.Productivity
	g.Run.Money = float32(math.Round(float64(g.Run.Money)))
	g.Run.EndRoundMoney = 0

	g.Run.CurrentSeason.Next()
	g.GetNextEvent()

}

func (g *Game) GetNextEvent() {

	newEvent := g.NewRandomEvent()
	newEvent.RoundIndex = g.Run.CurrentRound
	g.Run.Events[newEvent.RoundIndex] = newEvent
}

func (g *Game) ProcessNextEvent() {
	event := g.Run.Events[g.Run.CurrentRound]
	log.Printf("process event %v", event.Name)
	for _, effect := range event.Effects {
		if effect.IsPriceChange {
			log.Printf("Price of %v change by %v", effect.ProductImpacted, effect.PriceChange)
			current := g.Run.Products[effect.ProductImpacted].Price
			g.Run.Products[effect.ProductImpacted].Price = current * (1 + effect.PriceChange)
		}
	}
}

func (g *Game) PlaceTech(tech *Technology, x, y float32) {

	if g.Run.RoundActionsRemaining < 1 {
		g.Data["Message"] = "Unable to build Technology, out of actions"
		g.Data["MessageCounter"] = g.Seconds + 5
		return
	}

	row := int((x + TILE_WIDTH/2) / TILE_WIDTH)
	col := int((y + TILE_HEIGHT/2) / TILE_HEIGHT)

	// Store calculated row and column in tech
	tech.Tile.Row = row
	tech.Tile.Column = col

	log.Printf("tech %v", len(g.Run.Technology))
	g.Run.Technology = append(g.Run.Technology, tech)
	tech.OnBuild(g, tech)
	log.Printf("tech afte %v", len(g.Run.Technology))
	g.DrawTechnology(tech)

}

func (g *Game) DrawTechHoverWindow(tech Technology, x, y float32) {
	toolTipRect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  200,
		Height: 100,
	}
	rl.DrawRectangleRec(toolTipRect, rl.White)
	rl.DrawRectangleLinesEx(toolTipRect, 1, rl.Black)
	rl.DrawText(tech.Name, int32(x+5), int32(y+5), 20, rl.Black)
	rl.DrawText(tech.Description, int32(x+5), int32(y+25), 10, rl.Black)

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
		g.DrawTechHoverWindow(tech, x+30, y+30)
	}

}

func DrawTechnologyWindow(g *Game, win *Window) {
	windowWidth := 900
	offset := 90
	rl.DrawRectangle(200, 50, int32(windowWidth), 500, rl.White)

	rl.DrawText("Technology", 205, 55, 30, rl.Black)

	for i, tech := range g.Run.Technology {
		g.drawRunTech(*tech, float32(210+(i*offset)), 90)

	}

}

func GenerateRandomEvents() []Event {
	results := []Event{}
	for i := range 4 * YEARS {
		results = append(results, Event{
			RoundIndex: i,
			Name:       fmt.Sprintf("Event %v", i),
		})

	}

	return results

}
