package main

type Technology struct {
	Name string
	Tile BoardSquare
}

type Person struct {
}

type Run struct {
	Technology []Technology
	People     []Person
	Money      int32
}

func (g *Game) InitRun() {

	g.Run = &Run{
		Money:      100,
		Technology: make([]Technology, 0),
		People:     make([]Person, 1),
	}
	g.Run.Technology = append(g.Run.Technology, g.CreateChickenCoopTech())

}

func (g *Game) EndRound() {
	g.Run.Money += 100
}

func (g *Game) PlaceTech(tech *Technology, x, y float32) {

	row := int((x + TILE_WIDTH/2) / TILE_WIDTH)
	col := int((y + TILE_HEIGHT/2) / TILE_HEIGHT)

	// Store calculated row and column in tech
	tech.Tile.Row = row
	tech.Tile.Column = col

	g.Run.Technology = append(g.Run.Technology, *tech)

}

func (g *Game) CreateChickenCoopTech() Technology {
	result := Technology{
		Name: "Chicken Coop",
		Tile: BoardSquare{},
	}
	result.Tile = BoardSquare{
		Tile:     g.Data["ChickenCoopTile"].(Tile),
		TileType: "Technology",
		Row:      10,
		Column:   10,
		Width:    2,
		Height:   2,
	}
	return result
}
