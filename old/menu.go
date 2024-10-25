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

// some camera stuff too
//	scene := g.Scenes["Board"]

//	Camera zoom controls
// scene.Camera.Zoom += float32(rl.GetMouseWheelMove()) * 0.05
// if scene.Camera.Zoom > 1.2 {
// 	scene.Camera.Zoom = 1.2
// } else if scene.Camera.Zoom < 0.8 {
// 	scene.Camera.Zoom = 0.8
// }
// if rl.IsKeyDown(rl.KeyRight) {
// 	scene.Camera.Target.X += 5
// 	if scene.Camera.Target.X > TILE_COLUMNS*TILE_WIDTH-float32(g.screenWidth-50) {
// 		scene.Camera.Target.X = TILE_COLUMNS*TILE_WIDTH - float32(g.screenWidth-50)
// 	}
// }
// if rl.IsKeyDown(rl.KeyLeft) {
// 	scene.Camera.Target.X -= 5
// 	if scene.Camera.Target.X < -200 {
// 		scene.Camera.Target.X = -200
// 	}
// }
// if rl.IsKeyDown(rl.KeyDown) {
// 	scene.Camera.Target.Y += 5
// 	if scene.Camera.Target.Y > TILE_ROWS*TILE_HEIGHT-float32(g.screenHeight-200) {
// 		scene.Camera.Target.Y = TILE_ROWS*TILE_HEIGHT - float32(g.screenHeight-200)
// 	}
// }
// if rl.IsKeyDown(rl.KeyUp) {
// 	scene.Camera.Target.Y -= 5
// 	if scene.Camera.Target.Y < -300 {
// 		scene.Camera.Target.Y = -300
// 	}
// }
//g.SelectTiles()
//	g.HandleLeftClick()

// mousePosition := rl.GetMousePosition()
