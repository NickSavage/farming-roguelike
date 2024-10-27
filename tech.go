package main

import (
	"fmt"
	"log"
	//	"sort"
	"errors"
	//
	// rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) InitTechnology() {
	log.Printf("init tech")
	tech := make(map[string]*Technology)

	tech["ChickenCoop"] = g.ChickenCoop()
	tech["WheatField"] = g.WheatField()
	tech["PotatoField"] = g.PotatoField()

	tech["Workstation"] = g.Workstation()
	tech["ChickenEggWarmer"] = g.ChickenEggWarmer()

	g.Technology = tech
}

func (g *Game) InitProduct(productType ProductType, price float32) {

	if _, exists := g.Run.Products[productType]; !exists {
		g.Run.Products[productType] = &Product{
			Type:     productType,
			Quantity: 0,
			Price:    price,
			Yield:    1,
		}
	}
}

func (g *Game) GetProductNames() []ProductType {
	results := []ProductType{}
	for _, product := range g.Run.Products {
		results = append(results, product.Type)
	}
	// todo sort
	//sort.Strings(results)
	return results
}

func (g *Game) PlaceTech(tech *Technology, space *TechnologySpace) error {
	space.IsFilled = true
	copy := &tech
	space.Technology = *copy
	tech.Space = space

	err := g.Run.SpendAction(tech.CostActions)
	if err != nil {
		return errors.New("cannot spend action")
	}

	err = g.Run.SpendMoney(tech.CostMoney)

	if err == nil {
		err := tech.OnBuild(g, tech)
		if err == nil {
			g.Run.Technology = append(g.Run.Technology, tech)
		}
	}
	return nil
}

func (g *Game) RemoveTech(tech *Technology) {

	tech.ToBeDeleted = true
	space := tech.Space
	space.IsFilled = false
	space.Technology = &Technology{}

	var results []*Technology
	for _, tech := range g.Run.Technology {
		if !tech.ToBeDeleted {
			results = append(results, tech)
		}
	}
	g.Run.Technology = results
}

func (g *Game) HandleClickTech(tech *Technology) string {
	return tech.OnClick(g, tech)
}

// chicken

func (g *Game) CreateChickenCoopTech() *Technology {

	result := g.ChickenCoop()
	result.Square = BoardSquare{
		//		Tile:         g.Data["ChickenCoopTile"].(Tile),
		TileType:    "Technology",
		Row:         10,
		Column:      10,
		Width:       2,
		Height:      2,
		Occupied:    true,
		MultiSquare: true,
	}

	return result
}

func (g *Game) ChickenCoop() *Technology {
	result := &Technology{
		Name:            "Chicken Coop",
		ProductType:     Chicken,
		TechnologyType:  BuildingSpace,
		Tile:            g.Data["ChickenCoopTile"].(Tile),
		TileWidth:       1,
		TileHeight:      1,
		TileFillSpace:   false,
		Square:          BoardSquare{},
		Description:     "asdasd",
		CostMoney:       50,
		CanBuild:        ChickenCoopCanBuild,
		OnBuild:         ChickenCoopOnBuild,
		OnClick:         ChickenCoopOnClick,
		OnRoundEnd:      ChickenCoopRoundEnd,
		RoundEndProduce: ChickenCoopProduce,
		TempYield:       1,
	}
	return result
}

func ChickenCoopCanBuild(g *Game, tech *Technology) bool {
	return true
}

func ChickenCoopOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, 5)
	return nil
}

func ChickenCoopProduce(g *Game, tech *Technology) float32 {
	return 5 * g.Run.Productivity * g.Run.Products["Chicken"].Yield * tech.TempYield
}

func ChickenCoopOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		produced := ChickenCoopProduce(g, tech)
		g.Run.Products["Chicken"].Quantity += produced

		tech.ReadyToHarvest = false
		return fmt.Sprintf("Chicken: %v", produced)
	}
	return ""
}

func ChickenCoopRoundEnd(g *Game, tech *Technology) {
	tech.ReadyToHarvest = true
}

// wheat

func (g *Game) CreateWheatTech() *Technology {

	result := g.WheatField()
	result.Square = BoardSquare{
		//		Tile:         g.Data["WheatTile"].(Tile),
		TileType: "Technology",
		Row:      8,
		Column:   8,
		Width:    5,
		Height:   5,
		Occupied: true,
	}
	g.Run.Products["Wheat"] = &Product{
		Type:     Wheat,
		Quantity: 0,
		Price:    1,
		Yield:    1,
	}

	return result
}

func (g *Game) WheatField() *Technology {
	return &Technology{
		Name:            "Wheat",
		ProductType:     Wheat,
		TechnologyType:  PlantSpace,
		Tile:            g.Data["WheatTile"].(Tile),
		TileWidth:       1,
		TileHeight:      1,
		TileFillSpace:   true,
		Square:          BoardSquare{},
		CostMoney:       50,
		Description:     "asdasd",
		CanBuild:        WheatFieldCanBuild,
		OnBuild:         WheatFieldOnBuild,
		OnClick:         WheatFieldOnClick,
		OnRoundEnd:      WheatFieldRoundEnd,
		RoundEndProduce: WheatFieldProduce,
		TempYield:       1,
		ReadyToTouch:    true,
	}
}

func WheatFieldCanBuild(g *Game, tech *Technology) bool {
	if g.Run.CurrentSeason == Spring {
		return true
	}
	return false
}

func WheatFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, 1)
	return nil
}

func WheatFieldProduce(g *Game, tech *Technology) float32 {
	result := float32(125) * g.Run.Productivity * g.Run.Products["Wheat"].Yield * tech.TempYield
	log.Printf("wheat %v", result)
	return result
}

