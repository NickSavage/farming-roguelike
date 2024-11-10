package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
	"os"
)

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
	g.Scenes["Settings"].Data["Return"] = "GameMenu"
}
func OnClickStats(gi engine.GameInterface) {

	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["Stats"])
}
func OnClickAbout(gi engine.GameInterface) {

	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}

func OnClickExit(gi engine.GameInterface) {

	// g := gi.(*Game)
	os.Exit(0)
}

func (g *Game) InitGameMenu() {
	log.Printf("init menu")

	scene := g.Scenes["GameMenu"]
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
	scene.Components = append(scene.Components, &newButton)
	continueButton := g.NewButton(
		"Continue Run",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2-60, 200, 50),
		OnClickContinueRun,
	)
	scene.Components = append(scene.Components, &continueButton)

	settingsButton := g.NewButton(
		"Settings",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+60, 200, 50),
		OnClickSettings,
	)
	scene.Components = append(scene.Components, &settingsButton)

	statsButton := g.NewButton(
		"Statistics",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+120, 200, 50),
		OnClickStats,
	)
	scene.Components = append(scene.Components, &statsButton)

	aboutButton := g.NewButton(
		"About",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+180, 200, 50),
		OnClickAbout,
	)
	scene.Components = append(scene.Components, &aboutButton)

	exitButton := g.NewButton(
		"Exit",
		rl.NewRectangle(float32(g.screenWidth)/2-100, float32(g.screenHeight)/2+240, 200, 50),
		OnClickExit,
	)
	scene.Components = append(scene.Components, &exitButton)

	scene.KeyBindingFunctions = make(map[string]func(engine.GameInterface))

	scene.Windows = make(map[string]*engine.Window)
	scene.Windows["Stats"] = &engine.Window{
		Name:       "Stats",
		Display:    false,
		DrawWindow: DrawStatsWindow,
		Buttons:    make([]engine.Button, 1),
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
	scene.Windows["Stats"].Components = append(scene.Windows["Stats"].Components, &statsCloseButton)

	scene.Windows["About"] = &engine.Window{
		Name:       "About",
		Display:    false,
		DrawWindow: DrawStatsWindow,
		Buttons:    make([]engine.Button, 1),
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
	scene.Windows["About"].Components = append(scene.Windows["About"].Components, &button)
}

func DrawGameMenu(gi engine.GameInterface) {

}

func UpdateGameMenu(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	if rl.IsKeyPressed(rl.KeyDown) {
		if scene.SelectedComponentIndex == len(scene.Components)-1 {
			scene.SelectedComponentIndex = 0
		} else {
			scene.SelectedComponentIndex += 1
		}
	}
	if scene.SelectedComponentIndex > 0 {
		for i, _ := range scene.Components {
			if i == scene.SelectedComponentIndex {
				scene.Components[i].Select()
			} else {
				scene.Components[i].Unselect()
			}
		}
	}

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

func CloseAboutWindow(gi engine.GameInterface) {

	g := gi.(*Game)
	scene := g.Scenes["GameMenu"]
	g.ActivateWindow(scene.Windows, scene.Windows["About"])
}
