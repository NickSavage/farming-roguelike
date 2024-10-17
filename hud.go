package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	// "log"
)

func (g *Game) HideOtherWindows() {
	scene := g.Scenes["HUD"]
	scene.Data["DisplayShopWindow"] = false
	scene.Data["DisplayTechWindow"] = false
}

func OnClickShopWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	if scene.Data["DisplayShopWindow"].(bool) == true {
		scene.Data["DisplayShopWindow"] = false
	} else {
		g.HideOtherWindows()
		scene.Data["DisplayShopWindow"] = true
	}

}

func OnClickTestButton(g *Game) {

	g.HideOtherWindows()
	g.Scenes["Board"].Data["PlaceTech"] = true
	g.Scenes["Board"].Data["PlaceTechSkip"] = true
	g.Scenes["Board"].Data["PlaceChosenTech"] = &g.Run.Technology[0]
}

func OnClickTechWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	if scene.Data["DisplayTechWindow"].(bool) == true {
		scene.Data["DisplayTechWindow"] = false
	} else {
		g.HideOtherWindows()
		scene.Data["DisplayTechWindow"] = true
	}
}

func (g *Game) InitHUD() {
	scene := g.Scenes["HUD"]

	endButton := Button{
		Rectangle: rl.Rectangle{
			X:      500,
			Y:      700,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "End Season",
		TextColor: rl.Black,
		OnClick:   OnClickEndRound,
	}
	scene.Buttons = append(scene.Buttons, endButton)
	techButton := Button{
		Rectangle: rl.Rectangle{
			X:      10,
			Y:      100,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "Technology",
		TextColor: rl.Black,
		OnClick:   OnClickTechWindowButton,
	}
	scene.Buttons = append(scene.Buttons, techButton)
	scene.Data["DisplayTechWindow"] = false

	shopButton := Button{
		Rectangle: rl.Rectangle{
			X:      10,
			Y:      150,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "Shop",
		TextColor: rl.Black,
		OnClick:   OnClickShopWindowButton,
	}
	scene.Buttons = append(scene.Buttons, shopButton)
	scene.Data["DisplayShopWindow"] = false

}

func UpdateHUD(g *Game) {
	scene := g.Scenes["HUD"]
	for _, button := range scene.Buttons {
		if g.WasButtonClicked(&button) {
			button.OnClick(g)
		}
	}

}

func DrawHUD(g *Game) {
	scene := g.Scenes["HUD"]
	height := int32(150)
	sidebarWidth := int32(200)
	rl.DrawRectangle(0, g.screenHeight-height, g.screenWidth, height, rl.Black)
	rl.DrawRectangle(0, 0, sidebarWidth, g.screenHeight-height, rl.Black)

	rl.DrawText(
		fmt.Sprintf("Actions: %v/%v", g.Run.RoundActionsRemaining, g.Run.RoundActions),
		30, 30, 20, rl.White,
	)
	rl.DrawText(fmt.Sprintf("Money: $%v", g.Run.Money), 30, 50, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Round: %v", g.Run.Round), 30, 70, 20, rl.White)
	g.DrawButtons(scene.Buttons)

	if g.Data["Message"].(string) != "" {
		rl.DrawText(g.Data["Message"].(string), 205, g.screenHeight-height+15, 20, rl.White)
		if g.Data["MessageCounter"].(int32) == g.Seconds {
			g.Data["Message"] = ""
			g.Data["MessageCounter"] = 0
		}

	}
	if scene.Data["DisplayTechWindow"].(bool) {
		g.DrawTechnologyWindow()
	}
	if scene.Data["DisplayShopWindow"].(bool) {
		g.DrawShopWindow()
	}

}
