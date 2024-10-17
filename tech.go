package main

func (g *Game) InitTechnology() {
	tech := make(map[string]Technology)
	chicken := Technology{
		Name:          "Chicken Coop",
		Tile:          BoardSquare{},
		Description:   "asdasd",
		Cost:          50,
		OnBuild:       ChickenCoopOnBuild,
		OnRoundEnd:    ChickenCoopRoundEnd,
		RoundEndValue: ChickenCoopRoundEndValue,
	}
	tech["ChickenCoop"] = chicken

	tech["WheatField"] = Technology{
		Name:            "Wheat",
		Tile:            BoardSquare{},
		Cost:            50,
		Description:     "asdasd",
		OnBuild:         WheatFieldOnBuild,
		OnRoundEnd:      WheatFieldRoundEnd,
		RoundCounterMax: 4,
		RoundCounter:    4,
		RoundEndValue:   WheatFieldRoundEndValue,
	}

	g.Data["Technology"] = tech
}

func (g *Game) CreateChickenCoopTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["ChickenCoop"]
	result.Tile = BoardSquare{
		Tile:        g.Data["ChickenCoopTile"].(Tile),
		TileType:    "Technology",
		Row:         10,
		Column:      10,
		Width:       2,
		Height:      2,
		Occupied:    true,
		MultiSquare: true,
	}
	return &result
}

func (g *Game) CreateWheatTech() *Technology {

	tech := g.Data["Technology"].(map[string]Technology)
	result := tech["WheatField"]
	result.Tile = BoardSquare{
		Tile:     g.Data["WheatTile"].(Tile),
		TileType: "Technology",
		Row:      8,
		Column:   8,
		Width:    5,
		Height:   5,
		Occupied: true,
	}
	return &result
}

func ChickenCoopOnBuild(g *Game, tech *Technology) {
	g.Run.Productivity += 0.05
	g.Run.RoundActionsRemaining -= 1
	g.Run.Money -= tech.Cost
}

func ChickenCoopRoundEndText(g *Game, tech *Technology) string {
	return "Chicken Coop: $5"
}

func ChickenCoopRoundEndValue(g *Game, tech *Technology) float32 {
	return 5
}

func ChickenCoopRoundEnd(g *Game, tech *Technology) {
	g.Run.EndRoundMoney += 5
}

func WheatFieldOnBuild(g *Game, tech *Technology) {

	g.Run.RoundActionsRemaining -= 1
	g.Run.Money -= tech.Cost
}

func WheatFieldRoundEndValue(g *Game, tech *Technology) float32 {
	if tech.RoundCounter-1 == 0 {
		return 125
	} else {
		return 0
	}
}
func WheatFieldRoundEndText(g *Game, tech *Technology) string {
	if tech.RoundCounter-1 == 0 {
		return "Wheat Field: $125"
	} else {
		return "Wheat Field: $0"
	}

}

func WheatFieldRoundEnd(g *Game, tech *Technology) {

	g.Run.RoundActionsRemaining -= 1
	tech.Tile.Tile.TileFrame.X += 45
	tech.RoundCounter -= 1
	g.Run.EndRoundMoney += WheatFieldRoundEndValue(g, tech)
	if tech.RoundCounter == 0 {
		tech.RoundCounter = tech.RoundCounterMax
	}

}
