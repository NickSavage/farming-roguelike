package main

import (
	"fmt"
	"log"
	"math"
	"nsavage/farming-roguelike/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// implements UIComponent
type ShopBuildingButton struct {
	g                *Game
	rect             rl.Rectangle
	Technology       *Technology
	ExpandedButton   bool // displays title and description
	Position         int
	Purchased        bool // whether its been purchased or not
	SelectDirections engine.SelectDirections
	Selected         bool
	CanBuild         bool //determined in InitShopRoundComponents
}

func (b *ShopBuildingButton) Render() {
	// g := b.g
	textColor := rl.Black

	backgroundColor := rl.White
	if b.Purchased || b.IsSelected() {
		backgroundColor = rl.LightGray
	}

	if !b.CanBuild {
		textColor = rl.LightGray
	}

	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, b.rect) {
		if b.CanBuild && !b.Purchased {
			backgroundColor = rl.LightGray
		}
	}

	x := b.rect.X
	y := b.rect.Y

	rl.DrawRectangleRec(b.rect, backgroundColor)
	rl.DrawRectangleLinesEx(b.rect, 1, rl.Black)
	//	log.Printf("button %v", b)
	if !b.Purchased {
		DrawTile(b.Technology.Tile, x+5, y+2)
		if b.ExpandedButton {
			rl.DrawText(b.Technology.Name, int32(x), int32(y+100), 20, textColor)
			rl.DrawText(b.Technology.Description, int32(x), int32(y+120), 10, textColor)
		}
	}

}
func (b ShopBuildingButton) OnClick() {
	if b.Purchased {
		return
	}
	ShopButtonOnClick(b.g, b)
}
func (b *ShopBuildingButton) Rect() rl.Rectangle {
	return b.rect
}

func (b *ShopBuildingButton) Select()                              { b.Selected = true }
func (b *ShopBuildingButton) Unselect()                            { b.Selected = false }
func (b *ShopBuildingButton) IsSelected() bool                     { return b.Selected }
func (b *ShopBuildingButton) Directions() *engine.SelectDirections { return &b.SelectDirections }

func (g *Game) NewShopButton(rect rl.Rectangle, tech *Technology) ShopBuildingButton {
	return ShopBuildingButton{
		g:                g,
		rect:             rect,
		Technology:       tech,
		SelectDirections: engine.SelectDirections{},
	}
}

// implements UIComponent
type ShopSeedButton struct {
	g                *Game
	rect             rl.Rectangle
	Technology       *Technology
	Position         int
	SelectDirections engine.SelectDirections
	Selected         bool
	CanBuild         bool //determined in InitShopRoundComponents
	ToBeDeleted      bool
}

func (b *ShopSeedButton) Render() {
	// g := b.g
	textColor := rl.Black

	backgroundColor := rl.White
	if b.IsSelected() || b.Technology.ToBeDeleted {
		backgroundColor = rl.LightGray
	}

	if !b.CanBuild {
		textColor = rl.LightGray
	}

	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, b.rect) {
		if b.CanBuild {
			backgroundColor = rl.LightGray
		}
	}

	x := b.rect.X
	y := b.rect.Y

	rl.DrawRectangleRec(b.rect, backgroundColor)
	rl.DrawRectangleLinesEx(b.rect, 1, rl.Black)
	DrawTile(b.Technology.Tile, x+5, y+2)
	rl.DrawText(b.Technology.Name, int32(x), int32(y+100), 20, textColor)

}
func (b *ShopSeedButton) OnClick() {
	ShopSeedButtonOnClick(b.g, b)
}
func (b *ShopSeedButton) Rect() rl.Rectangle {
	return b.rect
}

func (b *ShopSeedButton) Select()                              { b.Selected = true }
func (b *ShopSeedButton) Unselect()                            { b.Selected = false }
func (b *ShopSeedButton) IsSelected() bool                     { return b.Selected }
func (b *ShopSeedButton) Directions() *engine.SelectDirections { return &b.SelectDirections }

func (g *Game) NewShopSeedButton(rect rl.Rectangle, tech *Technology) ShopSeedButton {
	return ShopSeedButton{
		g:                g,
		rect:             rect,
		Technology:       tech,
		SelectDirections: engine.SelectDirections{},
	}
}

// implements UIComponent
type EventButton struct {
	g                *Game
	rect             rl.Rectangle
	Event            *Event
	SelectDirections engine.SelectDirections
	Selected         bool
}

func (b *EventButton) Render() {
	x := b.rect.X
	y := b.rect.Y
	rl.DrawRectangleRec(b.rect, rl.White)
	rl.DrawRectangleLinesEx(b.rect, 5, rl.Black)
	mousePosition := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePosition, b.rect) || b.Selected {
		rl.DrawRectangleRec(b.rect, rl.LightGray)
		rl.DrawRectangleLinesEx(b.rect, 5, rl.Green)
	}

	rl.DrawText(b.Event.Name, int32(x+5), int32(y+10), 30, rl.Black)

	rl.DrawText(b.Event.Description, int32(x+5), int32(y+45), 15, rl.Black)

	rl.DrawText(fmt.Sprintf("%v", b.Event.Severity), int32(x+140), int32(y+400-20), 15, rl.Black)
	for i, effect := range b.Event.Effects {
		adjustedY := y + 400 - 40
		if effect.IsPriceChange {
			newPrice := b.g.Run.Products[effect.ProductImpacted].Price * float32(1+effect.PriceChange)
			newPrice = float32(math.Round(float64(newPrice*100))) / 100

			displayChange := math.Round(float64(effect.PriceChange*100*100)) / 100
			text := fmt.Sprintf("%v: %v (%v%%)", effect.ProductImpacted, newPrice, displayChange)
			rl.DrawText(text, int32(x+5), int32(adjustedY)-int32((i*20)), 20, rl.Black)
		}
	}

}

