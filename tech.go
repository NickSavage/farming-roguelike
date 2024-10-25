package main

import (
	"fmt"
	"log"
	"sort"
	// rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) InitTechnology() {
	log.Printf("init tech")
	tech := make(map[string]*Technology)

	tech["ChickenCoop"] = g.ChickenCoop()
	tech["WheatField"] = g.WheatField()
	tech["PotatoField"] = g.PotatoField()

	tech["Workstation"] = g.Workstation()

	g.Technology = tech
}

func (g *Game) InitProduct(tech *Technology, price float32) {

	if _, exists := g.Run.Products[tech.ProductName]; !exists {
		g.Run.Products[tech.ProductName] = &Product{
			Name:     tech.ProductName,
			Quantity: 0,
			Price:    price,
		}
	}
}

func (g *Game) GetProductNames() []string {
	results := []string{}
	for _, product := range g.Run.Products {
		results = append(results, product.Name)
	}
	sort.Strings(results)
	return results
}

func (g *Game) RoundEndValue(tech *Technology, handler *TechnologyRoundHandler) float32 {
	units := handler.RoundEndProduce(g, tech)
	price := g.Run.Products[tech.ProductName].Price
	return units * price

}

func (g *Game) RoundEndText(tech *Technology, handler *TechnologyRoundHandler) string {

	units := handler.RoundEndProduce(g, tech)
	price := g.Run.Products[tech.ProductName].Price
	text := "$%v (%v units at $%v each)"
	return fmt.Sprintf(text, units*price, units, price)
}

func (g *Game) PlaceTech(tech *Technology, space *TechnologySpace) error {
	space.IsFilled = true
	space.Technology = tech

	if g.CanBuild(tech) {
		err := tech.OnBuild(g, tech)
		if err == nil {
			g.Run.Technology = append(g.Run.Technology, tech)
		}
	}
	return nil

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
		Name:           "Chicken Coop",
		ProductName:    "Chicken",
		TechnologyType: BuildingSpace,
		Tile:           g.Data["ChickenCoopTile"].(Tile),
		TileWidth:      1,
		TileHeight:     1,
		TileFillSpace:  false,
		Square:         BoardSquare{},
		Description:    "asdasd",
		CostMoney:      50,
		OnBuild:        ChickenCoopOnBuild,
		Redraw:         false,
		RoundHandler: []TechnologyRoundHandler{
			{
				OnRoundEnd:      ChickenCoopRoundEnd,
				RoundEndProduce: ChickenCoopProduce,
			},
		},
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

func ChickenCoopOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech, 5)
	return nil
}

func ChickenCoopProduce(g *Game, tech *Technology) float32 {
	return 5 * g.Run.Productivity
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
		Name:     "Wheat",
		Quantity: 0,
		Price:    1,
	}

	return result
}

func (g *Game) WheatField() *Technology {
	return &Technology{
		Name:           "Wheat",
		ProductName:    "Wheat",
		TechnologyType: PlantSpace,
		Tile:           g.Data["WheatTile"].(Tile),
		TileWidth:      1,
		TileHeight:     1,
		TileFillSpace:  true,
		Square:         BoardSquare{},
		CostMoney:      50,
		Description:    "asdasd",
		OnBuild:        WheatFieldOnBuild,
		Redraw:         false,
		RoundHandler: []TechnologyRoundHandler{
			{
				Season:          Spring,
				OnRoundEnd:      WheatFieldRoundSpring,
				RoundEndProduce: WheatFieldProduce,
			},
			{
				Season:          Summer,
				OnRoundEnd:      WheatFieldRoundSummer,
				RoundEndProduce: WheatFieldProduce,
			},
			{
				Season:          Autumn,
				OnRoundEnd:      WheatFieldRoundAutumn,
				RoundEndProduce: WheatFieldProduce,
			},
			{
				Season:          Winter,
				OnRoundEnd:      WheatFieldRoundWinter,
				RoundEndProduce: WheatFieldProduce,
			},
		},
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      true,
	}

}

func WheatFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech, 1)
	return nil
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
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundSummer(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundAutumn(g *Game, tech *Technology) {
	g.Run.Products["Wheat"].Quantity += WheatFieldProduce(g, tech)
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundWinter(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 0
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

	g.InitProduct(result, 5)
	return result
}

func (g *Game) PotatoField() *Technology {
	return &Technology{
		Name:           "Potato",
		ProductName:    "Potato",
		TechnologyType: PlantSpace,
		Tile:           g.Data["PotatoTile"].(Tile),
		TileWidth:      1,
		TileHeight:     1,
		TileFillSpace:  true,
		Square:         BoardSquare{},
		CostMoney:      50,
		Description:    "asdasd",
		OnBuild:        PotatoFieldOnBuild,
		Redraw:         false,
		RoundHandler: []TechnologyRoundHandler{
			{
				Season:          Spring,
				OnRoundEnd:      PotatoFieldRoundSpring,
				RoundEndProduce: PotatoFieldProduce,
			},
			{
				Season:          Summer,
				OnRoundEnd:      PotatoFieldRoundSummer,
				RoundEndProduce: PotatoFieldProduce,
			},
			{
				Season:          Autumn,
				OnRoundEnd:      PotatoFieldRoundAutumn,
				RoundEndProduce: PotatoFieldProduce,
			},
			{
				Season:          Winter,
				OnRoundEnd:      PotatoFieldRoundWinter,
				RoundEndProduce: PotatoFieldProduce,
			},
		},
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      true,
	}

}

func PotatoFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech, 1)
	return nil
}

func PotatoFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return float32(125) * g.Run.Productivity
	} else {
		return 0
	}
}

func PotatoFieldRoundSpring(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func PotatoFieldRoundSummer(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func PotatoFieldRoundAutumn(g *Game, tech *Technology) {
	g.Run.Products["Potato"].Quantity += PotatoFieldProduce(g, tech)
	tech.RoundHandlerIndex += 1
	tech.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func PotatoFieldRoundWinter(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 0
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
		Name:           "Workstation",
		ProductName:    "",
		TechnologyType: BuildingSpace,
		Tile:           g.Data["WorkstationTile"].(Tile),
		TileWidth:      1,
		TileHeight:     1,
		TileFillSpace:  false,
		Square:         BoardSquare{},
		CostMoney:      25,
		Description:    "asdasd",
		OnBuild:        WorkstationOnBuild,
		Redraw:         false,
		RoundHandler: []TechnologyRoundHandler{
			{
				OnRoundEnd: WorkstationRoundEnd,
			},
		},
		RoundCounterMax:   0,
		RoundCounter:      0,
		RoundHandlerIndex: 0,
		ShowEndRound:      false,
	}

}

func WorkstationOnBuild(g *Game, tech *Technology) error {
	g.Run.Productivity += 0.05

	return nil

}
func WorkstationRoundEnd(g *Game, tech *Technology) {

}
