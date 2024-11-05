package engine

import (
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
)

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
	OnChange func(*GameInterface, *Option)
}

type Button struct {
	Rectangle       rl.Rectangle
	Color           rl.Color
	HoverColor      rl.Color
	Text            string
	TextColor       rl.Color
	TextSize        int32
	Active          bool
	OnClickFunction func(GameInterface)
}

func DefaultOptionOnChange(g *GameInterface, o *Option) {}

func (dropdown *Dropdown) Render() {
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

func (button *Button) Render() {
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

func (button *Button) OnClick() {
	// if button.OnClickFunction != nil {
	// 	button.OnClickFunction()
	// }
}

func (button *Button) Rect() rl.Rectangle {
	return button.Rectangle
}
