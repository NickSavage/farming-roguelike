package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
)

func OnClickShopWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["ShopWindow"])
}

func OnClickTechWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["TechWindow"])
}

func OnClickOpenEndRoundPage1Window(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["EndRound1"])
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
		OnClick:   OnClickOpenEndRoundPage1Window,
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

	scene.Windows = make(map[string]*Window)
	scene.Windows["ShopWindow"] = &Window{
		Name:       "Shop Window",
		Display:    false,
		DrawWindow: DrawShopWindow,
	}
	scene.Windows["TechWindow"] = &Window{
		Name:       "Tech Window",
		Display:    false,
		DrawWindow: DrawTechnologyWindow,
	}

	scene.Windows["EndRound1"] = &Window{
		Name:       "End Round 1",
		Display:    false,
		DrawWindow: DrawEndRoundWindowPage1,
	}
	scene.Windows["EndRound2"] = &Window{
		Name:       "End Round 2",
		Display:    false,
		DrawWindow: DrawEndRoundWindowPage2,
	}
	scene.Windows["NextEvent"] = &Window{
		Name:       "Next Event",
		Display:    false,
		DrawWindow: DrawNextEventWindow,
	}

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
	for _, window := range scene.Windows {
		if window.Display {
			window.DrawWindow(g, window)
		}
	}
}

func DrawEndRoundWindowPage1(g *Game, window *Window) {

	windowRect := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(windowRect, rl.White)
	rl.DrawRectangleLinesEx(windowRect, 5, rl.Black)

	rl.DrawText("Income", int32(windowRect.X+5), int32(windowRect.Y+5), 30, rl.Black)
	var totalEarned float32 = 0

	var x, y int32
	for i, tech := range g.Run.Technology {
		x = int32(windowRect.X + 10)
		y = int32(windowRect.Y + 50 + float32(i*30))
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
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["EndRound2"])
	}
}

func DrawEndRoundWindowPage2(g *Game, win *Window) {

	windowRect := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(windowRect, rl.White)
	rl.DrawRectangleLinesEx(windowRect, 5, rl.Black)

	rl.DrawText("Investments", int32(windowRect.X+5), int32(windowRect.Y+5), 30, rl.Black)

	var actions float32 = float32(g.Run.RoundActions)

	var x, y int32
	for i, tech := range g.Run.Technology {
		x = int32(windowRect.X + 10)
		y = int32(windowRect.Y + 50 + float32(i*30))
		nextSeason := tech.RoundHandler[tech.RoundHandlerIndex]

		actions -= nextSeason.CostActions
		text := fmt.Sprintf(
			"%s: -%v actions -$%v money",
			tech.Name,
			nextSeason.CostActions,
			nextSeason.CostMoney,
		)
		rl.DrawText(text, x, y, 20, rl.Red)

	}
	text := fmt.Sprintf("Actions next season: %v", actions)
	rl.DrawText(text, x, y+30, 20, rl.Red)

	button := g.Scenes["HUD"].Data["EndRoundConfirmButton"].(Button)
	button.Rectangle.X = 500
	button.Rectangle.Y = 500

	g.DrawButton(button)

	previousButton := g.Button("Previous")
	previousButton.Rectangle.X = 300
	previousButton.Rectangle.Y = 500
	g.DrawButton(previousButton)
	if g.WasButtonClicked(&previousButton) {
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["EndRound1"])
	}

	if g.WasButtonClicked(&button) {
		OnClickEndRound(g)
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["NextEvent"])
	}
}

func DrawNextEventWindow(g *Game, win *Window) {

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
		button.OnClick(g)
		win.Display = false
	}

}

func OnClickConfirmNextEvent(g *Game) {

}
