package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
)

func OnClickShopWindowButton(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["Board"]
	// log.Printf("hey %v", scene)
	g.ActivateWindow(scene.Windows, scene.Windows["ShopWindow"])
}

func OnClickOpenEndRoundPage1Window(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["Board"]
	PreEndRound(g)
	g.ActivateWindow(scene.Windows, scene.Windows["EndRound1"])
}

func OnClickOpenMarketWindow(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["Board"]

	components := make([]engine.UIComponent, 0)
	products := g.GetProductNames()

	blank := engine.NewBlankComponent()
	blank.SelectDirections.Up = len(products)
	blank.SelectDirections.Down = 1
	components = append(components, &blank)

	x := float32(230 + 150*5)
	y := float32(50 + 80)
	for i, product := range products {
		rect := rl.NewRectangle(x, y+float32(i*30), 50, 25)
		button := g.NewSellButton(rect, g.Run.Products[product])
		button.SelectDirections.Up = i
		button.SelectDirections.Down = i + 2
		if i == len(products)-1 {
			button.SelectDirections.Down = 1
		}
		log.Printf("up %v down %v", button.SelectDirections.Up, button.SelectDirections.Down)
		components = append(components, &button)

	}

	scene.Windows["Prices"].Components = components
	g.ActivateWindow(scene.Windows, scene.Windows["Prices"])
	log.Printf("components %v directions %v", scene.Windows["Prices"].Components)
	for _, component := range components {
		log.Printf("component %v", component)
	}
	log.Printf("?")
}

func OnClickOpenSettings(gi engine.GameInterface) {

	g := gi.(*Game)
	g.Scenes["Settings"].Data["Return"] = "Board"
	g.ActivateScene("Settings")

}

func OnClickEndRoundConfirmButton(gi engine.GameInterface) {

	g := gi.(*Game)
	OnClickEndRound(g)
	g.ActivateWindow(g.Scenes["Board"].Windows, g.Scenes["Board"].Windows["NextEvent"])
}

func CloseAllWindows(gi engine.GameInterface) {
	g := gi.(*Game)
	scene := g.Scenes["Board"]
	found := false
	for _, window := range scene.Windows {
		if window.Display {
			log.Printf("found %v", window)
			found = true
		}
		window.Display = false
	}
	if !found {
		log.Printf("activate")
		OnClickOpenSettings(g)
		g.Scenes["Settings"].Data["Return"] = "Board"
		g.ActivateScene("Settings")

	}

}

