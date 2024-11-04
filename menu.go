package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
)

func OnClickContinueRun(g *Game) {

	g.InitRun(true)
	g.ActivateScene("Board")
	g.Scenes["HUD"].Active = true
}

func OnClickNewRun(g *Game) {

	g.InitRun(false)
	g.ActivateScene("Board")
	g.Scenes["HUD"].Active = true
}

func OnClickSettings(g *Game) {
	g.ActivateScene("Settings")
	g.Scenes["Settings"].Data["Return"] = "GameMenu"
}
func OnClickStats(g *Game) {}
func OnClickAbout(g *Game) {}

func OnClickExit(g *Game) {
	os.Exit(0)
}

func (g *Game) InitGameMenu() {

	scene := g.Scenes["GameMenu"]
	newButton := Button{
		Rectangle: rl.NewRectangle(
			float32(g.screenWidth)/2-100,
			float32(g.screenHeight)/2,
			200,
			50,
		),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       "New Run",
		TextColor:  rl.Black,
		OnClick:    OnClickNewRun,
		Active:     true,
	}
	scene.Buttons = append(scene.Buttons, newButton)
	continueButton := Button{
		Rectangle: rl.NewRectangle(
			float32(g.screenWidth)/2-100,
			float32(g.screenHeight)/2-60,
			200,
			50,
		),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       "Continue Run",
		TextColor:  rl.Black,
		OnClick:    OnClickContinueRun,
		Active:     g.ExistingSave,
	}
	scene.Buttons = append(scene.Buttons, continueButton)
	settingsButton := Button{
		Rectangle: rl.NewRectangle(
			float32(g.screenWidth)/2-100,
			float32(g.screenHeight)/2+60,
			200,
			50,
		),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       "Settings",
		TextColor:  rl.Black,
		OnClick:    OnClickSettings,
		Active:     true,
	}
	scene.Buttons = append(scene.Buttons, settingsButton)
	statsButton := Button{
		Rectangle: rl.NewRectangle(
			float32(g.screenWidth)/2-100,
			float32(g.screenHeight)/2+120,
			200,
			50,
		),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       "Statistics",
		TextColor:  rl.Black,
		OnClick:    OnClickStats,
		Active:     true,
	}
	scene.Buttons = append(scene.Buttons, statsButton)
	aboutButton := Button{
		Rectangle: rl.NewRectangle(
			float32(g.screenWidth)/2-100,
			float32(g.screenHeight)/2+180,
			200,
			50,
		),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       "About",
		TextColor:  rl.Black,
		OnClick:    OnClickAbout,
		Active:     true,
	}
	scene.Buttons = append(scene.Buttons, aboutButton)
	exitButton := Button{
		Rectangle: rl.NewRectangle(
			float32(g.screenWidth)/2-100,
			float32(g.screenHeight)/2+240,
			200,
			50,
		),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       "Exit",
		TextColor:  rl.Black,
		OnClick:    OnClickExit,
		Active:     true,
	}
	scene.Buttons = append(scene.Buttons, exitButton)

}

func DrawGameMenu(g *Game) {

}

func UpdateGameMenu(g *Game) {
}
