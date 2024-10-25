package main

import (
	"errors"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"math/rand"
)

const YEARS int = 8

func (g *Game) InitRun() {

	g.Run = &Run{
		Money:        100,
		Productivity: 1.0,
		CurrentRound: 1,
		Technology:   make([]*Technology, 0),
		People:       make([]Person, 1),
		Events:       GenerateRandomEvents(),
		Products:     make(map[string]*Product),
	}
	g.InitTechSpaces()
	//	g.Run.Technology = append(g.Run.Technology, g.CreateChickenCoopTech())
	// g.Run.Technology = append(g.Run.Technology, g.CreateWheatTech())

}

func (g *Game) InitTechSpaces() {
	spaces := []*TechnologySpace{
		{
			TechnologyType: PlantSpace,
			Row:            2,
			Column:         1,
			Width:          5,
			Height:         5,
			IsFilled:       false,
		},
		{
			TechnologyType: PlantSpace,
			Row:            2,
			Column:         7,
			Width:          5,
			Height:         5,
			IsFilled:       false,
		},
		{
			TechnologyType: PlantSpace,
			Row:            2,
			Column:         13,
			Width:          5,
			Height:         5,
			IsFilled:       false,
		},
		{
			TechnologyType: BuildingSpace,
			Row:            9,
			Column:         1,
			Width:          2,
			Height:         2,
			IsFilled:       false,
		},
		{
			TechnologyType: BuildingSpace,
			Row:            9,
			Column:         4,
			Width:          2,
			Height:         2,
			IsFilled:       false,
		},
		{
			TechnologyType: BuildingSpace,
			Row:            9,
			Column:         7,
			Width:          2,
			Height:         2,
			IsFilled:       false,
		},
	}
	g.Run.TechnologySpaces = spaces

	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	for _, space := range g.Run.TechnologySpaces {
		for x := range space.Width {
			for y := range space.Height {
				grid[space.Row+x][space.Column+y].IsTechnologySpace = true
				grid[space.Row+x][space.Column+y].TechnologySpace = space
			}
		}

	}

}

func PreEndRound(g *Game) {
	g.Run.Yield = g.Run.GenerateYield()
}

func OnClickEndRound(g *Game) {

	for _, tech := range g.Run.Technology {
		tech.RoundHandler[tech.RoundHandlerIndex].OnRoundEnd(g, tech)
		g.Run.EndRoundMoney -= tech.RoundHandler[tech.RoundHandlerIndex].CostMoney
	}
	g.Run.EndRoundMoney += g.sellAllProducts()
	g.Run.Money += g.Run.EndRoundMoney * g.Run.Yield
	g.Run.Money = float32(math.Round(float64(g.Run.Money)))
	g.Run.EndRoundMoney = 0

	g.Run.CurrentRound += 1
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
	for _, effect := range event.Effects {
		if effect.IsPriceChange {
			log.Printf("Price of %v change by %v", effect.ProductImpacted, effect.PriceChange)
			current := g.Run.Products[effect.ProductImpacted].Price
			g.Run.Products[effect.ProductImpacted].Price = current * (1 + effect.PriceChange)
		}
	}
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

func (g *Game) drawExistingTechIcon(tech Technology, x, y float32) {

	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  tech.Tile.TileFrame.Width,
		Height: tech.Tile.TileFrame.Height,
	}
	DrawTile(tech.Square.Tile, x, y)

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
		g.drawExistingTechIcon(*tech, float32(210+(i*offset)), 90)

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

func (r *Run) CanSpendMoney(money float32) bool {
	if r.Money >= money {
		return true
	}
	return false
}

func (r *Run) SpendMoney(money float32) error {

	if r.CanSpendMoney(money) {
		r.Money -= money
		return nil
	}
	return errors.New("cannot spend money, not enough money")
}

func (g *Game) sellAllProducts() float32 {
	var result float32 = 0
	for _, product := range g.Run.Products {
		result += +g.SellProduct(product)
	}

	log.Printf("sell %v", result)
	return result
}

func (g *Game) SellProduct(product *Product) float32 {
	result := +product.Quantity * product.Price
	log.Printf("selling %v %v = %v", product.Quantity, product.Name, result)
	product.Quantity = 0
	// TODO: add to round money for reporting?

	return result
}

func (r *Run) CalculateNetWorth() float32 {
	var result float32 = 0
	for _, product := range r.Products {
		value := product.Price * product.Quantity
		result += value
	}
	result += r.Money
	return result
}

// yield

func (r *Run) GenerateYield() float32 {
	f := rand.Float64()
	scaledF := float32(f*1.2 + 0.8)
	return scaledF

}