func (g *Game) InitHUD() {
	scene := g.Scenes["Board"]
	log.Printf("init hud")

	g.SidebarWidth = int32(200)

	scene.Windows = make(map[string]*engine.Window)
	scene.Windows["ShopWindow"] = &engine.Window{
		Name:       "Shop Window",
		Display:    false,
		DrawWindow: DrawShopWindow,
		Rectangle: rl.Rectangle{
			X:      float32((g.screenWidth-g.SidebarWidth-890)/2 + g.SidebarWidth),
			Y:      55,
			Width:  890,
			Height: 500,
		},
	}

	scene.Windows["EndRound1"] = &engine.Window{
		Name:       "End Round 1",
		Display:    false,
		DrawWindow: DrawEndRoundWindowPage1,
		Components: make([]engine.UIComponent, 0),
	}
	blank := engine.NewBlankComponent()
	blank.SelectDirections.Right = 1
	blank.SelectDirections.Left = 1
	nextButton := g.NewButton(
		"Next Page",
		rl.NewRectangle(500, 500, 150, 40),
		OnClickEndRoundConfirmButton,
	)
	scene.Windows["EndRound1"].Components = append(scene.Components, &blank)
	scene.Windows["EndRound1"].Components = append(scene.Components, &nextButton)
	scene.Windows["NextEvent"] = &engine.Window{
		Name:       "Next Event",
		Display:    false,
		DrawWindow: DrawNextEventWindow,
	}
	scene.Windows["Prices"] = &engine.Window{
		Name:       "Prices",
		Display:    false,
		DrawWindow: DrawMarketWindow,
	}
	scene.Windows["GameOver"] = &engine.Window{
		Name:       "Game Over",
		Display:    false,
		DrawWindow: DrawGameOverWindow,
	}

	scene.KeyBindingFunctions = make(map[string]func(engine.GameInterface))
	scene.KeyBindingFunctions["CloseAllWindows"] = CloseAllWindows
	//	scene.KeyBindingFunctions["OpenShop"] = OnClickShopWindowButton

	g.LoadSceneShortcuts("Board")
	log.Printf("shorcuts %v", scene.KeyBindings)

	blank = engine.NewBlankComponent()
	blank.SelectDirections.Up = 4
	blank.SelectDirections.Down = len(scene.Components) + 1

	scene.Components = append(scene.Components, &blank)
	// 1
	shopButton := g.NewButton(
		"Shop",
		rl.NewRectangle(10, 240, 150, 40),
		OnClickShopWindowButton,
	)
	shopButton.SelectDirections.Up = 4
	shopButton.SelectDirections.Down = len(scene.Components) + 1
	shopButton.SelectDirections.Right = 5
	scene.Components = append(scene.Components, &shopButton)
	// 2
	priceButton := g.NewButton(
		"Market",
		rl.NewRectangle(10, 290, 150, 40),
		OnClickOpenMarketWindow,
	)
	priceButton.SelectDirections.Up = len(scene.Components) - 1
	priceButton.SelectDirections.Down = len(scene.Components) + 1
	priceButton.SelectDirections.Right = 5
	scene.Components = append(scene.Components, &priceButton)
	// 3
	viewEndRoundButton := g.NewButton(
		"End Round",
		rl.NewRectangle(10, 340, 150, 40),
		OnClickOpenEndRoundPage1Window,
	)
	viewEndRoundButton.SelectDirections.Up = len(scene.Components) - 1
	viewEndRoundButton.SelectDirections.Down = len(scene.Components) + 1
	viewEndRoundButton.SelectDirections.Right = 5

	scene.Components = append(scene.Components, &viewEndRoundButton)
	// 4
	settingsButton := g.NewButton(
		"Settings",
		rl.NewRectangle(10, float32(g.screenHeight)-60, 150, 40),
		OnClickOpenSettings,
	)
	settingsButton.SelectDirections.Up = len(scene.Components) - 1
	settingsButton.SelectDirections.Down = 1
	settingsButton.SelectDirections.Right = 5
	scene.Components = append(scene.Components, &settingsButton)
}

func UpdateHUD(g *Game) {
	scene := g.Scenes["Board"]

	if g.GameOverTriggered {
		g.ActivateWindow(scene.Windows, scene.Windows["GameOver"])
		g.GameOverTriggered = false
	}

}

func DrawHUD(g *Game) {
	scene := g.Scenes["Board"]
	height := int32(150)
	//	rl.DrawRectangle(0, g.screenHeight-height, g.screenWidth, height, rl.Black)

	DrawSidebar(g)

	if g.Data["Message"].(string) != "" {
		rl.DrawText(g.Data["Message"].(string), 205, g.screenHeight-height+15, 20, rl.White)
		if g.Data["MessageCounter"].(int32) == g.Seconds {
			g.Data["Message"] = ""
			g.Data["MessageCounter"] = 0
		}

	}
	open := false
	for _, window := range scene.Windows {
		if window.Display {
			window.DrawWindow(g, window)
			open = true
		}
	}
	g.WindowOpen = open
}

func DrawSidebar(g *Game) {

	//	scene := g.Scenes["Board"]
	rl.DrawRectangle(0, 0, g.SidebarWidth, g.screenHeight, rl.Black)

	rl.DrawText(
		fmt.Sprintf("Actions: %v/%v", g.Run.ActionsRemaining, g.Run.ActionsMaximum),
		30, 30, 20, rl.White,
	)
	rl.DrawText(fmt.Sprintf("Round: %v", g.Run.CurrentRound), 30, 50, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Required Money: $%v", g.Run.MoneyRequirement), 30, 70, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Money: $%v", g.Run.Money), 30, 90, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Productivity: %v", g.Run.Productivity), 30, 110, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Season: %v", g.Run.CurrentSeason.String()), 30, 130, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Yield: %v", g.Run.Yield), 30, 150, 20, rl.White)

}

