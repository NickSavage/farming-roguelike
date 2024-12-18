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

	tech["Field"] = g.CreateFieldTech()

	tech["Chicken Coop"] = g.CreateChickenCoopTech()
	tech["Wheat Field"] = g.CreateWheatTech()
	tech["Potato Field"] = g.CreatePotatoTech()
	tech["Carrot Field"] = g.CreateCarrotTech()

	tech["Flour Mill"] = g.CreateFlourMillTech()

	tech["Cow Pasture"] = g.CreateCowPastureTech()
	tech["Cow Slaughterhouse"] = g.CreateCowSlaughterhouseTech()

	tech["Workstation"] = g.CreateWorkstationTech()
	tech["Fertilizer"] = g.CreateFertilizerTech()
	tech["Chicken Egg Warmer"] = g.CreateChickenEggWarmer()

	tech["Cell Tower"] = g.CreateCellTowerTech()
	tech["Solar Panels"] = g.CreateSolarPanelTech()

	g.Technology = tech
}

// to check timing of whether the seed should be planted or not
func (g *Game) CheckSeedPlanting(tech *Technology) bool {

	if tech.ProductType == Wheat {
		if g.Run.CurrentSeason == Spring {
			return true
		}
		return false
	}
	if tech.ProductType == Potato {
		if g.Run.CurrentSeason == Spring {
			return true
		}
		return false
	}
	if tech.ProductType == Carrot {
		if g.Run.CurrentSeason == Spring {
			return true
		}
		return false
	}
	return false

}

func (g *Game) CanBuild(tech *Technology) bool {

	if !g.Run.CanSpendMoney(tech.CostMoney) {
		log.Printf("cant spend money")
		return false
	}
	if !g.Run.CanSpendAction(tech.CostActions) {
		log.Printf("can't spend actions")
		return false
	}
	_, err := g.GetOpenSpace(tech)
	if err != nil {
		log.Printf("no space")
		return false
	}
	if tech.TechnologyType == Seed {
		// check if available field
		found := false
		for _, space := range g.Run.TechnologySpaces {
			if !space.IsField {
				continue
			}
			if len(space.PlantedSeeds) > 0 {
				continue
			}
			found = true
		}
		if found {
			check := g.CheckSeedPlanting(tech)
			if !check {
				return false
			}
		}
		log.Printf("%v can build %v", tech.Name, found)
		return found

	}
	if tech.TechnologyType == PlantSpace {
		if tech.Name == "Wheat Field" {
			if g.Run.CurrentSeason == Spring {
				return true
			}
		} else if tech.Name == "Potato Field" {
			if g.Run.CurrentSeason == Spring {
				return true
			}
		} else if tech.Name == "Carrot Field" {
			if g.Run.CurrentSeason == Spring || g.Run.CurrentSeason == Autumn {
				return true
			}
		} else if tech.Name == "Cow Pasture" {
			return true
		}
		return false

	}
	return true
}

