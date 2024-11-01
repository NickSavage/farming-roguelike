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

	events, err := g.InitEvents()
	if err != nil {
		log.Fatal(err)
	}

	g.Run = &Run{
		Money:            100,
		Productivity:     1.0,
		CurrentRound:     1,
		CurrentSeason:    Spring,
		NextSeason:       Summer,
		Technology:       make([]*Technology, 0),
		People:           make([]Person, 1),
		PossibleEvents:   events,
		Events:           []Event{{BlankEvent: true}},
		Products:         make(map[ProductType]*Product),
		ActionsMaximum:   5,
		ActionsRemaining: 5,
	}
	g.InitTechSpaces()

	g.Run.CurrentRoundShopPlants = g.ShopRandomPlants(2)
}

func (g *Game) InitTechSpaces() {
	spaces := []*TechnologySpace{
		{
			TechnologyType: PlantSpace,
			Row:            7,
			Column:         1,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         true,
		},
		{
			TechnologyType: PlantSpace,
			Row:            7,
			Column:         7,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         true,
		},
		{
			TechnologyType: PlantSpace,
			Row:            7,
			Column:         13,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         true,
		},
		{
			TechnologyType: PlantSpace,
			Row:            19,
			Column:         1,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         false,
		},
		{
			TechnologyType: PlantSpace,
			Row:            19,
			Column:         7,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         false,
		},
		{
			TechnologyType: PlantSpace,
			Row:            19,
			Column:         13,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         false,
		},
		{
			TechnologyType: BuildingSpace,
			Row:            13,
			Column:         1,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
		},
		{
			TechnologyType: BuildingSpace,
			Row:            13,
			Column:         4,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
		},
		{
			TechnologyType: BuildingSpace,
			Row:            13,
			Column:         7,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
		},
		{
			TechnologyType: CellTowerSpace,
			Row:            16,
			Column:         2,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
		},
	}
	g.Run.TechnologySpaces = spaces

	scene := g.Scenes["Board"]
	grid := scene.Data["Grid"].([][]BoardSquare)
	for _, space := range g.Run.TechnologySpaces {
		if !space.Active {
			continue
		}
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
		tech.OnRoundEnd(g, tech)
	}
	g.Run.EndRoundMoney += g.sellAllProducts()
	g.Run.Money += g.Run.EndRoundMoney * g.Run.Yield
	g.Run.Money = float32(math.Round(float64(g.Run.Money)))
	g.Run.EndRoundMoney = 0
	g.Run.ActionsRemaining = g.Run.ActionsMaximum

	g.Run.CurrentRound += 1
	g.Run.CurrentSeason.Next()
	g.Run.NextSeason.Next()
	g.GetNextEvents()

	g.Run.CurrentRoundShopPlants = g.ShopRandomPlants(2)

}

func (g *Game) GetNextEvents() {

	g.Run.EventChoices = g.PickEventChoices(2)
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

func (r *Run) CanSpendAction(actions int) bool {
	if r.ActionsRemaining >= actions {
		return true
	}
	return false
}

func (r *Run) SpendAction(actions int) error {

	if r.CanSpendAction(actions) {
		r.ActionsRemaining -= actions
		return nil
	}
	return errors.New("cannot spend action, not enough actions")
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
	log.Printf("selling %v %v = %v", product.Quantity, string(product.Type), result)
	product.Quantity = 0
	// TODO: add to round money for reporting?
	product.TotalEarned += result

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
	scaledF := float32(f*1.2 + 0.8*(1-f))
	return scaledF

}

func (g *Game) ShopRandomPlants(needed int) []*Technology {
	// Filter technologies with Type == PlantSpace and collect their keys into a slice
	plantSpaceTechnologies := make([]*Technology, 0)
	for _, tech := range g.Technology {
		if tech.TechnologyType == PlantSpace {
			plantSpaceTechnologies = append(plantSpaceTechnologies, tech)
		}
	}

	// Select a slice of unique random technologies from the filtered list
	result := []*Technology{}
	keysToPickFrom := make([]string, 0)
	for key, _ := range g.Technology {
		if _, found := findPlantSpaceTech(g.Technology, key); found {
			keysToPickFrom = append(keysToPickFrom, key)
		}
	}

	for i := 0; i < needed && len(keysToPickFrom) > 0; i++ {
		index := rand.Intn(len(keysToPickFrom))
		selectedKey := keysToPickFrom[index]
		keysToPickFrom = append(keysToPickFrom[:index], keysToPickFrom[index+1:]...)
		result = append(result, g.Technology[selectedKey])
	}

	return result
}

func findPlantSpaceTech(techMap map[string]*Technology, key string) (bool, bool) {
	_, found := techMap[key]
	if found && techMap[key].TechnologyType == PlantSpace {
		return true, true
	}
	return false, false
}
