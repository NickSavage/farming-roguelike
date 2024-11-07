package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nsavage/farming-roguelike/engine"
	"os"

	"github.com/gen2brain/raylib-go/raylib"
)

type ShopButton struct {
	X               float32
	Y               float32
	Width           int32
	Height          int32
	Title           string
	Description     string
	Image           Tile
	BackgroundColor rl.Color
	OnClick         func(*Game)
	Technology      *Technology
}

type KeyBindingJSON struct {
	Default      int32  `json:"default"`
	Name         string `json:"name"`
	FunctionName string `json:"functionName"`
	Scene        string `json:"scene"`
	Configurable bool   `json:"configurable"`
}

type Tile struct {
	Texture   rl.Texture2D
	TileFrame rl.Rectangle
	Color     rl.Color
}

func InitEngine() {
	log.Printf("hello world")
}

func DrawTile(t Tile, x, y float32) {

	rl.DrawTextureRec(
		t.Texture,
		t.TileFrame,
		rl.Vector2{
			X: x,
			Y: y,
		},
		t.Color,
	)

}

func (g *Game) ActivateScene(sceneName string) {
	for key, scene := range g.Scenes {
		if key == sceneName {
			scene.Active = true
		} else if scene.AutoDisable {
			scene.Active = false
		} else {
			// do nothing
		}
		g.Scenes[key] = scene
	}
}

func (g *Game) NewButton(text string, rect rl.Rectangle, onClick func(engine.GameInterface)) engine.Button {
	button := engine.Button{
		GameInterface:   g,
		Rectangle:       rect,
		Color:           rl.SkyBlue,
		HoverColor:      rl.LightGray,
		Text:            text,
		TextColor:       rl.Black,
		Active:          true,
		OnClickFunction: onClick,
	}
	return button
}

func (g *Game) Draw() {

	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		scene.DrawScene(g)

		for _, component := range scene.Components {
			component.Render()
		}

		open := false
		for _, window := range scene.Windows {
			if window.Display {
				window.DrawWindow(g, window)
				open = true
				for _, button := range window.Buttons {
					if button.Active {
						button.Render()
					}
				}

				for _, component := range window.Components {
					component.Render()
				}
			}

		}
		g.WindowOpen = open
	}
	rl.EndDrawing()
}

func (g *Game) Update() {
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		for _, binding := range scene.KeyBindings {
			if rl.IsKeyPressed(binding.Current) && !(g.ButtonSkip == binding.Current) {
				binding.OnPress(g)
				g.ButtonSkip = binding.Current
			}

		}

		for _, component := range scene.Components {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !g.ScreenSkip {
				mousePosition := rl.GetMousePosition()
				if rl.CheckCollisionPointRec(mousePosition, component.Rect()) {
					component.OnClick()
				}
			}
		}
		for _, window := range scene.Windows {
			if window.Display {
				for _, component := range window.Components {
					if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !g.ScreenSkip {
						mousePosition := rl.GetMousePosition()
						if rl.CheckCollisionPointRec(mousePosition, component.Rect()) {
							component.OnClick()
						}
					}
				}
			}

		}

		scene.UpdateScene(g)
	}

	if g.ScreenSkip {
		if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
			g.ScreenSkip = false
		}
	}

	if g.ButtonSkip != 0 {
		if rl.IsKeyUp(g.ButtonSkip) {
			g.ButtonSkip = 0
		}
	}
}

// tiles

func (g *Game) GetBoardCoordAtPoint(vec rl.Vector2) engine.BoardCoord {

	scene := g.Scenes["Board"]

	//	mousePosition := rl.GetMousePosition()
	X := int((vec.X + scene.Camera.Target.X) / scene.Camera.Zoom / float32(TILE_WIDTH))
	Y := int((vec.Y + scene.Camera.Target.Y) / scene.Camera.Zoom / float32(TILE_HEIGHT))
	return engine.BoardCoord{
		Row:    X,
		Column: Y,
	}
}

// window handling
func (g *Game) DisableAllWindows(windows map[string]*engine.Window) {
	for _, window := range windows {
		window.Display = false
	}
}

func (g *Game) ActivateWindow(windows map[string]*engine.Window, window *engine.Window) {
	g.ScreenSkip = true
	if window.Display {
		g.DisableAllWindows(windows)
	} else {
		g.DisableAllWindows(windows)
		window.Display = true
	}
}

// menus

// save files

func SaveRun(saveFile SaveFile) error {
	json, err := json.Marshal(saveFile)
	if err != nil {
		log.Printf("save marshalling error %v", err)
		return err
	}
	err = os.WriteFile("save.json", json, 0644)
	log.Printf("save write error %v", err)
	return err
}

func LoadRun() (SaveFile, error) {
	var save SaveFile
	file, err := os.Open("./save.json")
	if err != nil {
		fmt.Println(err)
		return save, err
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)

	err = jsonDecoder.Decode(&save)
	if err != nil {
		fmt.Println(err)
		return save, err
	}

	return save, nil
}

// key bindings

func (g *Game) LoadSceneShortcuts(sceneName string) {
	scene := g.Scenes[sceneName]
	for _, binding := range g.KeyBindingJSONs {
		if binding.Scene == sceneName {
			scene.KeyBindings[binding.Name] = &engine.KeyBinding{
				Current: binding.Default,
				Default: binding.Default,
				OnPress: scene.KeyBindingFunctions[binding.FunctionName],
			}

		}
	}

}

// ui
