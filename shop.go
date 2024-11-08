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
	if !tech.CanBuild(g, tech) {
		return
	}
	g.Scenes["Board"].Windows["ShopWindow"].Display = false
	space, err := g.GetOpenSpace(tech)

	if err == nil {
		g.PlaceTech(tech, space)
	}
}

func ShopClickChickenCoop(g *Game) {

	tech := g.CreateChickenCoopTech()
	g.ShopChooseTech(tech)
}

func ShopClickWheatField(g *Game) {

	tech := g.CreateWheatTech()
	g.ShopChooseTech(tech)
}
func ShopClickPotatoField(g *Game) {

	tech := g.CreatePotatoTech()
	g.ShopChooseTech(tech)
}

func ShopClickCarrotField(g *Game) {

	tech := g.CreateCarrotTech()
	g.ShopChooseTech(tech)
}
func ShopClickWorkstation(g *Game) {

	tech := g.CreateWorkstationTech()
	g.ShopChooseTech(tech)
}

func ShopClickChickenEggWarmer(g *Game) {
	tech := g.CreateChickenEggWarmer()
	g.ShopChooseTech(tech)
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
