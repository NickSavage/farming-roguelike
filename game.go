package main

import (
	"errors"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"math/rand"
	"nsavage/farming-roguelike/engine"
)

const YEARS int = 8
const ROUNDS int = YEARS * 4

func (g *Game) GetRun() interface{} {
	return g.Run
}

func (g *Game) GetScenes() map[string]*engine.Scene {
	return g.Scenes
}

func (g *Game) InitRun(loadSave bool) {

	run := &Run{
		Game:                  g,
		Money:                 100,
		MoneyRequirementStart: 200,
		MoneyRequirementRate:  2,
		Productivity:          1.0,
		CurrentRound:          1,
		CurrentYear:           1,
		CurrentSeason:         Spring,
		NextSeason:            Summer,
		Technology:            make([]*Technology, 0),
		Events:                []Event{{BlankEvent: true}},
		EventTracker:          make(map[string]bool),
		Products:              make(map[ProductType]*Product),
		ActionsMaximum:        5,
		ActionsRemaining:      5,
		AutoSellRoundEnd:      false,
	}

	g.Run = run
	g.InitTechSpaces()
	events, err := g.Run.InitEvents()
	if err != nil {
		log.Fatal(err)
	}

	g.Run.PossibleEvents = events
	g.Run.CurrentSeeds = []*Technology{
		g.CreateWheatTech(),
	}
	if loadSave {
		save, err := LoadRun()
		if err == nil {
			run.Money = save.Money
			run.Yield = save.Yield
			run.Productivity = save.Productivity
			run.CurrentRound = save.CurrentRound
			run.CurrentYear = save.CurrentYear
			run.CurrentSeason = save.CurrentSeason
			run.ActionsRemaining = save.ActionsRemaining
			run.ActionsMaximum = save.ActionsMaximum
			run.Products = save.Products
			run.Technology = g.UnpackTechnology(save.Technology)
			run.CurrentSeeds = g.UnpackSeeds(save.CurrentSeeds)
			run.Events = g.Run.UnpackEvents(save.Events)
			// todo: place tech
		}

	}

	g.Run.MoneyRequirement = g.Run.calculateMoneyRequirement()

	g.Run.CurrentRoundShopBuildings = g.ShopRandomBuildings(5)
	g.InitShopRoundComponents()
	g.ActiveRun = true
}

