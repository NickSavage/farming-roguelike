package main

import (
	"errors"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
)

func (g *Game) ShopChooseTech(tech *Technology) error {

	if !g.CanBuild(tech) {
		return errors.New("cannot build")
	}
	g.Scenes["Board"].Windows["ShopWindow"].Display = false
	space, err := g.GetOpenSpace(tech)

	if err == nil {
		g.PlaceTech(tech, space)
	}
	return nil
}

func ShopButtonOnClick(g *Game, b ShopBuildingButton) {
	window := g.Scenes["Board"].Windows["ShopWindow"]
	err := g.ShopChooseTech(b.Technology)
	log.Printf("err %v", err)
	if err == nil {
		// setting purchased flag
		// this is a bit of a cludge, think about another way at some point

		// don't need to do it for seeds
		if b.Technology.TechnologyType == Seed {
			return
		}
		var button ShopBuildingButton
		for i, _ := range window.Components {
			// log.Printf("component %v i %v", component, i)
			// log.Printf("button %v b %v", button.Position, b.Position)
			if button.Position == b.Position {
				button.Purchased = true
				window.Components[i] = &button
			}
		}
	}
}

func ShopSeedButtonOnClick(g *Game, b *ShopSeedButton) {
	// window := g.Scenes["Board"].Windows["ShopWindow"]
	g.ShopChooseTech(b.Technology)
	b.Technology.ToBeDeleted = true
}

func (g *Game) InitShopWindow() {
	g.InitShopRoundComponents()
}

// run each time the shop is opened, maybe should be each time the round is changed
func (g *Game) InitShopRoundComponents() {
	window := g.Scenes["Board"].Windows["ShopWindow"]
	window.Components = make([]engine.UIComponent, 0)
	buildings := g.Run.CurrentRoundShopBuildings

	n := 1
	blank := engine.NewBlankComponent()
	blank.SelectDirections.Right = 1
	blank.SelectDirections.Left = 5
	window.Components = append(window.Components, &blank)

	x := window.Rectangle.X
	y := window.Rectangle.Y

	for i, building := range buildings {
		// for i, _ := range buildings {
		rect := rl.NewRectangle(x+50+float32(i*160), y+45, 150, 300)
		button := g.NewShopButton(rect, building)
		button.SelectDirections.Right = n + 1
		button.SelectDirections.Left = n - 1

		button.Position = n
		button.ExpandedButton = true
		button.CanBuild = g.CanBuild(building)

		n += 1
		if n == 6 {
			button.SelectDirections.Right = 1
		}
		window.Components = append(window.Components, &button)
	}

	seeds := make([]*Technology, 0)
	for i, seed := range g.Run.CurrentSeeds {
		log.Printf("seed %v", seed.ToBeDeleted)
		if seed.ToBeDeleted {
			continue
		}
		seeds = append(seeds, seed)

		rect := rl.NewRectangle(x+50+float32(i*160), y+55+300, 150, 150)
		button := g.NewShopSeedButton(rect, seed)
		button.CanBuild = g.CanBuild(seed)
		button.Position = n
		window.Components = append(window.Components, &button)
		n += 1
	}
	g.Run.CurrentSeeds = seeds

}

func DrawShopWindow(gi engine.GameInterface, win *engine.Window) {
	g := gi.(*Game)
	window := g.Scenes["Board"].Windows["ShopWindow"]

	x := int32(window.Rectangle.X)
	y := int32(window.Rectangle.Y)

	rl.DrawRectangleRec(window.Rectangle, rl.White)
	rl.DrawRectangleLinesEx(window.Rectangle, 5, rl.Black)
	rl.DrawText("Shop", x+5, y+5, 30, rl.Black)

}

func (g *Game) ShopRandomBuildings(needed int) []*Technology {

	keysToPickFrom := make([]string, 0)
	for key, tech := range g.Technology {
		log.Printf("tech %v unlocked %v", tech, tech.Unlocked)
		if !tech.Unlocked {
			continue
		}
		keysToPickFrom = append(keysToPickFrom, key)
		// some sort of filtering is needed here

	}
	results := g.PickRandomTechnologies(needed, keysToPickFrom)

	return results
}
