package main

import (
	"fmt"
	"log"
	//	"sort"
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

func (g *Game) RoundEndValue(tech *Technology) float32 {
	units := tech.RoundEndProduce(g, tech)
	price := g.Run.Products[tech.ProductType].Price
	return units * price

}

func (g *Game) RoundEndText(tech *Technology) string {

	units := tech.RoundEndProduce(g, tech)
	price := g.Run.Products[tech.ProductType].Price
	text := "$%v (%v units at $%v each)"
	return fmt.Sprintf(text, units*price, units, price)
}

func (g *Game) PlaceTech(tech *Technology, space *TechnologySpace) error {
	space.IsFilled = true
	space.Technology = tech
	tech.Space = space

	if g.CanBuild(tech) {
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

// chicken

func (g *Game) CreateChickenCoopTech() *Technology {

	result := g.Technology["ChickenCoop"]
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
		Name:              "Chicken Coop",
		ProductType:       Chicken,
		TechnologyType:    BuildingSpace,
		Tile:              g.Data["ChickenCoopTile"].(Tile),
		TileWidth:         1,
		TileHeight:        1,
		TileFillSpace:     false,
		Square:            BoardSquare{},
		Description:       "asdasd",
		CostMoney:         50,
		CanBuild:          ChickenCoopCanBuild,
		OnBuild:           ChickenCoopOnBuild,
		Redraw:            false,
		OnRoundEnd:        ChickenCoopRoundEnd,
		RoundEndProduce:   ChickenCoopProduce,
		RoundHandlerIndex: 0,
		ShowEndRound:      true,
	}
	return result
}

func (g *Game) CanBuild(tech *Technology) bool {
	err := g.Run.SpendMoney(tech.CostMoney)
	if err != nil {
		log.Printf("err %v", err)
		return false
	}
	return true
}

func ChickenCoopCanBuild(g *Game) bool {
	return true
}

func ChickenCoopOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, 5)
	return nil
}

func ChickenCoopProduce(g *Game, tech *Technology) float32 {
	return 5 * g.Run.Productivity * g.Run.Products["Chicken"].Yield
}

func ChickenCoopRoundEnd(g *Game, tech *Technology) {
	g.Run.Products["Chicken"].Quantity += ChickenCoopProduce(g, tech)
	log.Printf("chicken %v", g.Run.Products["Chicken"].Quantity)
}

// wheat

func (g *Game) CreateWheatTech() *Technology {

	result := g.Technology["WheatField"]
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
	}

	return result
}

func (g *Game) WheatField() *Technology {
	return &Technology{
		Name:              "Wheat",
		ProductType:       Wheat,
		TechnologyType:    PlantSpace,
		Tile:              g.Data["WheatTile"].(Tile),
		TileWidth:         1,
		TileHeight:        1,
		TileFillSpace:     true,
		Square:            BoardSquare{},
		CostMoney:         50,
		Description:       "asdasd",
		CanBuild:          WheatFieldCanBuild,
		OnBuild:           WheatFieldOnBuild,
		Redraw:            false,
		OnRoundEnd:        WheatFieldRoundEnd,
		RoundEndProduce:   WheatFieldProduce,
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      true,
	}

}

func WheatFieldCanBuild(g *Game) bool {
	return true
}

func WheatFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, 1)
	return nil
}

func WheatFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return float32(125) * g.Run.Productivity * g.Run.Products["Wheat"].Yield
	} else {
		return 0
	}
}

func WheatFieldRoundEnd(g *Game, tech *Technology) {
	g.Run.Products["Wheat"].Quantity += WheatFieldProduce(g, tech)
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}

// potato

func (g *Game) CreatePotatoTech() *Technology {

	result := g.Technology["PotatoField"]
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
		Name:              "Potato",
		ProductType:       Potato,
		TechnologyType:    PlantSpace,
		Tile:              g.Data["PotatoTile"].(Tile),
		TileWidth:         1,
		TileHeight:        1,
		TileFillSpace:     true,
		Square:            BoardSquare{},
		CostMoney:         50,
		Description:       "asdasd",
		CanBuild:          PotatoFieldCanBuild,
		OnBuild:           PotatoFieldOnBuild,
		Redraw:            false,
		OnRoundEnd:        PotatoFieldRoundEnd,
		RoundEndProduce:   PotatoFieldProduce,
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      true,
	}

}

func PotatoFieldCanBuild(g *Game) bool {
	return true
}

func PotatoFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, 1)
	return nil
}

func PotatoFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return float32(125) * g.Run.Productivity * g.Run.Products["Potato"].Yield
	} else {
		return 0
	}
}
func PotatoFieldRoundEnd(g *Game, tech *Technology) {
	g.Run.Products["Potato"].Quantity += PotatoFieldProduce(g, tech)
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}

// workstation

func (g *Game) CreateWorkstationTech() *Technology {

	result := g.Technology["Workstation"]
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
		Name:              "Workstation",
		TechnologyType:    BuildingSpace,
		Tile:              g.Data["WorkstationTile"].(Tile),
		TileWidth:         1,
		TileHeight:        1,
		TileFillSpace:     false,
		Square:            BoardSquare{},
		CostMoney:         25,
		Description:       "asdasd",
		CanBuild:          WorkstationCanBuild,
		OnBuild:           WorkstationOnBuild,
		Redraw:            false,
		OnRoundEnd:        WorkstationRoundEnd,
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      false,
	}

}

func WorkstationCanBuild(g *Game) bool {
	return true
}

func WorkstationOnBuild(g *Game, tech *Technology) error {
	g.Run.Productivity += 0.05

	return nil

}
func WorkstationRoundEnd(g *Game, tech *Technology) {

}

func (g *Game) CreateChickenEggWarmer() *Technology {

	result := g.Technology["ChickenEggWarmer"]
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
		Name:              "ChickenEggWarmer",
		ProductType:       Chicken,
		TechnologyType:    BuildingSpace,
		Tile:              g.Data["ChickenEggWarmerTile"].(Tile),
		TileWidth:         1,
		TileHeight:        1,
		TileFillSpace:     false,
		Square:            BoardSquare{},
		CostMoney:         25,
		Description:       "asdasd",
		CanBuild:          ChickenEggWarmerCanBuild,
		OnBuild:           ChickenEggWarmerOnBuild,
		Redraw:            false,
		OnRoundEnd:        ChickenEggWarmerRoundEnd,
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      false,
	}

}

func ChickenEggWarmerCanBuild(g *Game) bool {
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
