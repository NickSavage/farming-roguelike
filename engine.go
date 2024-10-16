package main

import (
	"github.com/gen2brain/raylib-go/raylib"
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
}

type Tile struct {
	Texture   rl.Texture2D
	TileFrame rl.Rectangle
	Color     rl.Color
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

func (g *Game) DrawButtons(buttons []Button) {
	for _, button := range buttons {
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

}

func (g *Game) DrawShopButton(shopButton ShopButton, x, y float32) {
	rect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  float32(shopButton.Width),
		Height: float32(shopButton.Height),
	}
	rl.DrawRectangleLinesEx(rect, 1, rl.Black)
	rl.DrawRectangleRec(rect, shopButton.BackgroundColor)
	DrawTile(shopButton.Image, x+5, y+2)
	rl.DrawText(shopButton.Title, int32(x+50), int32(y+2), 20, rl.Black)
	rl.DrawText(shopButton.Description, int32(x+50), int32(y+22), 10, rl.Black)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, rect) {
			shopButton.OnClick(g)
		}
	}
}

func (g *Game) WasButtonClicked(button *Button) bool {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
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
}