func DrawEndRoundWindowPage1(gi engine.GameInterface, win *engine.Window) {
	log.Printf("win %v", win.Components)

	g := gi.(*Game)
	windowRect := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(windowRect, rl.White)
	rl.DrawRectangleLinesEx(windowRect, 5, rl.Black)

	rl.DrawText("Production", int32(windowRect.X+5), int32(windowRect.Y+5), 30, rl.Black)
	var subtotal float32 = 0
	var columnOffset int32 = 150

	var x, y int32
	var i int
	productNames := g.GetProductNames()

	for _, product := range productNames {
		product := g.Run.Products[product]

		x = int32(windowRect.X + 10)
		y = int32(windowRect.Y + 50 + float32(i*30))
		units := product.Quantity
		price := product.Price
		subtotal += units * price

		text := fmt.Sprintf("$%v (%v units at $%v each)", units*price, units, price)
		rl.DrawText(string(product.Type), x, y, 20, rl.Black)
		rl.DrawText(text, x+columnOffset, y, 20, rl.Black)

		i += 1
	}

	total := (subtotal * g.Run.Yield)
	yield := total - subtotal
	text := fmt.Sprintf("Subtotal: $%v", subtotal)
	rl.DrawText(text, x, y+30, 20, rl.Black)
	text = fmt.Sprintf("Yield: %v (%%%v)", yield, g.Run.Yield)
	rl.DrawText(text, x, y+50, 20, rl.Black)
	text = fmt.Sprintf("Total: %v", total)
	rl.DrawText(text, x, y+70, 20, rl.Black)
}

func (g *Game) HandleChooseEvent(event Event) {
	g.ActivateWindow(g.Scenes["Board"].Windows, g.Scenes["Board"].Windows["NextEvent"])
	g.ApplyEvent(event)
	g.ApplyPriceChanges(event)
	g.ScreenSkip = true

	//log.Printf("apply screen skip: mouse down %v", rl.IsMouseButtonPressed(rl.MouseLeftButton))
}

func DrawNextEventWindow(gi engine.GameInterface, win *engine.Window) {
	window := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 2, rl.Black)

}

func OnClickConfirmNextEvent(g *Game) {
	g.ActivateWindow(g.Scenes["Board"].Windows, g.Scenes["Board"].Windows["NextEvent"])
}

func DrawMarketWindow(gi engine.GameInterface, win *engine.Window) {
	g := gi.(*Game)
	//	scene := g.Scenes["Board"]

	window := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	rl.DrawText("Market Prices", 225, 60, 30, rl.Black)

	var i, x, y int32
	var columnOffset int32 = 150
	// var sellAllConfirm string = scene.Data["SellAllConfirm"].(string)

	products := g.GetProductNames()
	rl.DrawText("Products", int32(window.X+10), int32(window.Y+50), 25, rl.Black)
	rl.DrawText("Quantity", int32(window.X+10)+columnOffset, int32(window.Y+50), 25, rl.Black)
	rl.DrawText("Yield", int32(window.X+10)+columnOffset*2, int32(window.Y+50), 25, rl.Black)
	rl.DrawText("Spot Price", int32(window.X+10)+columnOffset*3, int32(window.Y+50), 25, rl.Black)
	rl.DrawText("Total Earned", int32(window.X+10)+columnOffset*4, int32(window.Y+50), 25, rl.Black)
	for _, product := range products {
		productName := string(product)
		x = int32(window.X + 10)
		y = int32(window.Y + 80 + float32(i*30))
		rl.DrawText(productName, x, y, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%v", g.Run.Products[product].Quantity), x+columnOffset, y, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%v", g.Run.Products[product].Yield), x+columnOffset*2, y, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%v", g.Run.Products[product].Price), x+columnOffset*3, y, 20, rl.Black)
		value := g.Run.Products[product].TotalEarned
		rl.DrawText(fmt.Sprintf("%v", value), x+columnOffset*4, y, 20, rl.Black)

		if g.Run.Products[product].Quantity == 0 {
			i += 1
			continue
		}
		i += 1

	}
	g.DrawMessages()
}

func DrawGameOverWindow(gi engine.GameInterface, win *engine.Window) {
	//	g := gi.(*Game)
	window := rl.NewRectangle(220, 50, 500, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	rl.DrawText("Game Over!", 225, 60, 30, rl.Black)
}
