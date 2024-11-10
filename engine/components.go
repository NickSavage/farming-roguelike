package engine

import (
	"github.com/gen2brain/raylib-go/raylib"
	"log"
)

type UIComponent interface {
	Render()
	OnClick()
	Rect() rl.Rectangle
	Select()
	Unselect()
	Directions() *SelectDirections
}

type SelectDirections struct {
	Up    int
	Left  int
	Right int
	Down  int
}

type Dropdown struct {
	Rectangle        rl.Rectangle
	Options          []*Option
	CurrentOption    *Option
	Color            rl.Color
	TextColor        rl.Color
	TextSize         int32
	IsOpen           bool
	SelectDirections SelectDirections
}

type Option struct {
	Text     string
	OnChange func(*GameInterface, *Option)
}

type Button struct {
	GameInterface
	Rectangle        rl.Rectangle
	Color            rl.Color
	HoverColor       rl.Color
	Text             string
	TextColor        rl.Color
	TextSize         int32
	Active           bool // whether button is active or disabled
	OnClickFunction  func(GameInterface)
	Selected         bool // whether user has selected the button with a controller
	SelectDirections SelectDirections
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

func (dropdown *Dropdown) Select()                       {}
func (dropdown *Dropdown) Unselect()                     {}
func (dropdown *Dropdown) Directions() *SelectDirections { return &dropdown.SelectDirections }

func (button Button) Render() {
	var boxColor rl.Color
	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
		if button.HoverColor == rl.Blank {
			button.HoverColor = button.Color
		}
		boxColor = button.HoverColor
	} else if button.Selected {
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

func (button Button) OnClick() {
	log.Printf("???")
	if button.OnClickFunction != nil {
		button.OnClickFunction(button.GameInterface)
	}
}

func (button Button) Rect() rl.Rectangle {
	return button.Rectangle
}

func (button *Button) WasButtonClicked() bool {
	// todo fix screenskip
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) { //&& !g.ScreenSkip {
		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
			return true
		}
	}
	return false
}

func (button *Button) Select() {
	button.Selected = true
}

func (button *Button) Unselect() {
	button.Selected = false
}
func (button *Button) Directions() *SelectDirections { return &button.SelectDirections }

type BlankComponent struct {
	SelectDirections SelectDirections
}

func (c *BlankComponent) Render() {

}

func (c *BlankComponent) OnClick() {}

func (c *BlankComponent) Rect() rl.Rectangle {
	return rl.NewRectangle(0, 0, 0, 0)
}

func (c *BlankComponent) Select()   {}
func (c *BlankComponent) Unselect() {}
func (c *BlankComponent) Directions() *SelectDirections {
	return &c.SelectDirections
}

func NewBlankComponent() BlankComponent {
	return BlankComponent{
		SelectDirections: SelectDirections{},
	}
}
