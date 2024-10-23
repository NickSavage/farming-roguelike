package main

import (
	"fmt"
	"log"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) InitTechnology() {
	tech := make(map[string]Technology)

	tech["ChickenCoop"] = g.ChickenCoop()
	tech["WheatField"] = g.WheatField()
	tech["Workstation"] = g.Workstation()

	g.Data["Technology"] = tech
}

func (g *Game) GetProductNames() []string {
	results := []string{}
	for _, product := range g.Run.Products {
		results = append(results, product.Name)
	}
	sort.Strings(results)
	return results
}

func (g *Game) CreateChickenCoopTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["ChickenCoop"]
	result.Tile = BoardSquare{
		Tile:         g.Data["ChickenCoopTile"].(Tile),
		TileType:     "Technology",
		Row:          10,
		Column:       10,
		Width:        2,
		Height:       2,
		Occupied:     true,
		MultiSquare:  true,
		IsTechnology: true,
	}
	g.Run.Products["Chicken"] = &Product{
		Name:     "Chicken",
		Quantity: 0,
		Price:    5,
	}

	return &result
}

func (g *Game) CreateWheatTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["WheatField"]
	result.Tile = BoardSquare{
		Tile:         g.Data["WheatTile"].(Tile),
		TileType:     "Technology",
		Row:          8,
		Column:       8,
		Width:        5,
		Height:       5,
		Occupied:     true,
		IsTechnology: true,
	}
	g.Run.Products["Wheat"] = &Product{
		Name:     "Wheat",
		Quantity: 0,
		Price:    1,
	}

	return &result
}

func (g *Game) ChickenCoop() Technology {
	result := Technology{
		Name:        "Chicken Coop",
		Tile:        BoardSquare{},
		Description: "asdasd",
		Cost:        50,
		OnBuild:     ChickenCoopOnBuild,
		Redraw:      false,
		RoundHandler: []TechnologyRoundHandler{
			{
				OnRoundEnd:    ChickenCoopRoundEnd,
				RoundEndValue: ChickenCoopRoundEndValue,
				RoundEndText:  ChickenCoopRoundEndText,
			},
		},
		RoundHandlerIndex: 0,
	}
	return result
}

func ChickenCoopCanBeBuilt(g *Game) bool {
	return true
}

func ChickenCoopOnBuild(g *Game, tech *Technology) {
	g.Run.Productivity += 0.05
	g.Run.RoundActionsRemaining -= 1
	g.Run.Money -= tech.Cost
}

func ChickenCoopRoundEndText(g *Game, tech *Technology) string {
	units := ChickenCoopProduce(g, tech)
	price := g.Run.Products["Chicken"].Price
	text := "$%v (%v units at $%v each)"
	return fmt.Sprintf(text, units*price, units, price)
}

func ChickenCoopRoundEndValue(g *Game, tech *Technology) float32 {
	return ChickenCoopProduce(g, tech) //* g.Run.Products["Chicken"].Price
}

func ChickenCoopProduce(g *Game, tech *Technology) float32 {
	return 5 * g.Run.Productivity

}

func ChickenCoopRoundEnd(g *Game, tech *Technology) {
	g.Run.Products["Chicken"].Quantity += ChickenCoopProduce(g, tech)
	log.Printf("chicken %v", g.Run.Products["Chicken"].Quantity)
}

func (g *Game) WheatField() Technology {
	return Technology{
		Name:        "Wheat",
		Tile:        BoardSquare{},
		Cost:        50,
		Description: "asdasd",
		OnBuild:     WheatFieldOnBuild,
		Redraw:      false,
		RoundHandler: []TechnologyRoundHandler{
			{
				Season:        Spring,
				OnRoundEnd:    WheatFieldRoundSpring,
				RoundEndValue: WheatFieldRoundEndValue,
				RoundEndText:  WheatFieldRoundEndText,
				CostActions:   1,
			},
			{
				Season:        Summer,
				OnRoundEnd:    WheatFieldRoundSummer,
				RoundEndValue: WheatFieldRoundEndValue,
				RoundEndText:  WheatFieldRoundEndText,
				CostActions:   1,
			},
			{
				Season:        Autumn,
				OnRoundEnd:    WheatFieldRoundAutumn,
				RoundEndValue: WheatFieldRoundEndValue,
				RoundEndText:  WheatFieldRoundEndText,
				CostActions:   1,
			},
			{
				Season:        Winter,
				OnRoundEnd:    WheatFieldRoundWinter,
				RoundEndValue: WheatFieldRoundEndValue,
				RoundEndText:  WheatFieldRoundEndText,
				CostActions:   1,
			},
		},
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
	}

}

