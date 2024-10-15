package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func ShopClickChickenCoop(g *Game) {
	g.Scenes["HUD"].Data["DisplayShopWindow"] = false
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.Scenes["Board"].Data["PlaceTechSkip"] = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = &g.Run.Technology[0]

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
	}
	scene.Data["ShopButtons"] = buttons
}

type ShopButton struct {
	X           float32
	Y           float32
	Width       int32
	Height      int32
	Title       string
	Description string
	Image       Tile
	OnClick     func(*Game)
}

func (g *Game) DrawShopButton(shopButton ShopButton, x, y float32) {
	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  float32(shopButton.Width),
		Height: float32(shopButton.Height),
	}
	rl.DrawRectangleLinesEx(rect, 1, rl.Black)
	DrawTile(shopButton.Image, x+5, y+2)
	rl.DrawText(shopButton.Title, int32(x+50), int32(y+2), 20, rl.Black)
	rl.DrawText(shopButton.Description, int32(x+50), int32(y+22), 10, rl.Black)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, rect) {
			shopButton.OnClick(g)
		}
	}
}

func (g *Game) DrawShopWindow() {
	scene := g.Scenes["Board"]

	rl.DrawRectangle(200, 50, 900, 500, rl.White)
	rl.DrawText("Shop", 205, 55, 30, rl.Black)

	buttons := scene.Data["ShopButtons"].([]ShopButton)
	for _, button := range buttons {
		g.DrawShopButton(button, 205, 90)

	}
}
