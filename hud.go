package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
)

func (g *Game) HideOtherWindows() {
	scene := g.Scenes["HUD"]
	scene.Data["DisplayShopWindow"] = false
	scene.Data["DisplayTechWindow"] = false
	scene.Data["DisplayEndRoundWindow"] = false
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

func OnClickTechWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	if scene.Data["DisplayTechWindow"].(bool) == true {
		scene.Data["DisplayTechWindow"] = false
	} else {
		g.HideOtherWindows()
		scene.Data["DisplayTechWindow"] = true
	}
}

func OnClickOpenEndRoundWindow(g *Game) {
	scene := g.Scenes["HUD"]
	// warn about remaining actions?

	if scene.Data["DisplayEndRoundWindow"].(bool) == true {
		scene.Data["DisplayEndRoundWindow"] = false
	} else {
		g.HideOtherWindows()
		g.ScreenSkip = true
		scene.Data["DisplayEndRoundWindowPage"] = 1
		scene.Data["DisplayEndRoundWindow"] = true
	}
}

func (g *Game) InitHUD() {
	scene := g.Scenes["HUD"]

	techButton := Button{
		Rectangle: rl.Rectangle{
			X:      10,
			Y:      150,
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
			Y:      200,
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

	viewEndRoundButton := Button{
		Rectangle: rl.Rectangle{
			X:      10,
			Y:      250,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "End Round",
		TextColor: rl.Black,
		OnClick:   OnClickOpenEndRoundWindow,
	}
	scene.Buttons = append(scene.Buttons, viewEndRoundButton)

	scene.Data["EndRoundToPageTwoButton"] = Button{
		Rectangle: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "Next Page",
		TextColor: rl.Black,
	}
	scene.Data["EndRoundConfirmButton"] = Button{
		Rectangle: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "End Round",
		TextColor: rl.Black,
		OnClick:   OnClickEndRound,
	}
	scene.Data["DisplayEndRoundWindow"] = false

	scene.Data["NextEventConfirmButton"] = Button{
		Rectangle: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  150,
			Height: 40,
		},
		Color:     rl.SkyBlue,
		Text:      "Confirm",
		TextColor: rl.Black,
		OnClick:   OnClickConfirmNextEvent,
	}
	scene.Data["DisplayNextEventWindow"] = false
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
	rl.DrawText(fmt.Sprintf("Round: %v", g.Run.CurrentRound), 30, 70, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Season: %v", g.Run.CurrentSeason.String()), 30, 90, 20, rl.White)
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
	if scene.Data["DisplayEndRoundWindow"].(bool) {
		g.DrawEndRoundWindow()
	}
	if scene.Data["DisplayNextEventWindow"].(bool) {
		g.DrawNextEventWindow()
	}

}

func (g *Game) DrawEndRoundWindow() {

	window := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	displayPage := g.Scenes["HUD"].Data["DisplayEndRoundWindowPage"].(int)
	if displayPage == 1 {
		g.DrawEndRoundWindowPage1(window)
	} else {
		g.DrawEndRoundWindowPage2(window)
	}

}

func (g *Game) DrawEndRoundWindowPage1(window rl.Rectangle) {
	if g.ScreenSkip {
		if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
			g.ScreenSkip = false
		}
	}
	var totalEarned float32 = 0

	var x, y int32
	for i, tech := range g.Run.Technology {
		x = int32(window.X + 10)
		y = int32(window.Y + 50 + float32(i*30))
		value := tech.RoundHandler[0].RoundEndValue(g, tech)
		totalEarned += value
		text := fmt.Sprintf("%s: $%v", tech.Name, value)
		rl.DrawText(text, x, y, 20, rl.Black)
	}

	text := fmt.Sprintf("Total: $%v", totalEarned)
	rl.DrawText(text, x, y+30, 20, rl.Black)

	button := g.Scenes["HUD"].Data["EndRoundToPageTwoButton"].(Button)
	button.Rectangle.X = 500
	button.Rectangle.Y = 500

	g.DrawButton(button)
	if g.WasButtonClicked(&button) {
		g.ScreenSkip = true
		g.Scenes["HUD"].Data["DisplayEndRoundWindowPage"] = 2
	}
}

func (g *Game) DrawEndRoundWindowPage2(window rl.Rectangle) {

	if g.ScreenSkip {
		if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
			g.ScreenSkip = false
		}
	}

	rl.DrawText("Investments", int32(window.X+5), int32(window.Y+5), 30, rl.Black)

	button := g.Scenes["HUD"].Data["EndRoundConfirmButton"].(Button)
	button.Rectangle.X = 500
	button.Rectangle.Y = 500

	g.DrawButton(button)
	if g.WasButtonClicked(&button) {
		g.ScreenSkip = true
		g.Scenes["HUD"].Data["DisplayEndRoundWindow"] = false
		button.OnClick(g)
	}
}

func (g *Game) DrawNextEventWindow() {
	//	scene := g.Scenes["HUD"]

	if g.ScreenSkip {
		if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
			g.ScreenSkip = false
		}

	}

	window := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	button := g.Scenes["HUD"].Data["NextEventConfirmButton"].(Button)
	button.Rectangle.X = 500
	button.Rectangle.Y = 500

	g.DrawButton(button)
	rl.DrawText("NEW EVENT", 225, 60, 30, rl.Black)
	rl.DrawText(g.Run.Events[g.Run.CurrentRound].Name, 225, 95, 15, rl.Black)

	if g.WasButtonClicked(&button) {
		g.Scenes["HUD"].Data["DisplayNextEventWindow"] = false
		g.ScreenSkip = true
		button.OnClick(g)
	}

}

func OnClickConfirmNextEvent(g *Game) {

}
