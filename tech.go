package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	//
	// rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) InitTechnology() {
	log.Printf("init tech")
	g.LoadInitialData()
	tech := make(map[string]*Technology)

	tech["ChickenCoop"] = g.ChickenCoop()
	tech["WheatField"] = g.WheatField()
	tech["PotatoField"] = g.PotatoField()
	tech["CarrotField"] = g.CarrotField()

	tech["Workstation"] = g.Workstation()
	tech["ChickenEggWarmer"] = g.ChickenEggWarmer()

	g.Technology = tech
}

func (g *Game) CreateTechFromInitialData(input InitialData) Technology {
	log.Printf("inptu %v", input)
	return Technology{
		Name:           input.Name,
		Description:    input.Description,
		TechnologyType: TechnologyType(input.TechnologyType),
		ProductType:    ProductType(input.ProductType),
		Rarity:         input.Rarity,
		Tile:           g.Data[input.TileConfig.ID].(Tile),
		TileWidth:      input.TileConfig.Width,
		TileHeight:     input.TileConfig.Height,
		TileFillSpace:  input.TileConfig.FillSpace,
		CostMoney:      input.CostMoney,
		Square:         BoardSquare{},
		TempYield:      1,
		ReadyToTouch:   true,
		InitialPrice:   input.Price,
		BaseProduction: input.Production,
	}
}

func (g *Game) LoadInitialData() {
	var initialData []InitialData

	// Load the JSON data from the file into memory
	file, err := os.Open("./assets/technology.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)

	err = jsonDecoder.Decode(&initialData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create an empty map to store the data
	dataMap := make(map[string]InitialData)

	// Iterate over each item in the initialData slice
	for _, item := range initialData {
		// Store the item in the map with its product type as key
		log.Printf("item %v", item)
		dataMap[item.Name] = item
	}
	log.Printf("data %v", dataMap)
	g.InitialData = dataMap
}

func (g *Game) InitProduct(productType ProductType, price float32) {

	log.Printf("init product %v  %v", productType, price)
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
	sort.Slice(results, func(i, j int) bool {
		strI := string(results[i])
		strJ := string(results[j])
		return strI < strJ
	})

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

	tech := g.CreateTechFromInitialData(g.InitialData["Chicken Coop"])
	tech.CanBuild = ChickenCoopCanBuild
	tech.OnBuild = ChickenCoopOnBuild
	tech.OnClick = ChickenCoopOnClick
	tech.OnRoundEnd = ChickenCoopRoundEnd
	tech.RoundEndProduce = ChickenCoopProduce
	return &tech
}

func ChickenCoopCanBuild(g *Game, tech *Technology) bool {
	return true
}

func ChickenCoopOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	tech.ReadyToTouch = false
	return nil
}

func ChickenCoopProduce(g *Game, tech *Technology) float32 {
	log.Printf("%v, %v, %v, %v", tech.BaseProduction, g.Run.Productivity, g.Run.Products["Chicken"].Yield, tech.TempYield)
	return tech.BaseProduction * g.Run.Productivity * g.Run.Products["Chicken"].Yield * tech.TempYield
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
		Price:    g.InitialData["Wheat"].Price,
		Yield:    1,
	}

	return result
}

func (g *Game) WheatField() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Wheat"])
	tech.CanBuild = WheatFieldCanBuild
	tech.OnBuild = WheatFieldOnBuild
	tech.OnClick = WheatFieldOnClick
	tech.OnRoundEnd = WheatFieldRoundEnd
	tech.RoundEndProduce = WheatFieldProduce
	tech.ShopButton = WheatShopButton
	return &tech
}

func WheatShopButton(g *Game) *ShopButton {
	result := &ShopButton{
		Width:      50,
		Height:     50,
		Image:      g.Data["WheatIcon"].(Tile),
		OnClick:    ShopClickWheatField,
		Technology: g.CreateWheatTech(),
	}
	return result
}

func WheatFieldCanBuild(g *Game, tech *Technology) bool {
	if g.Run.CurrentSeason == Spring {
		return true
	}
	return false
}

func WheatFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, g.InitialData["Wheat"].Price)
	return nil
}

func WheatFieldProduce(g *Game, tech *Technology) float32 {
	result := g.InitialData["Wheat"].Production * g.Run.Productivity * g.Run.Products["Wheat"].Yield * tech.TempYield
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

	g.InitProduct(result.ProductType, g.InitialData["Potato"].Price)
	return result
}

