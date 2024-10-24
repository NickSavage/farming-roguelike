package main

// func (g *Game) HandleLeftClick() {
// 	// todo: shouldn't work if a window is open

// 	scene := g.Scenes["Board"]

// 	if g.WindowOpen {
// 		return
// 	}
// 	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
// 		mousePosition := rl.GetMousePosition()

// 		if !CheckVecVisible(mousePosition) {
// 			return
// 		}

// 		grid := scene.Data["Grid"].([][]BoardSquare)
// 		x := int((mousePosition.X + scene.Camera.Target.X) / scene.Camera.Zoom / float32(TILE_WIDTH))
// 		y := int((mousePosition.Y + scene.Camera.Target.Y) / scene.Camera.Zoom / float32(TILE_HEIGHT))
// 		if scene.RenderMenu {
// 			// TODO build a rect of the menu and check if the click is within
// 			if scene.Menu.BoardSquare.Row == grid[x][y].Row ||
// 				scene.Menu.BoardSquare.Column == grid[x][y].Column {
// 				return
// 			}
// 		}
// 		if !(grid[x][y].IsTechnology || grid[x][y].IsTree) {

// 			scene.RenderMenu = false
// 			return
// 		}
// 		menu := &BoardRightClickMenu{
// 			Rectangle: rl.Rectangle{
// 				X:      mousePosition.X,
// 				Y:      mousePosition.Y,
// 				Height: 100,
// 				Width:  100,
// 			},
// 			BoardSquare: &grid[x][y],
// 			Items:       make([]BoardMenuItem, 0),
// 		}

// 		if grid[x][y].IsTree {
// 			menu.Items = TreeMenuItems()
// 		}
// 		g.ScreenSkip = true
// 		scene.RenderMenu = true
// 		scene.Menu = menu

// 		return
// 	}
// 	if scene.RenderMenu {
// 		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

// 			mousePosition := rl.GetMousePosition()
// 			if !rl.CheckCollisionPointRec(mousePosition, scene.Menu.Rectangle) {
// 				scene.RenderMenu = false
// 			}

// 		}
// 	}

// }
