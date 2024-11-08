package main

import (
	// "nsavage/farming-roguelike/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// implements UIComponent
type ShopBuildingButton struct {
	g              *Game
	rect           rl.Rectangle
	Technology     *Technology
	ExpandedButton bool // displays title and description
}

func (b ShopBuildingButton) Render() {
	g := b.g
	textColor := rl.Black
	canBuild := true

	backgroundColor := rl.White

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
		if canBuild {
			backgroundColor = rl.LightGray
		}
	}
	x := b.rect.X
	y := b.rect.Y

	rl.DrawRectangleRec(b.rect, backgroundColor)
	rl.DrawRectangleLinesEx(b.rect, 1, rl.Black)
	DrawTile(b.Technology.Tile, x+5, y+2)
	if b.ExpandedButton {
		rl.DrawText(b.Technology.Name, int32(x), int32(y+100), 20, textColor)
		rl.DrawText(b.Technology.Description, int32(x), int32(y+120), 10, textColor)
	}

}
func (b ShopBuildingButton) OnClick() {
	ShopButtonOnClick(b.g, b)
}
func (b ShopBuildingButton) Rect() rl.Rectangle {
	return b.rect
}

func (g *Game) NewShopButton(rect rl.Rectangle, tech *Technology) ShopBuildingButton {
	return ShopBuildingButton{
		g:          g,
		rect:       rect,
		Technology: tech,
	}
}
func (g *Game) NewShopPlantButton(rect rl.Rectangle, tech *Technology) ShopBuildingButton {
	return ShopBuildingButton{
		g:          g,
		rect:       rect,
		Technology: tech,
	}

}