func (g *Game) PotatoField() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Potato"])
	tech.CanBuild = PotatoFieldCanBuild
	tech.OnBuild = PotatoFieldOnBuild
	tech.OnClick = PotatoFieldOnClick
	tech.OnRoundEnd = PotatoFieldRoundEnd
	tech.RoundEndProduce = PotatoFieldProduce
	tech.ShopButton = PotatoShopButton
	return &tech
}

func PotatoShopButton(g *Game) *ShopButton {
	result := &ShopButton{
		Width:      100,
		Height:     100,
		Image:      g.Data["PotatoIcon"].(Tile),
		OnClick:    ShopClickPotatoField,
		Technology: g.CreatePotatoTech(),
	}
	return result
}

func PotatoFieldCanBuild(g *Game, tech *Technology) bool {
	if g.Run.CurrentSeason == Spring {
		return true
	}
	return false
}

func PotatoFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, g.InitialData["Potato"].Price)
	return nil
}

func PotatoFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn {
		return g.InitialData["Potato"].Production * g.Run.Productivity * g.Run.Products["Potato"].Yield * tech.TempYield
	} else {
		return 0
	}
}
func PotatoFieldRoundEnd(g *Game, tech *Technology) {
	if g.Run.NextSeason == Autumn {
		tech.ReadyToHarvest = true
	} else if g.Run.NextSeason == Winter {
		g.RemoveTech(tech)
	} else {
		tech.ReadyToTouch = true
	}

	tech.Tile.TileFrame.X += 45
}

func PotatoFieldOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToTouch {
		err := g.Run.SpendAction(1)
		if err == nil {
			tech.TempYield += 0.05
			tech.ReadyToTouch = false
			return fmt.Sprintf("Yield: %v", tech.TempYield)
		}
		tech.ReadyToTouch = false
	} else if tech.ReadyToHarvest {
		produced := PotatoFieldProduce(g, tech)
		g.Run.Products["Potato"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Potatoes: %v", produced)
	}
	return ""
}

// carrot

func (g *Game) CreateCarrotTech() *Technology {
	result := g.CarrotField()
	result.Square = BoardSquare{
		TileType: "Technology",
		Width:    5,
		Height:   5,
		Occupied: true,
	}

	return result
}

func (g *Game) CarrotField() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Carrot"])
	tech.CanBuild = CarrotFieldCanBuild
	tech.OnBuild = CarrotFieldOnBuild
	tech.OnClick = CarrotFieldOnClick
	tech.OnRoundEnd = CarrotFieldRoundEnd
	tech.RoundEndProduce = CarrotFieldProduce
	tech.ShopButton = CarrotShopButton
	return &tech
}

func CarrotShopButton(g *Game) *ShopButton {
	result := &ShopButton{
		Width:      80,
		Height:     120,
		Image:      g.Data["WorkstationTile"].(Tile),
		OnClick:    ShopClickCarrotField,
		Technology: g.CreateCarrotTech(),
	}
	return result
}

func CarrotFieldCanBuild(g *Game, tech *Technology) bool {
	if g.Run.CurrentSeason == Spring || g.Run.CurrentSeason == Autumn {
		return true
	}
	return false
}

func CarrotFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	return nil
}

func CarrotFieldProduce(g *Game, tech *Technology) float32 {
	if g.Run.CurrentSeason == Autumn && g.Run.NextSeason == Winter {
		return tech.BaseProduction * g.Run.Productivity * g.Run.Products["Carrot"].Yield * tech.TempYield
	} else {
		return 0
	}
}

func CarrotFieldRoundEnd(g *Game, tech *Technology) {
	if g.Run.NextSeason == Winter && g.Run.CurrentSeason != Autumn {
		tech.ReadyToHarvest = true
		tech.Tile.TileFrame.X += 50
	}

	tech.ReadyToTouch = false
}

func CarrotFieldOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		produced := CarrotFieldProduce(g, tech)
		g.Run.Products["Carrot"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Carrots: %v", produced)
	} else {
		return ""
	}
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
		TechnologyType: "BuildingSpace",
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
		TechnologyType: "BuildingSpace",
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
	g.Run.Products["Chicken"].Yield += 0.05

	return nil

}
func ChickenEggWarmerRoundEnd(g *Game, tech *Technology) {

}

func ChickenEggWarmerOnClick(g *Game, tech *Technology) string {
	return ""
}
