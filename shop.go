package main

import (
	"errors"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
)

func (g *Game) ShopChooseTech(tech *Technology) error {

	if !g.Run.CanSpendMoney(tech.CostMoney) {
		return errors.New("cannot spend money")
	}
	if !g.Run.CanSpendAction(tech.CostActions) {
		return errors.New("cannot spend action")
	}
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
	log.Printf("does this happen???")
	if err == nil {
		// this is a bit of a cludge, think about another way at some point
		var button ShopBuildingButton
		var ok bool
		for i, component := range window.Components {
			log.Printf("component %v i %v", component, i)
			if button, ok = component.(ShopBuildingButton); !ok {
				log.Printf("button is not a shopbuilding button")
				continue
			}
			log.Printf("button %v b %v", button.Position, b.Position)
			if button.Position == b.Position {
				log.Printf("does this happen?")
				button.Purchased = true
				window.Components[i] = button
			}
		}
	}
}

func (g *Game) InitShopWindow() {
	log.Printf("init shop")
}

// run each time the shop is opened, maybe should be each time the round is changed
func (g *Game) InitShopRoundComponents() {
	window := g.Scenes["Board"].Windows["ShopWindow"]
	window.Components = make([]engine.UIComponent, 0)

	buildings := g.Run.CurrentRoundShopBuildings

	n := 0

	var x float32 = float32(window.X)
	var y float32 = float32(window.Y)

	for i, building := range buildings {
		// for i, _ := range buildings {
		rect := rl.NewRectangle(x+50+float32(i*160), y+45, 150, 300)
		button := g.NewShopButton(rect, building)
		n += 1
		button.Position = n
		button.ExpandedButton = true
		window.Components = append(window.Components, button)
	}

	plants := g.Run.CurrentRoundShopPlants

	for i, plant := range plants {
		rect := rl.NewRectangle(x+50+float32(i*105), y+355, 100, 100)
		button := g.NewShopPlantButton(rect, plant)
		n += 1
		button.Position = n
		window.Components = append(window.Components, button)
	}
}

func DrawShopWindow(gi engine.GameInterface, win *engine.Window) {
	g := gi.(*Game)
	window := g.Scenes["Board"].Windows["ShopWindow"]

	var x int32 = window.X
	var y int32 = window.Y
	rl.DrawRectangle(x, y, window.Width, window.Height, rl.White)
	rl.DrawText("Shop", x+5, y+5, 30, rl.Black)

}
