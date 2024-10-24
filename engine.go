package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

type Button struct {
	Rectangle rl.Rectangle
	Color     rl.Color
	Text      string
	TextColor rl.Color
	TextSize  int32
	OnClick   func(*Game)
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
	Name        string
	Active      bool
	AutoDisable bool
	DrawScene   func(*Game)
	UpdateScene func(*Game)
	Buttons     []Button
	skip        bool
	Data        map[string]interface{}
	Camera      rl.Camera2D
	Windows     map[string]*Window
	Menu        *BoardRightClickMenu
	RenderMenu  bool
}

type Window struct {
	Name       string
	DrawWindow func(*Game, *Window)
	Display    bool
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
		Rectangle: rl.NewRectangle(x, y, 150, 40),
		Color:     rl.SkyBlue,
		Text:      text,
		TextColor: rl.Black,
		OnClick:   onClick,
	}
}

func (g *Game) DrawButton(button Button) {

	rl.DrawRectangle(button.Rectangle.ToInt32().X, button.Rectangle.ToInt32().Y, button.Rectangle.ToInt32().Width, button.Rectangle.ToInt32().Height, button.Color)
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
	}
	rl.EndDrawing()
}

func (g *Game) Update() {
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		scene.UpdateScene(g)
	}

	if g.ScreenSkip {
		if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
			g.ScreenSkip = false
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
	log.Printf("triffer")
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
