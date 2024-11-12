package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
	"os"
)

func OnClickAbandonRun(gi *engine.GameInterface) {}

func OnClickContinueRun(gi engine.GameInterface) {

	g := gi.(*Game)
	g.InitRun(true)
	g.ActivateScene("Board")
}

func OnClickNewRun(gi engine.GameInterface) {
	g := gi.(*Game)

	g.InitRun(false)
	g.ActivateScene("Board")
}

func OnClickSettings(gi engine.GameInterface) {

	g := gi.(*Game)
	g.ActivateScene("Settings")
	g.Scenes["Settings"].SelectedComponentIndex = 0
	g.Scenes["Settings"].Data["Return"] = "GameMenu"
}
func OnClickStats(gi engine.GameInterface) {

	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	scene.Windows["Stats"].SelectedComponentIndex = 0
	g.ActivateWindow(scene.Windows, scene.Windows["Stats"])
}
func OnClickAbout(gi engine.GameInterface) {

	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	scene.Windows["About"].SelectedComponentIndex = 0
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}

func OnClickExit(gi engine.GameInterface) {

	// g := gi.(*Game)
	os.Exit(0)
}

func (g *Game) InitGameMenu() {
	log.Printf("init menu")

	// 0
	scene := g.Scenes["GameMenu"]
	components := make([]engine.UIComponent, 0)

	blank := engine.NewBlankComponent()
	if g.ActiveRun {
		blank.SelectDirections.Up = 5
	} else {
		blank.SelectDirections.Up = 6
	}
	blank.SelectDirections.Down = len(components) + 1

	components = append(components, &blank)
	if g.ActiveRun {
		abandonButton := g.NewButton(
			"Abandon Run",
			rl.NewRectangle(
				float32(g.screenWidth)/2-100,
				float32(g.screenHeight)/2,
				200,
				50,
			),
			OnClickContinueRun,
		)
		abandonButton.SelectDirections.Up = 5
		abandonButton.SelectDirections.Down = len(components) + 1
		components = append(components, &abandonButton)

	} else {

		continueButton := g.NewButton(
			"Continue Run",
			rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2-60, 200, 50),
			OnClickContinueRun,
		)
		continueButton.SelectDirections.Up = 6
		continueButton.SelectDirections.Down = len(components) + 1
		components = append(components, &continueButton)

		// 1
		newButton := g.NewButton(
			"New Run",
			rl.NewRectangle(
				float32(g.screenWidth)/2-100,
				float32(g.screenHeight)/2,
				200,
				50,
			),
			OnClickNewRun,
		)
		newButton.SelectDirections.Up = len(components) - 1
		newButton.SelectDirections.Down = len(components) + 1
		components = append(components, &newButton)
	}
	// 2
	settingsButton := g.NewButton(
		"Settings",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+60, 200, 50),
		OnClickSettings,
	)
	settingsButton.SelectDirections.Up = len(components) - 1
	settingsButton.SelectDirections.Down = len(components) + 1
	components = append(components, &settingsButton)
	// 3
	statsButton := g.NewButton(
		"Statistics",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+120, 200, 50),
		OnClickStats,
	)
	statsButton.SelectDirections.Up = len(components) - 1
	statsButton.SelectDirections.Down = len(components) + 1
	components = append(components, &statsButton)

	// 4
	aboutButton := g.NewButton(
		"About",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+180, 200, 50),
		OnClickAbout,
	)
	aboutButton.SelectDirections.Up = len(components) - 1
	aboutButton.SelectDirections.Down = len(components) + 1
	components = append(components, &aboutButton)

	// 5
	exitButton := g.NewButton(
		"Exit",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+240, 200, 50),
		OnClickExit,
	)
	exitButton.SelectDirections.Up = len(components) - 1
	exitButton.SelectDirections.Down = 1
	components = append(components, &exitButton)
	scene.Components = components

	scene.KeyBindingFunctions = make(map[string]func(engine.GameInterface))

	scene.Windows = make(map[string]*engine.Window)
	scene.Windows["Stats"] = &engine.Window{
		Name:       "Stats",
		Display:    false,
		DrawWindow: DrawStatsWindow,
	}
	statsCloseButton := g.NewButton(
		"Close",
		rl.NewRectangle(
			50,
			50,
			200,
			50,
		),
		CloseStatsWindow,
	)

	components = make([]engine.UIComponent, 0)
	blank = engine.NewBlankComponent()
	blank.SelectDirections.Up = 1
	blank.SelectDirections.Down = 1
	components = append(components, &blank)
	components = append(components, &statsCloseButton)
	scene.Windows["Stats"].Components = components

	scene.Windows["About"] = &engine.Window{
		Name:       "About",
		Display:    false,
		DrawWindow: DrawStatsWindow,
	}

	button := g.NewButton(
		"Close",
		rl.NewRectangle(
			50,
			50,
			200,
			50,
		),
		CloseAboutWindow,
	)

	components = make([]engine.UIComponent, 0)
	blank = engine.NewBlankComponent()
	blank.SelectDirections.Up = 1
	blank.SelectDirections.Down = 1
	components = append(components, &blank)
	components = append(components, &button)

	scene.Windows["About"].Components = components
}

func DrawGameMenu(gi engine.GameInterface) {

}

func UpdateGameMenu(gi engine.GameInterface) {
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

func CloseStatsWindow(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["Stats"])
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

func CloseAboutWindow(gi engine.GameInterface) {

	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}

func ToggleMenu(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	if scene.Active {
		if g.ActiveRun {
			g.ActivateScene("Board")
		}
	} else {
		g.InitGameMenu()
		g.ActivateScene("GameMenu")

	}

}