func (g *Game) InitTechSpaces() {
	scene := g.Scenes["Board"]
	spaces := []*TechnologySpace{
		{ // 5
			Game:           g,
			ID:             0,
			TechnologyType: Field,
			Row:            7,
			Column:         1,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 11,
				Down:  6,
				Up:    10,
			},
		},
		{ // 6
			Game:           g,
			ID:             1,
			TechnologyType: Field,
			Row:            7,
			Column:         7,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 11,
				Down:  7,
				Up:    5,
			},
		},
		{ // 7
			Game:           g,
			ID:             2,
			TechnologyType: Field,
			Row:            7,
			Column:         13,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 11,
				Down:  5,
				Up:    6,
			},
		},
		{ // 8
			Game:           g,
			ID:             3,
			TechnologyType: Field,
			Row:            19,
			Column:         1,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         false,
			SelectDirections: engine.SelectDirections{
				Left:  14,
				Right: 8,
				Down:  9,
				Up:    10,
			},
		},
		{ // 9
			Game:           g,
			ID:             4,
			TechnologyType: Field,
			Row:            19,
			Column:         7,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         false,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 9,
				Down:  10,
				Up:    8,
			},
		},
		{ // 10
			Game:           g,
			ID:             5,
			TechnologyType: Field,
			Row:            19,
			Column:         13,
			Width:          5,
			Height:         5,
			IsFilled:       false,
			Active:         false,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 10,
				Down:  8,
				Up:    9,
			},
		},
		{ // 11
			Game:           g,
			ID:             6,
			TechnologyType: BuildingSpace,
			Row:            13,
			Column:         1,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 14,
				Down:  12,
				Up:    13,
			},
		},
		{ // 12
			Game:           g,
			ID:             7,
			TechnologyType: BuildingSpace,
			Row:            13,
			Column:         4,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 14,
				Down:  13,
				Up:    11,
			},
		},
		{ // 13
			Game:           g,
			ID:             8,
			TechnologyType: BuildingSpace,
			Row:            13,
			Column:         7,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  1,
				Right: 14,
				Down:  11,
				Up:    12,
			},
		},
		{ // 14
			Game:           g,
			ID:             9,
			TechnologyType: CellTowerSpace,
			Row:            16,
			Column:         2,
			Width:          2,
			Height:         2,
			IsFilled:       false,
			Active:         true,
			SelectDirections: engine.SelectDirections{
				Left:  11,
				Right: 14,
				Down:  14,
				Up:    14,
			},
		},
	}
	for _, space := range spaces {
		scene.Components = append(scene.Components, space)
	}
	g.Run.TechnologySpaces = spaces

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
		log.Printf("round end tech %v", tech)
		tech.OnRoundEnd(g, tech)
	}
	if g.Run.AutoSellRoundEnd {
		g.Run.EndRoundMoney += g.sellAllProducts()

	}
	g.Run.Money += g.Run.EndRoundMoney * g.Run.Yield
	g.Run.Money = float32(math.Round(float64(g.Run.Money)))
	g.Run.EndRoundMoney = 0
	g.Run.ActionsRemaining = g.Run.ActionsMaximum

	g.Run.CurrentRound += 1
	g.Run.CurrentSeason.Next()
	g.Run.NextSeason.Next()
	g.GetNextEvents()

	// end of year stuff
	if g.Run.CurrentSeason == Spring {
		g.Run.CurrentYear += 1
		if g.CheckGameOver() {
			g.GameOver = true
			g.GameOverTriggered = true
			g.EndGame()
		}

		// if game isn't over, increment this
		g.Run.MoneyRequirement = g.Run.calculateMoneyRequirement()
		g.Run.CurrentRoundShopBuildings = g.ShopRandomBuildings(5)
		g.InitShopRoundComponents()
	}
	g.Run.SaveRun()
	g.WriteSettingsToDisk()

}

func (g *Game) CheckGameOver() bool {
	if g.Run.Money < g.Run.MoneyRequirement {
		return true
	}
	if g.Run.CurrentRound > ROUNDS {
		return true
	}
	return false
}

func (r *Run) calculateMoneyRequirement() float32 {
	return float32(float64(r.MoneyRequirementStart) * math.Pow(float64(r.MoneyRequirementRate), float64(r.CurrentYear)-1))
}

func (g *Game) EndGame() {

}

