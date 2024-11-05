package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gen2brain/raylib-go/raylib"
)

type UIComponent interface {
	Render()
	OnClick()
	Rect() rl.Rectangle
}

type Button struct {
	Rectangle  rl.Rectangle
	Color      rl.Color
	HoverColor rl.Color
	Text       string
	TextColor  rl.Color
	TextSize   int32
	OnClick    func(*Game)
	Active     bool
}

type Dropdown struct {
	Rectangle     rl.Rectangle
	Options       []*Option
	CurrentOption *Option
	Color         rl.Color
	TextColor     rl.Color
	TextSize      int32
	IsOpen        bool
}

type Option struct {
	Text     string
	OnChange func(*Game, *Option)
}

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
type Scene struct {
	Name                string
	Active              bool
	AutoDisable         bool
	DrawScene           func(*Game)
	UpdateScene         func(*Game)
	Buttons             []Button
	skip                bool
	Data                map[string]interface{}
	Camera              rl.Camera2D
	Windows             map[string]*Window
	Menu                *BoardRightClickMenu
	RenderMenu          bool
	Messages            []Message
	KeyBindingFunctions map[string]func(*Game)
	KeyBindings         map[string]*KeyBinding
	Components          []UIComponent
}

type KeyBinding struct {
	Current      int32
	Default      int32
	Name         string
	FunctionName string
	Scene        string
	Configurable bool
	OnPress      func(*Game)
}

type KeyBindingJSON struct {
	Default      int32  `json:"default"`
	Name         string `json:"name"`
	FunctionName string `json:"functionName"`
	Scene        string `json:"scene"`
	Configurable bool   `json:"configurable"`
}

type Message struct {
	Text  string
	Vec   rl.Vector2
	Timer int
}

type Window struct {
	Name       string
	DrawWindow func(*Game, *Window)
	Display    bool
	Buttons    []Button
}

type Tile struct {
	Texture   rl.Texture2D
	TileFrame rl.Rectangle
	Color     rl.Color
}

type BoardRightClickMenu struct {
	Rectangle   rl.Rectangle
	BoardSquare *BoardSquare
	Items       []BoardMenuItem
}

type BoardMenuItem struct {
	Rectangle       rl.Rectangle
	Text            string
	OnClick         func(*Game)
	CheckIsDisabled func(*Game, *BoardSquare) bool
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

func (g *Game) CloseButton(x, y float32, onClick func(*Game)) Button {
	closeButton := g.Button("X", x, y, onClick)
	closeButton.Rectangle.Width = 40
	return closeButton
}

func (g *Game) Button(text string, x, y float32, onClick func(*Game)) Button {
	return Button{
		Rectangle:  rl.NewRectangle(x, y, 150, 40),
		Color:      rl.SkyBlue,
		HoverColor: rl.LightGray,
		Text:       text,
		TextColor:  rl.Black,
		OnClick:    onClick,
		Active:     true,
	}
}

func (g *Game) DrawButton(button Button) {
	var boxColor rl.Color
	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
		if button.HoverColor == rl.Blank {
			button.HoverColor = button.Color
		}
		boxColor = button.HoverColor
	} else {
		boxColor = button.Color
	}
	rl.DrawRectangle(button.Rectangle.ToInt32().X, button.Rectangle.ToInt32().Y, button.Rectangle.ToInt32().Width, button.Rectangle.ToInt32().Height, boxColor)
	textSize := button.TextSize
	if textSize == 0 {
		textSize = int32(button.Rectangle.Height - 15)
	}

	rl.DrawText(
		button.Text,
		button.Rectangle.ToInt32().X+5,
		button.Rectangle.ToInt32().Y+5,
		textSize,
		button.TextColor,
	)
}

func (g *Game) DrawButtons(buttons []Button) {
	for _, button := range buttons {
		g.DrawButton(button)
	}

}

func (g *Game) WasButtonClicked(button *Button) bool {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !g.ScreenSkip {
		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
			return true
		}
	}
	return false
}

