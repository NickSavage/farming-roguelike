package main

type Technology struct {
	Name  string
	Tiles []BoardSquare
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
		Technology: make([]Technology, 1),
		People:     make([]Person, 1),
	}
	g.Run.Technology = append(g.Run.Technology, g.CreateChickenCoopTech())

}

func (g *Game) EndRound() {
	g.Run.Money += 100
}

func (g *Game) CreateChickenCoopTech() Technology {
	result := Technology{
		Name:  "Chicken Coop",
		Tiles: make([]BoardSquare, 1),
	}
	result.Tiles = append(result.Tiles, BoardSquare{
		Tile:     g.Data["ChickenCoopTile"].(Tile),
		TileType: "Technology",
		Row:      10,
		Column:   10,
	})
	return result
}