func (g *Game) CreateTechFromInitialData(input InitialData) *Technology {
	log.Printf("input %v unlock %v", input.Name, input.DefaultUnlocked)
	return &Technology{
		Name:           input.Name,
		Description:    input.Description,
		TechnologyType: TechnologyType(input.TechnologyType),
		ProductType:    ProductType(input.ProductType),
		Rarity:         input.Rarity,
		Tile:           g.Data[input.TileConfig.ID].(Tile),
		TileWidth:      input.TileConfig.Width,
		TileHeight:     input.TileConfig.Height,
		TileFillSpace:  input.TileConfig.FillSpace,
		ShopIcon:       input.ShopIcon,
		CostMoney:      input.CostMoney,
		CostActions:    input.CostActions,
		Square:         BoardSquare{},
		TempYield:      1,
		ReadyToTouch:   true,
		InitialPrice:   input.Price,
		BaseProduction: input.Production,
		Input:          input.Input,
		Unlocked:       input.DefaultUnlocked,
		ShowInShop:     input.ShowInShop,
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
	unlocks := []UnlockJSON{}

	for _, item := range initialData {
		item.DefaultUnlocked = true
		if item.Unlock != nil {
			item.Unlock.TechnologyName = item.Name
			unlocks = append(unlocks, *item.Unlock)
			item.DefaultUnlocked = false
		}
		dataMap[item.Name] = item
	}
	g.InitialData = dataMap
	g.UnlockBaseData = unlocks
}

func (g *Game) InitProduct(productType ProductType, price float32) {

	log.Printf("init product %v  %v", productType, price)
	log.Printf("stats %v", g.ProductStats[productType])
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

func (r *Run) AddSeed(tech *Technology) {
	r.CurrentSeeds = append(r.CurrentSeeds, tech)
}
func (g *Game) PlaceSeed(tech *Technology, space *TechnologySpace) error {
	tech.Space = space
	copy := *tech
	space.PlantedSeeds = append(space.PlantedSeeds, &copy)

	err := g.Run.SpendAction(tech.CostActions)
	if err != nil {
		return errors.New("cannot spend action")
	}

	err = g.Run.SpendMoney(tech.CostMoney)

	if err == nil {
		err := tech.OnBuild(g, tech)
		if err == nil {
			g.Run.Technology = append(g.Run.Technology, &copy)
		}
	}
	results := []*Technology{}
	for _, tech := range g.Run.CurrentSeeds {
		if !tech.ToBeDeleted {
			results = append(results, tech)
		}
	}
	g.Run.CurrentSeeds = results
	return nil
}

func (g *Game) PlaceTech(tech *Technology, space *TechnologySpace) error {
	if tech.TechnologyType == Seed {
		return g.PlaceSeed(tech, space)
	}
	space.IsFilled = true
	tech.Space = space
	copy := *tech
	space.Technology = &copy

	err := g.Run.SpendAction(tech.CostActions)
	if err != nil {
		return errors.New("cannot spend action")
	}

	err = g.Run.SpendMoney(tech.CostMoney)

	if err == nil {
		if tech.TechnologyType == Field {
			log.Printf("yes this happens")
			space.IsField = true
		}
		err := tech.OnBuild(g, tech)
		if err == nil {
			g.Run.Technology = append(g.Run.Technology, space.Technology)
		}
	}
	return nil
}

func (g *Game) GetOpenSpace(tech *Technology) (*TechnologySpace, error) {

	for _, space := range g.Run.TechnologySpaces {
		if tech.TechnologyType == Seed && space.IsField {
			if len(space.PlantedSeeds) != 0 {
				continue
			}
			return space, nil
		}
		if space.IsFilled {
			continue
		}
		if tech.Name == "Solar Panels" {
			if space.TechnologyType != PlantSpace {
				continue
			}

		} else if space.TechnologyType != tech.TechnologyType {
			continue
		}
		return space, nil
	}
	return &TechnologySpace{}, errors.New("no empty space")
}

func (g *Game) RemoveTech(tech *Technology) {

	log.Printf("space %v", tech.Space)
	tech.ToBeDeleted = true
	if tech.TechnologyType == Seed {
		new := make([]*Technology, 0)
		for _, seed := range tech.Space.PlantedSeeds {
			if seed.ToBeDeleted {
				continue
			}
			new = append(new, seed)
		}
		tech.Space.PlantedSeeds = new

	} else {
		space := tech.Space
		space.IsFilled = false
		space.Technology = &Technology{}
	}

	var results []*Technology
	for _, tech := range g.Run.Technology {
		if !tech.ToBeDeleted {
			results = append(results, tech)
		}
	}
	g.Run.Technology = results
}

func (g *Game) HandleClickTech(space *TechnologySpace) string {
	if space.TechnologyType == Field {
		log.Printf("test")
		log.Printf("click %v", space.PlantedSeeds)
		result := ""
		for _, seed := range space.PlantedSeeds {
			result += seed.OnClick(g, seed)
		}
		return result
	} else {
		return space.Technology.OnClick(g, space.Technology)
	}
}

func (g *Game) RoundEndProduce(tech *Technology) float32 {
	if tech.Name == "Solar Panels" {
		return tech.BaseProduction * g.Run.Products[tech.ProductType].Yield * tech.TempYield
	}
	return tech.BaseProduction * g.Run.Productivity * g.Run.Products[tech.ProductType].Yield * tech.TempYield
}

func (g *Game) ShopButton(tech *Technology) *ShopButton {
	result := &ShopButton{
		Width:      50,
		Height:     50,
		Image:      g.Data[tech.ShopIcon].(Tile),
		Technology: tech,
	}
	return result
}

// fields

func (g *Game) CreateFieldTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Field"])
	tech.OnBuild = FieldOnBuild
	tech.OnClick = FieldOnClick
	tech.OnRoundEnd = FieldRoundEnd
	return tech

}

func FieldOnBuild(g *Game, tech *Technology) error {
	tech.ReadyToTouch = false
	return nil

}
func FieldRoundEnd(g *Game, tech *Technology) {

}
func FieldOnClick(g *Game, tech *Technology) string {
	return ""
}

// chicken

func (g *Game) CreateChickenCoopTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Chicken Coop"])
	tech.OnBuild = ChickenCoopOnBuild
	tech.OnClick = ChickenCoopOnClick
	tech.OnRoundEnd = ChickenCoopRoundEnd
	return tech
}

func ChickenCoopOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	tech.ReadyToTouch = false
	return nil
}

func ChickenCoopOnClick(g *Game, tech *Technology) string {
	log.Printf("?")
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

	tech := g.CreateTechFromInitialData(g.InitialData["Wheat Field"])
	tech.OnBuild = WheatFieldOnBuild
	tech.OnClick = WheatFieldOnClick
	tech.OnRoundEnd = WheatFieldRoundEnd
	return tech
}

func WheatFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	tech.ReadyToHarvest = false
	tech.ReadyToTouch = true
	return nil
}

func WheatFieldRoundEnd(g *Game, tech *Technology) {
	log.Printf("does this run?")
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
	log.Printf("tech %v", tech)
	if tech.ReadyToTouch {
		log.Printf("touch")
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

		copy := *tech
		g.Run.AddSeed(&copy)

		g.RemoveTech(tech)

		return fmt.Sprintf("Wheat: %v", produced)
	}
	return ""
}

func (g *Game) CreateFlourMillTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Flour Mill"])
	tech.OnBuild = FlourMillOnBuild
	tech.OnClick = FlourMillOnClick
	tech.OnRoundEnd = FlourMillRoundEnd
	return tech

}

func FlourMillOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	g.InitProduct(Flour, 5)
	tech.ReadyToTouch = false
	return nil

}
func FlourMillRoundEnd(g *Game, tech *Technology) {
	tech.ReadyToHarvest = true
}
func FlourMillOnClick(g *Game, tech *Technology) string {

	if tech.ReadyToHarvest {
		err := g.Run.SpendAction(1)
		if err != nil {
			return ""
		}

		var input float32
		input = tech.Input.MaximumInput
		market := g.ConsumeOrBuyProduct(g.Run.Products[Wheat], tech.Input.MaximumInput)

		produced := input * tech.Input.OutputPerInput * g.Run.Productivity
		g.Run.Products[Flour].Quantity += produced
		tech.ReadyToHarvest = false
		if market > 0 {
			return fmt.Sprintf("Flour: %v (-$%v)", produced, market)

		} else {
			return fmt.Sprintf("Flour: %v (-%v Wheat)", produced, input)
		}
	}
	return ""
}

func FlourMillUnlockOtherCost(g *Game) bool {
	var total float32

	if current, exists := g.ProductStats[Wheat]; exists {
		total = current.TotalProduction
	} else {
		total = 0
	}
	if total > 500 {
		return true
	}
	return false
}
func FlourMillUnlockOtherCostDescription(g *Game) string {

	var total float32

	if current, exists := g.ProductStats[Wheat]; exists {
		total = current.TotalProduction
	} else {
		total = 0
	}

	return fmt.Sprintf("Sell %v/%v Wheat", total, 500)
}