func (g *Game) Draw() {

	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		scene.DrawScene(g)
		for _, button := range scene.Buttons {
			if button.Active {
				g.DrawButton(button)
			}
		}

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
						g.DrawButton(button)
					}
				}

				for _, component := range scene.Components {
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
		for _, button := range scene.Buttons {
			if g.WasButtonClicked(&button) {
				button.OnClick(g)
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
				for _, button := range window.Buttons {
					if g.WasButtonClicked(&button) {
						button.OnClick(g)
					}
				}
			}

		}

		scene.UpdateScene(g)
	}

	if g.ScreenSkip {
		if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
			g.ScreenSkip = false

			//			log.Printf("remove screen skip: mouse down %v", rl.IsMouseButtonPressed(rl.MouseLeftButton))
		}
	}

	if g.ButtonSkip != 0 {
		if rl.IsKeyUp(g.ButtonSkip) {
			g.ButtonSkip = 0
		}
	}
}

// tiles

func (g *Game) GetBoardCoordAtPoint(vec rl.Vector2) BoardCoord {

	scene := g.Scenes["Board"]

	//	mousePosition := rl.GetMousePosition()
	X := int((vec.X + scene.Camera.Target.X) / scene.Camera.Zoom / float32(TILE_WIDTH))
	Y := int((vec.Y + scene.Camera.Target.Y) / scene.Camera.Zoom / float32(TILE_HEIGHT))
	return BoardCoord{
		Row:    X,
		Column: Y,
	}
}

// window handling
func (g *Game) DisableAllWindows(windows map[string]*Window) {
	for _, window := range windows {
		window.Display = false
	}
}

func (g *Game) ActivateWindow(windows map[string]*Window, window *Window) {
	g.ScreenSkip = true
	if window.Display {
		g.DisableAllWindows(windows)
	} else {
		g.DisableAllWindows(windows)
		window.Display = true
	}
}

// menus

func (g *Game) DrawContextMenu(scene *Scene) {
	if !scene.RenderMenu {
		return
	}

	mousePosition := rl.GetMousePosition()

	var color rl.Color
	var textColor rl.Color

	x := scene.Menu.Rectangle.X
	y := scene.Menu.Rectangle.Y

	square := scene.Menu.BoardSquare

	for _, item := range scene.Menu.Items {
		rec := item.Rectangle
		rec.X = x
		rec.Y = y
		if rl.CheckCollisionPointRec(mousePosition, rec) {
			color = rl.Gray
		} else {
			color = rl.White
		}
		if !item.CheckIsDisabled(g, square) {
			textColor = rl.Black
		} else {
			textColor = rl.LightGray

		}
		rl.DrawRectangleRec(rec, color)
		rl.DrawText(item.Text, int32(rec.X+5), int32(rec.Y+5), 15, textColor)

		y = rec.Y + rec.Height

		if g.ScreenSkip {
			if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
				g.ScreenSkip = false
			}
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !g.ScreenSkip {
			mousePosition := rl.GetMousePosition()
			if rl.CheckCollisionPointRec(mousePosition, rec) {
				if !item.CheckIsDisabled(g, square) {
					item.OnClick(g)
					scene.RenderMenu = false
				}
			}

		}
	}

}

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
			scene.KeyBindings[binding.Name] = &KeyBinding{
				Current: binding.Default,
				Default: binding.Default,
				OnPress: scene.KeyBindingFunctions[binding.FunctionName],
			}

		}
	}

}

// ui

func DefaultOptionOnChange(g *Game, o *Option) {}

func (dropdown *Dropdown) Render() {
	log.Printf("adfadfad")
	rl.DrawRectangleRec(dropdown.Rectangle, dropdown.Color)
	rl.DrawRectangleLinesEx(dropdown.Rectangle, 1, rl.Black)
	rl.DrawText(
		dropdown.CurrentOption.Text,
		dropdown.Rectangle.ToInt32().X+5,
		dropdown.Rectangle.ToInt32().Y+5,
		dropdown.TextSize,
		dropdown.TextColor,
	)
	if dropdown.IsOpen {
		for _, option := range dropdown.Options {
			rect := dropdown.Rectangle
			rect.Y += dropdown.Rectangle.Height

			rl.DrawRectangleRec(rect, dropdown.Color)
			rl.DrawRectangleLinesEx(rect, 1, rl.Black)
			rl.DrawText(
				option.Text,
				int32(rect.X+5),
				int32(rect.Y+5),
				dropdown.TextSize,
				dropdown.TextColor,
			)
		}
	}
}

func (dropdown *Dropdown) OnClick() {
	dropdown.IsOpen = !dropdown.IsOpen
}

func (dropdown *Dropdown) Rect() rl.Rectangle {
	if dropdown.IsOpen {
		rect := dropdown.Rectangle
		rect.Height = rect.Height * float32(len(dropdown.Options))
		return rect

	}
	return dropdown.Rectangle
}
