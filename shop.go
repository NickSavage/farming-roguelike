package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
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
	g.Scenes["HUD"].Windows["ShopWindow"].Display = false
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

func ShopClickWorkstation(g *Game) {

	tech := g.CreateWorkstationTech()
	g.ShopChooseTech(tech)
}

func ShopClickChickenEggWarmer(g *Game) {
	tech := g.CreateChickenEggWarmer()
	g.ShopChooseTech(tech)
}

func (g *Game) DrawShopButton(shopButton ShopButton, x, y float32) {
	textColor := rl.Black
	log.Printf("shop %v", shopButton.Technology)

	if !g.Run.CanSpendMoney(shopButton.Technology.CostMoney) ||
		!shopButton.Technology.CanBuild(g, shopButton.Technology) {
		textColor = rl.LightGray
	}
	if !g.Run.CanSpendAction(shopButton.Technology.CostActions) {
		textColor = rl.LightGray
	}
	_, err := g.GetOpenSpace(shopButton.Technology)
	if err != nil {
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
			Image:      g.Data["PotatoTile"].(Tile),
			OnClick:    ShopClickPotatoField,
			Technology: tech["PotatoField"],
		},
		ShopButton{
			Width:      400,
			Height:     50,
			Image:      g.Data["WorkstationTile"].(Tile),
			OnClick:    ShopClickWorkstation,
			Technology: tech["Workstation"],
		},
		ShopButton{
			Width:      400,
			Height:     50,
			Image:      g.Data["ChickenEggWarmerShopTile"].(Tile),
			OnClick:    ShopClickChickenEggWarmer,
			Technology: tech["ChickenEggWarmer"],
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
	g.DrawShopButton(buttons[3], 205, 255)
	g.DrawShopButton(buttons[4], 205, 310)
	// for _, button := range buttons {
	// 	g.DrawShopButton(button, 205, 90)

	// }
}
