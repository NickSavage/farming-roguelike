package main

func (g *Game) InitTechnology() {
	tech := make(map[string]Technology)

	tech["ChickenCoop"] = g.ChickenCoop()
	tech["WheatField"] = g.WheatField()

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
	return "Chicken Coop: $5"
}

func ChickenCoopRoundEndValue(g *Game, tech *Technology) float32 {
	return 5
}

func ChickenCoopRoundEnd(g *Game, tech *Technology) {
	g.Run.EndRoundMoney += 5
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
				CostActions:   1,
			},
			{
				Season:        Summer,
				OnRoundEnd:    WheatFieldRoundSummer,
				RoundEndValue: WheatFieldRoundEndValue,
				CostActions:   1,
			},
			{
				Season:        Autumn,
				OnRoundEnd:    WheatFieldRoundAutumn,
				RoundEndValue: WheatFieldRoundEndValue,
				CostActions:   1,
			},
			{
				Season:        Winter,
				OnRoundEnd:    WheatFieldRoundWinter,
				RoundEndValue: WheatFieldRoundEndValue,
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
	if g.Run.CurrentSeason == Autumn {
		return "Wheat Field: $125"
	} else {
		return "Wheat Field: $0"
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
	g.Run.EndRoundMoney += 125
	tech.RoundHandlerIndex += 1
	tech.Tile.Tile.TileFrame.X += 45
	tech.Redraw = true
}
func WheatFieldRoundWinter(g *Game, tech *Technology) {
	tech.RoundHandlerIndex += 0
	tech.Tile.Tile.TileFrame.X += 45
	tech.Redraw = true
}
