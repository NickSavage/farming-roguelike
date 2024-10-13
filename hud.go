package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
)

func OnClickTestButton(g *Game) {
	g.EndRound()
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
		OnClick:   OnClickTestButton,
	}
	scene.Buttons = append(scene.Buttons, endButton)

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
	sidebarWidth := int32(150)
	rl.DrawRectangle(0, g.screenHeight-height, g.screenWidth, height, rl.Black)
	rl.DrawRectangle(0, 0, sidebarWidth, g.screenHeight-height, rl.Black)

	rl.DrawText(fmt.Sprintf("Money: $%v", g.Run.Money), 30, 30, 20, rl.White)
	g.DrawButtons(scene.Buttons)

}
