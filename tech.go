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

	tech["ChickenCoop"] = g.CreateChickenCoopTech()
	tech["WheatField"] = g.CreateWheatTech()
	tech["PotatoField"] = g.CreatePotatoTech()
	tech["CarrotField"] = g.CreateCarrotTech()

	tech["Workstation"] = g.CreateWorkstationTech()
	tech["ChickenEggWarmer"] = g.CreateChickenEggWarmer()

	g.Technology = tech
}

func (g *Game) CreateTechFromInitialData(input InitialData) Technology {
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

func (g *Game) RoundEndProduce(tech *Technology) float32 {

	return tech.BaseProduction * g.Run.Productivity * g.Run.Products[tech.ProductType].Yield * tech.TempYield
}

// chicken

func (g *Game) CreateChickenCoopTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Chicken Coop"])
	tech.CanBuild = ChickenCoopCanBuild
	tech.OnBuild = ChickenCoopOnBuild
	tech.OnClick = ChickenCoopOnClick
	tech.OnRoundEnd = ChickenCoopRoundEnd
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

func ChickenCoopOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		produced := g.RoundEndProduce(tech)
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

	tech := g.CreateTechFromInitialData(g.InitialData["Wheat"])
	tech.CanBuild = WheatFieldCanBuild
	tech.OnBuild = WheatFieldOnBuild
	tech.OnClick = WheatFieldOnClick
	tech.OnRoundEnd = WheatFieldRoundEnd
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
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	return nil
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
		produced := g.RoundEndProduce(tech)
		g.Run.Products["Wheat"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Wheat: %v", produced)
	}
	return ""
}

// potato

func (g *Game) CreatePotatoTech() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Potato"])
	tech.CanBuild = PotatoFieldCanBuild
	tech.OnBuild = PotatoFieldOnBuild
	tech.OnClick = PotatoFieldOnClick
	tech.OnRoundEnd = PotatoFieldRoundEnd
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
		produced := g.RoundEndProduce(tech)
		g.Run.Products["Potato"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Potatoes: %v", produced)
	}
	return ""
}

// carrot

func (g *Game) CreateCarrotTech() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Carrot"])
	tech.CanBuild = CarrotFieldCanBuild
	tech.OnBuild = CarrotFieldOnBuild
	tech.OnClick = CarrotFieldOnClick
	tech.OnRoundEnd = CarrotFieldRoundEnd
	tech.ShopButton = CarrotShopButton
	return &tech
}

func CarrotShopButton(g *Game) *ShopButton {
	result := &ShopButton{
		Width:      80,
		Height:     120,
		Image:      g.Data["CarrotShopIcon"].(Tile),
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

func CarrotFieldRoundEnd(g *Game, tech *Technology) {
	if g.Run.NextSeason == Winter {
		tech.ReadyToHarvest = true
	}

	tech.ReadyToTouch = false
}

func CarrotFieldOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		produced := g.RoundEndProduce(tech)
		g.Run.Products["Carrot"].Quantity += produced
		g.RemoveTech(tech)
		return fmt.Sprintf("Carrots: %v", produced)
	} else {
		return ""
	}
}

// workstation

func (g *Game) CreateWorkstationTech() *Technology {
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

	tech := g.CreateTechFromInitialData(g.InitialData["Carrot"])

	tech.CanBuild = ChickenEggWarmerCanBuild
	tech.OnBuild = ChickenEggWarmerOnBuild
	tech.OnClick = ChickenEggWarmerOnClick
	tech.OnRoundEnd = ChickenEggWarmerRoundEnd
	return &tech
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
