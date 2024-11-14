package main

import (
	"fmt"
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
	g := b.g
	textColor := rl.Black
	canBuild := true

	backgroundColor := rl.White
	if b.Purchased {
		backgroundColor = rl.LightGray
	}

	if !g.Run.CanSpendMoney(b.Technology.CostMoney) ||
		!g.CanBuild(b.Technology) {
		textColor = rl.LightGray
		canBuild = false
	}
	if !g.Run.CanSpendAction(b.Technology.CostActions) {
		textColor = rl.LightGray
		canBuild = false
	}
	_, err := g.GetOpenSpace(b.Technology)
	if err != nil {
		textColor = rl.LightGray
		canBuild = false
	}

	mousePosition := rl.GetMousePosition()

	if rl.CheckCollisionPointRec(mousePosition, b.rect) {
		if canBuild && !b.Purchased {
			backgroundColor = rl.LightGray
		}
	}
	if b.IsSelected() {
		backgroundColor = rl.LightGray
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
