package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
)

func (g *Game) ShopChooseTech(tech *Technology) {

	if !g.Run.CanSpendMoney(tech.CostMoney) {
		return
	}
	if !g.Run.CanSpendAction(tech.CostActions) {
		return
	}
	if !g.CanBuild(tech) {
		return
	}
	g.Scenes["Board"].Windows["ShopWindow"].Display = false
	space, err := g.GetOpenSpace(tech)

	if err == nil {
		g.PlaceTech(tech, space)
	}
}

func ShopButtonOnClick(g *Game, b ShopBuildingButton) {
	g.ShopChooseTech(b.Technology)
}

func (g *Game) InitShopWindow() {
	log.Printf("init shop")
}

// run each time the shop is opened, maybe should be each time the round is changed
func (g *Game) InitShopRoundComponents() {
	window := g.Scenes["Board"].Windows["ShopWindow"]
	window.Components = make([]engine.UIComponent, 0)

	buildings := g.Run.CurrentRoundShopBuildings

	for i, building := range buildings {
		// for i, _ := range buildings {
		rect := rl.NewRectangle(float32(215+i*160), 90, 150, 300)
		button := g.NewShopButton(rect, building)
		button.ExpandedButton = true
		window.Components = append(window.Components, button)
	}

	plants := g.Run.CurrentRoundShopPlants

	var x float32 = 215
	for i, plant := range plants {
		rect := rl.NewRectangle(x+float32(i*100), 400, 100, 100)
		button := g.NewShopPlantButton(rect, plant)
		window.Components = append(window.Components, button)
	}
}

func DrawShopWindow(gi engine.GameInterface, win *engine.Window) {
	// g := gi.(*Game)
	// scene := g.Scenes["Board"]

	rl.DrawRectangle(200, 50, 900, 500, rl.White)
	rl.DrawText("Shop", 205, 55, 30, rl.Black)

}