func (b *EventButton) OnClick() {
	b.g.HandleChooseEvent(*b.Event)
}

func (b *EventButton) Rect() rl.Rectangle {
	return b.rect
}

func (b *EventButton) Select()                              { b.Selected = true }
func (b *EventButton) Unselect()                            { b.Selected = false }
func (b *EventButton) IsSelected() bool                     { return b.Selected }
func (b *EventButton) Directions() *engine.SelectDirections { return &b.SelectDirections }

func (g *Game) NewEventButton(rect rl.Rectangle, event *Event) EventButton {
	return EventButton{
		g:                g,
		rect:             rect,
		Event:            event,
		SelectDirections: engine.SelectDirections{},
	}

}

type SellButton struct {
	g                *Game
	rect             rl.Rectangle
	Product          *Product
	SelectDirections engine.SelectDirections
	Selected         bool
}

func (b *SellButton) Render() {

	if b.Product.Quantity == 0 {
		return
	}
	bgColor := rl.SkyBlue
	if b.Selected {
		bgColor = rl.LightGray
	}
	rl.DrawRectangleRec(b.rect, bgColor)
	rl.DrawRectangleLinesEx(b.rect, 1, rl.Black)
	rl.DrawText(fmt.Sprintf("Sell"), int32(b.rect.X+5), int32(b.rect.Y), 20, rl.Black)

}

func (b *SellButton) OnClick() {

	g := b.g
	scene := g.Scenes["Board"]
	result := g.SellProduct(b.Product)
	g.Run.Money += result

	message := engine.Message{
		Text:  fmt.Sprintf("%v", result),
		Vec:   rl.Vector2{X: b.rect.X, Y: b.rect.Y},
		Timer: 30,
	}
	scene.Messages = append(scene.Messages, message)
}

func (b *SellButton) Rect() rl.Rectangle {
	return b.rect
}

func (b *SellButton) Select()                              { b.Selected = true }
func (b *SellButton) Unselect()                            { b.Selected = false }
func (b *SellButton) IsSelected() bool                     { return b.Selected }
func (b *SellButton) Directions() *engine.SelectDirections { return &b.SelectDirections }

func (g *Game) NewSellButton(rect rl.Rectangle, product *Product) SellButton {
	return SellButton{
		g:                g,
		rect:             rect,
		Product:          product,
		SelectDirections: engine.SelectDirections{},
	}

}

type TechnologySpace struct {
	Game             *Game
	ID               int
	Technology       *Technology
	TechnologyType   TechnologyType
	IsField          bool
	PlantedSeeds     []*Technology
	Row              int
	Column           int
	Width            int // in tiles
	Height           int // in tiles
	IsFilled         bool
	Active           bool // whether the game displays or not
	SelectDirections engine.SelectDirections
	Selected         bool
}

func (space *TechnologySpace) Render() {

	g := space.Game
	scene := space.Game.Scenes["Board"]
	if !space.Active {
		return
	}
	boxColor := rl.Blue
	if space.Selected {
		boxColor = rl.Green
	}

	vec := g.GetVecFromCoords(engine.BoardCoord{Row: space.Row, Column: space.Column})
	x := vec.X
	y := vec.Y
	width := float32(space.Width * TILE_WIDTH)
	height := float32(space.Height * TILE_HEIGHT)
	rect := rl.NewRectangle(x, y, width, height)
	rl.DrawRectangleRec(rect, boxColor)
	if !space.IsFilled {
		return
	}
	mousePosition := rl.GetMousePosition()

	if !g.WindowOpen && rl.CheckCollisionPointRec(mousePosition, rect) {

		space.Technology.Tile.Color = rl.Green
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			result := g.HandleClickTech(space)
			message := engine.Message{
				Text:  result,
				Vec:   rl.Vector2{X: x, Y: y},
				Timer: 30,
			}
			scene.Messages = append(scene.Messages, message)
		}

	} else {
		if space.Technology.ReadyToHarvest {
			space.Technology.Tile.Color = rl.Blue
		} else {
			space.Technology.Tile.Color = rl.White
		}
	}

	if space.Technology.TileFillSpace {
		for i := range space.Width {
			for j := range space.Height {
				DrawTile(
					space.Technology.Tile,
					float32(float32(x)+float32(i*TILE_WIDTH)),
					float32(float32(y)+float32(j*TILE_HEIGHT)),
				)
			}
		}

	} else {
		DrawTile(space.Technology.Tile, float32(x), float32(y))

	}
	if space.IsField {
		for _, seed := range space.PlantedSeeds {
			rl.DrawText(string(seed.ProductType), int32(x), int32(y), 5, rl.Black)

		}
	}

}

func (space *TechnologySpace) OnClick() {
	if space.IsFilled {
		space.Technology.OnClick(space.Game, space.Technology)
	}

}

func (space *TechnologySpace) Rect() rl.Rectangle {

	vec := space.Game.GetVecFromCoords(engine.BoardCoord{Row: space.Row, Column: space.Column})
	return rl.Rectangle{
		X:      vec.X,
		Y:      vec.Y,
		Width:  float32(space.Width * TILE_WIDTH),
		Height: float32(space.Height) * TILE_HEIGHT,
	}

}
func (space *TechnologySpace) Select() {
	log.Printf("selected")
	space.Selected = true

}

func (space *TechnologySpace) Unselect() {
	space.Selected = false
}

func (space *TechnologySpace) IsSelected() bool {
	return space.Selected
}

func (space *TechnologySpace) Directions() *engine.SelectDirections {
	return &space.SelectDirections
}