func WheatFieldRoundEnd(g *Game, tech *Technology) {
	if g.Run.NextSeason == Autumn {
		tech.ReadyToHarvest = true
	} else if g.Run.NextSeason == Winter {
		g.RemoveTech(tech)
	} else {
		tech.ReadyToTouch = true
	}
	tech.Tile.TileFrame.X += 45
}
func WheatFieldOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToTouch {
		err := g.Run.SpendAction(1)
		if err == nil {
			tech.TempYield += 0.05
			tech.ReadyToTouch = false
			return fmt.Sprintf("Yield: %v", tech.TempYield)
		}
		tech.ReadyToTouch = false

	} else if tech.ReadyToHarvest {
		produced := WheatFieldProduce(g, tech)
		g.Run.Products["Wheat"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Wheat: %v", produced)
	}
	return ""
}

// potato

func (g *Game) CreatePotatoTech() *Technology {

	result := g.PotatoField()
	result.Square = BoardSquare{
		//	Tile:         g.Data["PotatoTile"].(Tile),
		TileType: "Technology",
		Row:      8,
		Column:   8,
		Width:    4,
		Height:   4,
		Occupied: true,
	}

	g.InitProduct(result.ProductType, 5)
	return result
}

func (g *Game) PotatoField() *Technology {
	return &Technology{
		Name:            "Potato",
		ProductType:     Potato,
		TechnologyType:  PlantSpace,
		Tile:            g.Data["PotatoTile"].(Tile),
		TileWidth:       1,
		TileHeight:      1,
		TileFillSpace:   true,
		Square:          BoardSquare{},
		CostMoney:       50,
		Description:     "asdasd",
		CanBuild:        PotatoFieldCanBuild,
		OnBuild:         PotatoFieldOnBuild,
		OnClick:         PotatoFieldOnClick,
		OnRoundEnd:      PotatoFieldRoundEnd,
		RoundEndProduce: PotatoFieldProduce,
		TempYield:       1,
		ReadyToTouch:    true,
	}

}

func PotatoFieldCanBuild(g *Game, tech *Technology) bool {
	return true
}

func PotatoFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, 1)
	return nil
}

func PotatoFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return float32(125) * g.Run.Productivity * g.Run.Products["Potato"].Yield * tech.TempYield
	} else {
		return 0
	}
}
func PotatoFieldRoundEnd(g *Game, tech *Technology) {
	if g.Run.NextSeason == Autumn {
		tech.ReadyToHarvest = true
	}
	tech.Tile.TileFrame.X += 45
}

func PotatoFieldOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		produced := PotatoFieldProduce(g, tech)
		g.Run.Products["Potato"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Potatoes: %v", produced)
	}
	return ""
}

// workstation

func (g *Game) CreateWorkstationTech() *Technology {

	result := g.Workstation()
	result.Square = BoardSquare{
		//	Tile:         g.Data["WorkstationTile"].(Tile),
		TileType: "Technology",
		Row:      1,
		Column:   1,
		Width:    1,
		Height:   1,
		Occupied: true,
	}

	return result
}

func (g *Game) Workstation() *Technology {
	return &Technology{
		Name:           "Workstation",
		TechnologyType: BuildingSpace,
		Tile:           g.Data["WorkstationTile"].(Tile),
		TileWidth:      1,
		TileHeight:     1,
		TileFillSpace:  false,
		Square:         BoardSquare{},
		CostMoney:      25,
		CostActions:    1,
		Description:    "asdasd",
		CanBuild:       WorkstationCanBuild,
		OnBuild:        WorkstationOnBuild,
		OnClick:        WorkstationOnClick,
		OnRoundEnd:     WorkstationRoundEnd,
	}

}

func WorkstationCanBuild(g *Game, tech *Technology) bool {
	// if g.Run.CanSpendAction(g)
	return true
}

func WorkstationOnBuild(g *Game, tech *Technology) error {
	g.Run.Productivity += 0.05

	return nil

}
func WorkstationRoundEnd(g *Game, tech *Technology) {

}
func WorkstationOnClick(g *Game, tech *Technology) string {
	return ""
}

func (g *Game) CreateChickenEggWarmer() *Technology {

	result := g.CreateChickenEggWarmer()
	result.Square = BoardSquare{
		//	Tile:         g.Data["WorkstationTile"].(Tile),
		TileType: "Technology",
		Row:      1,
		Column:   1,
		Width:    1,
		Height:   1,
		Occupied: true,
	}

	return result
}

func (g *Game) ChickenEggWarmer() *Technology {
	return &Technology{
		Name:           "ChickenEggWarmer",
		ProductType:    Chicken,
		TechnologyType: BuildingSpace,
		Tile:           g.Data["ChickenEggWarmerTile"].(Tile),
		TileWidth:      1,
		TileHeight:     1,
		TileFillSpace:  false,
		Square:         BoardSquare{},
		CostMoney:      25,
		Description:    "asdasd",
		CanBuild:       ChickenEggWarmerCanBuild,
		OnBuild:        ChickenEggWarmerOnBuild,
		OnClick:        ChickenEggWarmerOnClick,
		OnRoundEnd:     ChickenEggWarmerRoundEnd,
	}

}

func ChickenEggWarmerCanBuild(g *Game, tech *Technology) bool {
	hasCoop := false
	for _, tech := range g.Run.Technology {
		if tech.Name == "Chicken Coop" {
			hasCoop = true
			break
		}
	}
	return hasCoop
}

func ChickenEggWarmerOnBuild(g *Game, tech *Technology) error {
	g.InitProduct("Chicken", 5)
	g.Run.Products["Chicken"].Yield += 0.05

	return nil

}
func ChickenEggWarmerRoundEnd(g *Game, tech *Technology) {

}

func ChickenEggWarmerOnClick(g *Game, tech *Technology) string {
	return ""
}
