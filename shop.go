package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func ShopClickChickenCoop(g *Game) {
	g.Scenes["HUD"].Windows["ShopWindow"].Display = false
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.ScreenSkip = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = g.CreateChickenCoopTech()
}

func ShopClickWheatField(g *Game) {
	g.Scenes["HUD"].Windows["ShopWindow"].Display = false
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.ScreenSkip = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = g.CreateWheatTech()

}

func (g *Game) InitShopWindow() {
	scene := g.Scenes["Board"]
	buttons := []ShopButton{
		ShopButton{
			Width:       400,
			Height:      50,
			Title:       "Chicken Coop",
			Description: "sdasda",
			Image:       g.Data["ChickenCoopShopTile"].(Tile),
			OnClick:     ShopClickChickenCoop,
		},
		ShopButton{
			Width:       400,
			Height:      50,
			Title:       "Wheat Field",
			Description: "sdasda",
			Image:       g.Data["WheatTile"].(Tile),
			OnClick:     ShopClickWheatField,
		},
	}
	scene.Data["ShopButtons"] = buttons
}

func DrawShopWindow(g *Game, window *Window) {
	scene := g.Scenes["Board"]

	rl.DrawRectangle(200, 50, 900, 500, rl.White)
	rl.DrawText("Shop", 205, 55, 30, rl.Black)

	buttons := scene.Data["ShopButtons"].([]ShopButton)
	g.DrawShopButton(buttons[0], 205, 90)
	g.DrawShopButton(buttons[1], 205, 145)
	// for _, button := range buttons {
	// 	g.DrawShopButton(button, 205, 90)

	// }
}
