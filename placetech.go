package main

// import (
// 	"errors"
// 	"github.com/gen2brain/raylib-go/raylib"
// 	"log"
// )

// func (g *Game) InitPlaceTech() {

// 	scene := g.Scenes["Board"]
// 	scene.Data["PlaceTechCancelButton"] = ShopButton{
// 		Width:           200,
// 		Height:          40,
// 		Title:           "Cancel Placement",
// 		OnClick:         OnClickCancelTechPlacement,
// 		BackgroundColor: rl.SkyBlue,
// 	}

// }

// // placing tech

// func OnClickCancelTechPlacement(g *Game) {
// 	scene := g.Scenes["Board"]
// 	scene.Data["PlaceTech"] = false

// }

// func (g *Game) DrawPlaceTech() {
// 	scene := g.Scenes["Board"]
// 	if scene.Data["PlaceTech"] == nil || !scene.Data["PlaceTech"].(bool) {
// 		return
// 	}

// 	if g.ScreenSkip {
// 		g.ScreenSkip = false
// 		return
// 	}
// 	chosenTech := scene.Data["PlaceChosenTech"].(*Technology)
// 	mousePosition := rl.GetMousePosition()

// 	cancelButton := scene.Data["PlaceTechCancelButton"].(ShopButton)
// 	//	g.DrawShopButton(cancelButton, 200, 50)
// 	if rl.CheckCollisionPointRec(mousePosition, rl.Rectangle{
// 		X:      200,
// 		Y:      50,
// 		Width:  float32(cancelButton.Width),
// 		Height: float32(cancelButton.Height),
// 	}) {
// 		// don't display placement if you're over the cancel button
// 		return
// 	}
// 	mouseTileX := float32(mousePosition.X) - (chosenTech.Square.Tile.TileFrame.Width / 2)
// 	mouseTileY := float32(mousePosition.Y) - (chosenTech.Square.Tile.TileFrame.Height / 2)
// 	if g.CheckTilesOccupied(chosenTech.Square, mousePosition.X, mousePosition.Y) {
// 		occupiedTile := chosenTech.Square.Tile
// 		occupiedTile.Color = rl.Red
// 		rl.DrawRectangle(
// 			int32(mouseTileX),
// 			int32(mouseTileY),
// 			int32(chosenTech.Square.Width*TILE_WIDTH),
// 			int32(chosenTech.Square.Height*TILE_HEIGHT),
// 			rl.Red,
// 		)
// 		DrawTile(
// 			occupiedTile,
// 			mouseTileX,
// 			mouseTileY,
// 		)

// 	} else {
// 		rl.DrawRectangle(
// 			int32(mouseTileX),
// 			int32(mouseTileY),
// 			int32(chosenTech.Square.Width*TILE_WIDTH),
// 			int32(chosenTech.Square.Height*TILE_HEIGHT),
// 			rl.Green,
// 		)
// 		DrawTile(
// 			chosenTech.Square.Tile,
// 			mouseTileX,
// 			mouseTileY,
// 		)
// 		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

// 			log.Printf("draw place tech")
// 			scene.Data["PlaceTech"] = false
// 			g.PlaceTech(
// 				chosenTech,
// 				mouseTileX,
// 				mouseTileY,
// 			)
// 		}

// 	}

// }