// potato

func (g *Game) CreatePotatoTech() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Potato Field"])
	tech.OnBuild = PotatoFieldOnBuild
	tech.OnClick = PotatoFieldOnClick
	tech.OnRoundEnd = PotatoFieldRoundEnd
	return tech
}

func PotatoFieldOnBuild(g *Game, tech *Technology) error {
	// prob shouldn't be using initial data here
	g.InitProduct(tech.ProductType, g.InitialData["Potato Field"].Price)
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
	tech := g.CreateTechFromInitialData(g.InitialData["Carrot Field"])
	tech.OnBuild = CarrotFieldOnBuild
	tech.OnClick = CarrotFieldOnClick
	tech.OnRoundEnd = CarrotFieldRoundEnd
	return tech
}

func CarrotFieldOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	tech.ReadyToTouch = false
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

	tech := g.CreateTechFromInitialData(g.InitialData["Workstation"])
	tech.OnBuild = WorkstationOnBuild
	tech.OnClick = WorkstationOnClick
	tech.OnRoundEnd = WorkstationRoundEnd
	return tech

}

func WorkstationOnBuild(g *Game, tech *Technology) error {
	g.Run.Productivity += 0.05
	tech.ReadyToTouch = false

	return nil

}
func WorkstationRoundEnd(g *Game, tech *Technology) {

}
func WorkstationOnClick(g *Game, tech *Technology) string {
	return ""
}

func (g *Game) CreateChickenEggWarmer() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Chicken Egg Warmer"])

	tech.OnBuild = ChickenEggWarmerOnBuild
	tech.OnClick = ChickenEggWarmerOnClick
	tech.OnRoundEnd = ChickenEggWarmerRoundEnd
	return tech
}

func ChickenEggWarmerOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	g.Run.Products["Chicken"].Yield += 0.05
	tech.ReadyToTouch = false
	return nil

}
func ChickenEggWarmerRoundEnd(g *Game, tech *Technology) {

}

func ChickenEggWarmerOnClick(g *Game, tech *Technology) string {
	return ""
}

func ChickenEggWarmerUnlockOtherCost(g *Game) bool {
	var total float32

	if current, exists := g.ProductStats[Chicken]; exists {
		total = current.TotalProduction
	} else {
		total = 0
	}
	if total > 100 {
		return true
	}
	return false
}
func ChickenEggWarmerUnlockOtherCostDescription(g *Game) string {

	var total float32

	if current, exists := g.ProductStats[Chicken]; exists {
		total = current.TotalProduction
	} else {
		total = 0
	}

	return fmt.Sprintf("Sell %v/%v Chicken", total, 100)
}

// workstation

func (g *Game) CreateCellTowerTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Cell Tower"])
	tech.OnBuild = CellTowerOnBuild
	tech.OnClick = CellTowerOnClick
	tech.OnRoundEnd = CellTowerRoundEnd
	return tech

}

func CellTowerOnBuild(g *Game, tech *Technology) error {
	tech.ReadyToTouch = false

	return nil

}
func CellTowerRoundEnd(g *Game, tech *Technology) {
	g.Run.EndRoundMoney += tech.BaseProduction

}
func CellTowerOnClick(g *Game, tech *Technology) string {
	return ""
}

// solar panels

func (g *Game) CreateSolarPanelTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Solar Panels"])
	tech.OnBuild = SolarPanelOnBuild
	tech.OnClick = SolarPanelOnClick
	tech.OnRoundEnd = SolarPanelRoundEnd
	return tech

}

