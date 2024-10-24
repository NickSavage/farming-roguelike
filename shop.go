package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

func ShopClickChickenCoop(g *Game) {

	tech := g.Technology["ChickenCoop"]
	if !g.Run.CanSpendMoney(tech.CostMoney) {
		return
	}
	g.Scenes["HUD"].Windows["ShopWindow"].Display = false
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.ScreenSkip = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = g.CreateChickenCoopTech()
}

func ShopClickWheatField(g *Game) {

	tech := g.Technology["WheatField"]
	if !g.Run.CanSpendMoney(tech.CostMoney) {
		return
	}
	g.Scenes["HUD"].Windows["ShopWindow"].Display = false
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.ScreenSkip = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = g.CreateWheatTech()

}
func ShopClickWorkstation(g *Game) {
	tech := g.Technology["Workstation"]
	if !g.Run.CanSpendMoney(tech.CostMoney) {
		return
	}
	g.Scenes["HUD"].Windows["ShopWindow"].Display = false
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.ScreenSkip = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = g.CreateWorkstationTech()

}

func (g *Game) DrawShopButton(shopButton ShopButton, x, y float32) {
	textColor := rl.Black
	log.Printf("shop %v", shopButton.Technology)
	if !g.Run.CanSpendMoney(shopButton.Technology.CostMoney) {
		textColor = rl.LightGray

	}
	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  float32(shopButton.Width),
		Height: float32(shopButton.Height),
	}
	rl.DrawRectangleLinesEx(rect, 1, rl.Black)
	rl.DrawRectangleRec(rect, shopButton.BackgroundColor)
	DrawTile(shopButton.Image, x+5, y+2)
	rl.DrawText(shopButton.Technology.Name, int32(x+50), int32(y+2), 20, textColor)
	rl.DrawText(shopButton.Technology.Description, int32(x+50), int32(y+22), 10, textColor)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, rect) {
			shopButton.OnClick(g)
		}
	}
}

func (g *Game) InitShopWindow() {
	log.Printf("init shop")
	tech := g.Technology
	scene := g.Scenes["Board"]
	buttons := []ShopButton{
		ShopButton{
			Width:      400,
			Height:     50,
			Image:      g.Data["ChickenCoopShopTile"].(Tile),
			OnClick:    ShopClickChickenCoop,
			Technology: tech["ChickenCoop"],
		},
		ShopButton{
			Width:      400,
			Height:     50,
			Image:      g.Data["WheatTile"].(Tile),
			OnClick:    ShopClickWheatField,
			Technology: tech["WheatField"],
		},
		ShopButton{
			Width:      400,
			Height:     50,
			Image:      g.Data["WorkstationTile"].(Tile),
			OnClick:    ShopClickWorkstation,
			Technology: tech["Workstation"],
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
	g.DrawShopButton(buttons[2], 205, 200)
	// for _, button := range buttons {
	// 	g.DrawShopButton(button, 205, 90)

	// }
}