func WheatFieldCanBeBuilt(g *Game) bool {
	return true
}

func WheatFieldOnBuild(g *Game, tech *Technology) {

	g.Run.Products["Wheat"] = &Product{
		Name:     "Wheat",
		Quantity: 0,
		Price:    1,
	}
	g.Run.RoundActionsRemaining -= 1
	g.Run.Money -= tech.Cost
}

func WheatFieldRoundEndValue(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return 125
	} else {
		return 0
	}
}
func WheatFieldRoundEndText(g *Game, tech *Technology) string {
	units := WheatFieldProduce(g, tech)
	price := g.Run.Products["Wheat"].Price
	return fmt.Sprintf(
		"$%v (%v units at $%v each)",
		units*price,
		units,
		price,
	)
}

func WheatFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return float32(125) * g.Run.Productivity
	} else {
		return 0
	}
}

func WheatFieldRoundSpring(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 1
	tech.Tile.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundSummer(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 1
	tech.Tile.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundAutumn(g *Game, tech *Technology) {
	g.Run.Products["Wheat"].Quantity += 125
	tech.RoundHandlerIndex += 1
	tech.Tile.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundWinter(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 0
	tech.Tile.Tile.TileFrame.X += 45
	tech.Redraw = true
}

// workstation

func (g *Game) CreateWorkstationTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["Workstation"]
	result.Tile = BoardSquare{
		Tile:         g.Data["WorkstationTile"].(Tile),
		TileType:     "Technology",
		Row:          1,
		Column:       1,
		Width:        1,
		Height:       1,
		Occupied:     true,
		IsTechnology: true,
	}

	return &result
}

func (g *Game) Workstation() Technology {
	return Technology{
		Name:        "Workstation",
		Tile:        BoardSquare{},
		Cost:        25,
		Description: "asdasd",
		OnBuild:     WorkstationOnBuild,
		Redraw:      false,
		RoundHandler: []TechnologyRoundHandler{
			{
				OnRoundEnd:    WorkstationRoundEnd,
				RoundEndValue: WorkstationRoundEndValue,
				RoundEndText:  WorkstationRoundEndText,
				CostActions:   0,
			},
		},
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
	}

}

func WorkstationOnBuild(g *Game, tech *Technology) {
	g.Run.Productivity += 0.05
	g.Run.Money -= tech.Cost

}
func WorkstationRoundEnd(g *Game, tech *Technology) {

}
func WorkstationRoundEndValue(g *Game, tech *Technology) float32 {
	return 0

}
func WorkstationRoundEndText(g *Game, tech *Technology) string {
	return ""

}

// trees

func TreeMenuItems() []BoardMenuItem {
	results := []BoardMenuItem{}
	results = append(results, BoardMenuItem{
		Rectangle: rl.Rectangle{
			X:      0,
			Y:      0,
			Height: 30,
			Width:  150,
		},
		Text:            "Chop (1 action)",
		OnClick:         ChopTree,
		CheckIsDisabled: IsChopActionDisabled,
	})
	results = append(results, BoardMenuItem{
		Rectangle: rl.Rectangle{
			X:      0,
			Y:      0,
			Height: 30,
			Width:  150,
		},
		Text:            "Test",
		OnClick:         BlankAction,
		CheckIsDisabled: IsBlankActionDisabled,
	})

	return results
}

func IsChopActionDisabled(g *Game, square *BoardSquare) bool {
	if !g.Run.CanSpendAction(1) {
		return true
	}
	return false

}

func ChopTree(g *Game) {

	scene := g.Scenes["Board"]
	log.Printf("square %v", scene.Menu.BoardSquare)
	err := g.Run.SpendAction(1)
	if err == nil {
		g.RemoveTechnology(scene.Menu.BoardSquare)
	}
}

// generic actions

func BlankAction(g *Game) {}

func IsBlankActionDisabled(g *Game, square *BoardSquare) bool {
	return false
}
