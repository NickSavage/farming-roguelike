package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
	"os"
)

func OnClickContinueRun(g *Game) {

	g.InitRun(true)
	g.ActivateScene("Board")
}

func OnClickNewRun(g *Game) {

	g.InitRun(false)
	g.ActivateScene("Board")
}

func OnClickSettings(g *Game) {
	g.ActivateScene("Settings")
	g.Scenes["Settings"].Data["Return"] = "GameMenu"
}
func OnClickStats(g *Game) {
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["Stats"])
}
func OnClickAbout(g *Game) {
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}

func OnClickExit(g *Game) {
	os.Exit(0)
}

func (g *Game) InitGameMenu() {
	log.Printf("init menu")

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

	scene.KeyBindingFunctions = make(map[string]func(*Game))

	scene.Windows = make(map[string]*engine.Window)
	scene.Windows["Stats"] = &engine.Window{
		Name:       "Stats",
		Display:    false,
		DrawWindow: DrawStatsWindow,
		Buttons:    make([]engine.Button, 1),
	}
	// scene.Windows["Stats"].Buttons[0] = Button{
	// 	Rectangle: rl.NewRectangle(
	// 		50,
	// 		50,
	// 		200,
	// 		50,
	// 	),
	// 	Color:      rl.SkyBlue,
	// 	HoverColor: rl.LightGray,
	// 	Text:       "Close",
	// 	TextColor:  rl.Black,
	// 	OnClick:    CloseStatsWindow,
	// 	Active:     true,
	// }

	scene.Windows["About"] = &engine.Window{
		Name:       "About",
		Display:    false,
		DrawWindow: DrawStatsWindow,
		Buttons:    make([]engine.Button, 1),
	}
	// scene.Windows["About"].Buttons[0] = Button{
	// 	Rectangle: rl.NewRectangle(
	// 		50,
	// 		50,
	// 		200,
	// 		50,
	// 	),
	// 	Color:      rl.SkyBlue,
	// 	HoverColor: rl.LightGray,
	// 	Text:       "Close",
	// 	TextColor:  rl.Black,
	// 	OnClick:    CloseStatsWindow,
	// 	Active:     true,
	// }
}

func DrawGameMenu(g *Game) {
	shopButton := engine.Button{
		GameInterface: g,
		Rectangle:     rl.NewRectangle(10, 240, 150, 40),
		Color:         rl.SkyBlue,
		HoverColor:    rl.LightGray,
		Text:          "Shop",
		TextColor:     rl.Black,
		Active:        true,
		// OnClickFunction: OnClickTestWindowButton,
		// OnClickFunction: func(gi engine.GameInterface) {
		// 	// Use closure to capture `g` directly if needed
		// 	log.Printf("Game specific run: %+v", g.Scenes)
		// },
	}
	g.Scenes["GameMenu"].Components = append(g.Scenes["GameMenu"].Components, shopButton)
	// shopButton.Render()
}

func UpdateGameMenu(g *Game) {
}

func DrawStatsWindow(gi engine.GameInterface, win *engine.Window) {
	g := gi.(*Game)
	rect := rl.Rectangle{
		X:      20,
		Y:      20,
		Width:  float32(g.screenWidth) - 40,
		Height: float32(g.screenHeight) - 40,
	}
	rl.DrawRectangleRec(rect, rl.White)
	rl.DrawRectangleLinesEx(rect, 5, rl.Black)

}

func CloseStatsWindow(g *Game) {
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}

func DrawAboutWindow(g *Game, win *engine.Window) {
	rect := rl.Rectangle{
		X:      20,
		Y:      20,
		Width:  float32(g.screenWidth) - 40,
		Height: float32(g.screenHeight) - 40,
	}
	rl.DrawRectangleRec(rect, rl.White)
	rl.DrawRectangleLinesEx(rect, 5, rl.Black)

}

func CloseAboutWindow(g *Game) {
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}
