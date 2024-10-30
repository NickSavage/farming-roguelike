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

func (g *Game) DrawShopButton(shopButton ShopButton, x, y float32) {
	textColor := rl.Black
	canBuild := true

	shopButton.BackgroundColor = rl.White
	if !g.Run.CanSpendMoney(shopButton.Technology.CostMoney) ||
		!shopButton.Technology.CanBuild(g, shopButton.Technology) {
		textColor = rl.LightGray
		canBuild = false
	}
	if !g.Run.CanSpendAction(shopButton.Technology.CostActions) {
		textColor = rl.LightGray
		canBuild = false
	}
	_, err := g.GetOpenSpace(shopButton.Technology)
	if err != nil {
		textColor = rl.LightGray
		canBuild = false
	}
	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  float32(shopButton.Width),
		Height: float32(shopButton.Height),
	}

	mousePosition := rl.GetMousePosition()

	if rl.CheckCollisionPointRec(mousePosition, rect) {

		if canBuild {

			shopButton.BackgroundColor = rl.LightGray
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			g.ShopChooseTech(shopButton.Technology)
		}
	}
	rl.DrawRectangleRec(rect, shopButton.BackgroundColor)
	rl.DrawRectangleLinesEx(rect, 1, rl.Black)
	DrawTile(shopButton.Image, x+5, y+2)
	rl.DrawText(shopButton.Technology.Name, int32(x), int32(y+50), 20, textColor)
	rl.DrawText(shopButton.Technology.Description, int32(x), int32(y+70), 10, textColor)
}

func (g *Game) DrawPlantPurchaseButton(shopButton *ShopButton, x, y float32) {

	//	textColor := rl.Black

	shopButton.BackgroundColor = rl.White

	if !g.Run.CanSpendMoney(shopButton.Technology.CostMoney) ||
		!shopButton.Technology.CanBuild(g, shopButton.Technology) {
		shopButton.Image.Color = rl.LightGray
		// textColor = rl.LightGray
	}
	if !g.Run.CanSpendAction(shopButton.Technology.CostActions) {
		shopButton.Image.Color = rl.LightGray
		// textColor = rl.LightGray
	}
	_, err := g.GetOpenSpace(shopButton.Technology)
	if err != nil {
		shopButton.Image.Color = rl.LightGray
		// textColor = rl.LightGray
	}

	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  float32(100),
		Height: float32(100),
	}

	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, rect) {
		shopButton.BackgroundColor = rl.LightGray

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			g.ShopChooseTech(shopButton.Technology)
		}
	}
	rl.DrawRectangleLinesEx(rect, 1, rl.Black)
	rl.DrawRectangleRec(rect, shopButton.BackgroundColor)
	DrawTile(shopButton.Image, x+5, y+5)
}

func (g *Game) InitShopWindow() {
	log.Printf("init shop")
	tech := g.Technology
	scene := g.Scenes["Board"]
	buttons := []ShopButton{
		ShopButton{
			Width:      150,
			Height:     300,
			Image:      g.Data["ChickenCoopShopTile"].(Tile),
			OnClick:    ShopClickChickenCoop,
			Technology: tech["ChickenCoop"],
		},
		ShopButton{
			Width:      150,
			Height:     300,
			Image:      g.Data["WorkstationTile"].(Tile),
			OnClick:    ShopClickWorkstation,
			Technology: tech["Workstation"],
		},
		ShopButton{
			Width:      150,
			Height:     300,
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
	g.DrawShopButton(buttons[0], 215, 90)
	g.DrawShopButton(buttons[1], 375, 90)
	g.DrawShopButton(buttons[2], 535, 90)

	plants := g.Run.CurrentRoundShopPlants

	var x float32 = 215
	for i, plant := range plants {
		log.Printf("draw %v", plant)
		g.DrawPlantPurchaseButton(g.ShopButton(plant), x+float32(i*100), 400)
	}

	// g.DrawPlantPurchaseButton(WheatShopButton(g), 215, 400)
	// g.DrawPlantPurchaseButton(PotatoShopButton(g), 325, 400)
	// g.DrawPlantPurchaseButton(buttons[2], 260, 400)
	// for _, button := range buttons {
	// 	g.DrawShopButton(button, 205, 90)

	// }
}