func SolarPanelOnBuild(g *Game, tech *Technology) error {
	tech.ReadyToTouch = false
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	return nil

}
func SolarPanelRoundEnd(g *Game, tech *Technology) {
	tech.ReadyToHarvest = true

}
func SolarPanelOnClick(g *Game, tech *Technology) string {
	log.Printf("?")

	if tech.ReadyToHarvest {
		produced := g.RoundEndProduce(tech)
		g.Run.Products["Solar"].Quantity += produced

		moneyEarned := g.SellProduct(g.Run.Products["Solar"])
		g.Run.Money += moneyEarned

		return fmt.Sprintf("Solar: %v", moneyEarned)
	}
	return ""
}

// fertilizer

func (g *Game) CreateFertilizerTech() *Technology {
	tech := g.CreateTechFromInitialData(g.InitialData["Fertilizer"])

	tech.OnBuild = FertilizerOnBuild
	tech.OnClick = FertilizerOnClick
	tech.OnRoundEnd = FertilizerRoundEnd

	return tech
}

func FertilizerOnBuild(g *Game, tech *Technology) error {
	g.Run.Productivity += 0.05
	tech.ReadyToTouch = false

	return nil
}
func FertilizerRoundEnd(g *Game, tech *Technology) {

}
func FertilizerOnClick(g *Game, tech *Technology) string {
	return ""
}

// cow slaughterhouse

func (g *Game) CreateCowSlaughterhouseTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Cow Slaughterhouse"])

	tech.OnBuild = CowSlaughterhouseOnBuild
	tech.OnClick = CowSlaughterhouseOnClick
	tech.OnRoundEnd = CowSlaughterhouseRoundEnd

	return tech
}

func CowSlaughterhouseOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	g.InitProduct(Cow, 2)
	return nil

}

func CowSlaughterhouseOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		err := g.Run.SpendAction(1)
		if err != nil {
			return ""
		}

		var input float32
		input = tech.Input.MaximumInput
		market := g.ConsumeOrBuyProduct(g.Run.Products["Cow"], tech.Input.MaximumInput)

		produced := input * tech.Input.OutputPerInput * g.Run.Productivity
		g.Run.Products[Beef].Quantity += produced
		tech.ReadyToHarvest = false
		if market > 0 {
			return fmt.Sprintf("Beef: %v (-$%v)", produced, market)

		} else {
			return fmt.Sprintf("Beef: %v (-%v Cow)", produced, input)
		}
	}
	return ""

}

func CowSlaughterhouseRoundEnd(g *Game, tech *Technology) {

	tech.ReadyToHarvest = true
}

func CowSlaughterhouseUnlockOtherCost(g *Game) bool {
	var total float32

	if current, exists := g.ProductStats[Cow]; exists {
		total = current.TotalProduction
	} else {
		total = 0
	}
	if total > 1000 {
		return true
	}
	return false
}
func CowSlaughterhouseUnlockOtherCostDescription(g *Game) string {

	var total float32

	if current, exists := g.ProductStats[Cow]; exists {
		total = current.TotalProduction
	} else {
		total = 0
	}

	return fmt.Sprintf("Sell %v/%v Cows", total, 1000)
}

// cow pasture

func (g *Game) CreateCowPastureTech() *Technology {

	tech := g.CreateTechFromInitialData(g.InitialData["Cow Pasture"])

	tech.OnBuild = CowPastureOnBuild
	tech.OnClick = CowPastureOnClick
	tech.OnRoundEnd = CowPastureRoundEnd

	return tech
}

func CowPastureOnBuild(g *Game, tech *Technology) error {
	g.InitProduct(tech.ProductType, tech.InitialPrice)
	g.InitProduct(Cow, 2)
	return nil

}

func CowPastureOnClick(g *Game, tech *Technology) string {
	if tech.ReadyToHarvest {
		err := g.Run.SpendAction(1)
		if err == nil {
			produced := g.RoundEndProduce(tech)
			g.Run.Products["Cow"].Quantity += produced
			tech.ReadyToHarvest = false
			return fmt.Sprintf("Cow: %v", produced)
		}
	}
	return ""

}

func CowPastureRoundEnd(g *Game, tech *Technology) {

	tech.ReadyToHarvest = true
}