func (g *Game) GetNextEvents() {
	choices := 2
	g.Run.EventChoices = g.PickEventChoices(choices)
	window := g.Scenes["Board"].Windows["NextEvent"]
	components := make([]engine.UIComponent, 0)

	blank := engine.NewBlankComponent()
	blank.SelectDirections.Left = choices
	blank.SelectDirections.Right = 1
	components = append(components, &blank)
	var x float32
	for i, event := range g.Run.EventChoices {

		x = float32(240 + (i * 300))
		rect := rl.Rectangle{X: x, Y: 60, Width: 300, Height: 400}
		button := g.NewEventButton(rect, &event)
		button.SelectDirections.Left = i + 1 - 1
		button.SelectDirections.Right = i + 1 + 1
		components = append(components, &button)
	}
	window.Components = components
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

func (g *Game) ConsumeOrBuyProduct(product *Product, maximumInput float32) float32 {

	var input float32
	var market float32 = 0
	if product.Quantity > maximumInput {
		product.Quantity -= input

	} else {
		market = (maximumInput - product.Quantity) * product.Price
		product.Quantity = 0
		g.Run.Money -= market
	}
	return market
}

func (g *Game) SellProduct(product *Product) float32 {
	result := +product.Quantity * product.Price
	log.Printf("selling %v %v = %v", product.Quantity, string(product.Type), result)
	g.RecordProductStat(product, product.Quantity, result)
	product.TotalProduced += product.Quantity
	product.TotalEarned += result
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
	scaledF := float32(f*1.2 + 0.8*(1-f))
	return scaledF

}

func (g *Game) PickRandomTechnologies(needed int, keys []string) []*Technology {

	results := []*Technology{}
	for i := 0; i < needed && len(keys) > 0; i++ {
		index := rand.Intn(len(keys))
		selectedKey := keys[index]
		keys = append(keys[:index], keys[index+1:]...)
		results = append(results, g.Technology[selectedKey])
	}
	return results
}

// save files

func createTechSave(tech *Technology) TechnologySave {
	saved := TechnologySave{
		Name:           tech.Name,
		ReadyToHarvest: tech.ReadyToHarvest,
		ReadyToTouch:   tech.ReadyToTouch,
		TempYield:      tech.TempYield,
		SpaceID:        tech.Space.ID,
	}
	if tech.TechnologyType == Field {
		seeds := make([]TechnologySave, 0)
		for _, seed := range tech.Space.PlantedSeeds {
			seeds = append(seeds, createTechSave(seed))
		}
		saved.Seeds = seeds
	}
	return saved

}

func (r *Run) PackTechnology(techs []*Technology) []TechnologySave {
	results := []TechnologySave{}
	for _, tech := range techs {
		if tech.TechnologyType == Seed {
			// we'll handle this elsewhere
			continue
		}
		saved := createTechSave(tech)
		results = append(results, saved)
	}
	return results
}

func (g *Game) unpackTech(save TechnologySave) *Technology {

	new := *g.Technology[save.Name]
	copy := &new
	copy.ReadyToHarvest = save.ReadyToHarvest
	copy.ReadyToTouch = save.ReadyToTouch
	copy.TempYield = save.TempYield
	return copy

}

func (g *Game) UnpackSeeds(saved []TechnologySave) []*Technology {
	results := []*Technology{}
	for _, save := range saved {
		copy := g.unpackTech(save)
		results = append(results, copy)

	}
	return results
}

func (g *Game) UnpackTechnology(saved []TechnologySave) []*Technology {
	results := []*Technology{}
	for _, save := range saved {
		log.Printf("tech %v", save.Name)

		copy := g.unpackTech(save)
		for _, space := range g.Run.TechnologySpaces {
			if space.ID != save.SpaceID {
				continue
			}
			space.Technology = copy
			space.IsFilled = true
			copy.Space = space
			if copy.TechnologyType == Field {
				space.IsField = true
				seeds := make([]*Technology, 0)
				for _, saveSeed := range save.Seeds {
					seed := g.unpackTech(saveSeed)
					seeds = append(seeds, seed)

				}
				space.PlantedSeeds = seeds
			}
		}
		results = append(results, copy)
	}
	log.Printf("unpacked %v", results)
	return results

}

func (r *Run) PackEvents() []EventSave {
	results := []EventSave{}
	for _, event := range r.Events {
		save := EventSave{
			RoundIndex:  event.RoundIndex,
			Name:        event.Name,
			Description: event.Description,
			Effects:     event.Effects,
			BlankEvent:  event.BlankEvent,
			Repeatable:  event.Repeatable,
		}
		results = append(results, save)
	}
	return results
}

func (r *Run) UnpackEvents(saved []EventSave) []Event {
	results := []Event{}

	for _, save := range saved {
		new := Event{
			RoundIndex:  save.RoundIndex,
			Name:        save.Name,
			Description: save.Description,
			Effects:     save.Effects,
			BlankEvent:  save.BlankEvent,
			Repeatable:  save.Repeatable,
			OnTrigger:   r.triggerFunctions[save.Name],
		}
		results = append(results, new)
	}
	return results
}

func (r *Run) SaveRun() {

	saveFile := RunSaveFile{
		Money:                 r.Money,
		MoneyRequirementStart: r.MoneyRequirementStart,
		MoneyRequirementRate:  r.MoneyRequirementRate,
		Yield:                 r.Yield,
		Productivity:          r.Productivity,
		CurrentRound:          r.CurrentRound,
		CurrentYear:           r.CurrentYear,
		CurrentSeason:         r.CurrentSeason,
		ActionsRemaining:      r.ActionsRemaining,
		ActionsMaximum:        r.ActionsMaximum,
		EventTracker:          r.EventTracker,
		Technology:            r.PackTechnology(r.Technology),
		Products:              r.Products,
		Events:                r.PackEvents(),
		CurrentSeeds:          r.PackTechnology(r.CurrentSeeds),
	}
	SaveRun(saveFile)
}
