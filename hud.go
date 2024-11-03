package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
)

func OnClickNull(g *Game) {}

func OnClickShopWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["ShopWindow"])
}

func OnClickOpenEndRoundPage1Window(g *Game) {
	scene := g.Scenes["HUD"]
	PreEndRound(g)
	g.ActivateWindow(scene.Windows, scene.Windows["EndRound1"])
}

func OnClickOpenMarketWindow(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["Prices"])
}

func OnClickOpenSettings(g *Game) {
	g.ActivateScene("Settings")

}

func OnClickEndRoundPageTwoButton(g *Game) {

}

func OnClickEndRoundConfirmButton(g *Game) {

}

func OpenSellWindow(g *Game, product *Product) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["Sell"])
}

func CloseAllWindows(g *Game) {
	scene := g.Scenes["HUD"]
	for _, window := range scene.Windows {
		window.Display = false
	}

}

func (g *Game) InitHUD() {
	scene := g.Scenes["HUD"]
	log.Printf("init hud")

	g.SidebarWidth = int32(200)

	scene.Windows = make(map[string]*Window)
	scene.Windows["ShopWindow"] = &Window{
		Name:       "Shop Window",
		Display:    false,
		DrawWindow: DrawShopWindow,
	}

	scene.Windows["EndRound1"] = &Window{
		Name:       "End Round 1",
		Display:    false,
		DrawWindow: DrawEndRoundWindowPage1,
	}
	scene.Windows["EndRound2"] = &Window{
		Name:       "End Round 2",
		Display:    false,
		DrawWindow: DrawEndRoundWindowPage2,
	}
	scene.Windows["NextEvent"] = &Window{
		Name:       "Next Event",
		Display:    false,
		DrawWindow: DrawNextEventWindow,
	}
	scene.Windows["Prices"] = &Window{
		Name:       "Prices",
		Display:    false,
		DrawWindow: DrawMarketWindow,
	}
	scene.Windows["Sell"] = &Window{
		Name:       "Sell",
		Display:    false,
		DrawWindow: DrawSellWindow,
	}
	scene.Windows["GameOver"] = &Window{
		Name:       "Game Over",
		Display:    false,
		DrawWindow: DrawGameOverWindow,
	}

	scene.Data["SellAllConfirm"] = ""

	scene.KeyBindings[rl.KeyEscape] = &KeyBinding{
		Current: rl.KeyEscape,
		Default: rl.KeyEscape,
		OnPress: CloseAllWindows,
	}
}

func UpdateHUD(g *Game) {
	scene := g.Scenes["HUD"]
	for _, button := range scene.Buttons {
		if g.WasButtonClicked(&button) {
			button.OnClick(g)
		}
	}
	if g.GameOverTriggered {
		g.ActivateWindow(scene.Windows, scene.Windows["GameOver"])
		g.GameOverTriggered = false
	}

}

func DrawHUD(g *Game) {
	scene := g.Scenes["HUD"]
	height := int32(150)
	//	rl.DrawRectangle(0, g.screenHeight-height, g.screenWidth, height, rl.Black)

	DrawSidebar(g)
	g.DrawButtons(scene.Buttons)

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

	rl.DrawRectangle(0, 0, g.SidebarWidth, g.screenHeight, rl.Black)

	// rl.DrawText(
	// 	fmt.Sprintf("Net Worth: %v", g.Run.CalculateNetWorth()),
	// 	30,
	// 	30,
	// 	20,
	// 	rl.White,
	// )

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

	buttons := []*Button{}
	// techButton := g.Button("Technology", 10, 190, OnClickTechWindowButton)
	// buttons = append(buttons, &techButton)
	shopButton := g.Button("Shop", 10, 240, OnClickShopWindowButton)
	buttons = append(buttons, &shopButton)
	priceButton := g.Button("Market", 10, 290, OnClickOpenMarketWindow)
	buttons = append(buttons, &priceButton)
	viewEndRoundButton := g.Button("End Round", 10, 340, OnClickOpenEndRoundPage1Window)
	buttons = append(buttons, &viewEndRoundButton)
	settingsButton := g.Button("Settings", 10, 390, OnClickOpenSettings)
	buttons = append(buttons, &settingsButton)

	for _, button := range buttons {
		g.DrawButton(*button)
		if g.WasButtonClicked(button) {
			button.OnClick(g)
		}
	}

}

func DrawEndRoundWindowPage1(g *Game, window *Window) {

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

	button := g.Button("Next Page", 500, 500, OnClickEndRoundPageTwoButton)

	g.DrawButton(button)
	// if g.WasButtonClicked(&button) {
	// 	g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["EndRound2"])
	// }
	if g.WasButtonClicked(&button) {
		OnClickEndRound(g)
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["NextEvent"])
	}
}

func DrawEndRoundWindowPage2(g *Game, win *Window) {

	windowRect := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(windowRect, rl.White)
	rl.DrawRectangleLinesEx(windowRect, 5, rl.Black)

	rl.DrawText("Investments", int32(windowRect.X+5), int32(windowRect.Y+5), 30, rl.Black)
	button := g.Button("End Round", 500, 500, OnClickEndRoundConfirmButton)

	g.DrawButton(button)

	previousButton := g.Button("Previous", 300, 500, OnClickOpenEndRoundPage1Window)
	g.DrawButton(previousButton)
	if g.WasButtonClicked(&previousButton) {
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["EndRound1"])
	}
	if g.WasButtonClicked(&button) {
		OnClickEndRound(g)
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["NextEvent"])
	}

}

func (g *Game) HandleChooseEvent(event Event) {
	g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["NextEvent"])
	g.ApplyEvent(event)
	g.ApplyPriceChanges(event)
	g.ScreenSkip = true

	//log.Printf("apply screen skip: mouse down %v", rl.IsMouseButtonPressed(rl.MouseLeftButton))
}

func DrawNextEventWindow(g *Game, win *Window) {

	window := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 2, rl.Black)

	events := g.Run.EventChoices
	var x int32
	var y int32 = 60

	for i, event := range events {
		x = int32(240 + (i * 300))
		rect := rl.NewRectangle(float32(x)+5, float32(y), 300, 400)

		rl.DrawRectangleRec(rect, rl.White)
		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, rect) {
			rl.DrawRectangleRec(rect, rl.LightGray)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				g.HandleChooseEvent(event)
			}
		}
		rl.DrawRectangleLinesEx(rect, 5, rl.Black)

		rl.DrawText(event.Name, x+5, y+10, 30, rl.Black)

		rl.DrawText(event.Description, x+5, y+45, 15, rl.Black)

		for i, effect := range event.Effects {
			if effect.IsPriceChange {
				newPrice := g.Run.Products[effect.ProductImpacted].Price * float32(1+effect.PriceChange)
				newPrice = float32(math.Round(float64(newPrice*100))) / 100

				displayChange := math.Round(float64(effect.PriceChange*100*100)) / 100
				text := fmt.Sprintf("%v: %v (%v%%)", effect.ProductImpacted, newPrice, displayChange)
				rl.DrawText(text, x+5, y+int32(60+(i*20)), 20, rl.Black)
			}
		}

	}
}

func OnClickConfirmNextEvent(g *Game) {
	g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["NextEvent"])
}

func (g *Game) DrawSellButton(x, y float32) Button {
	result := g.Button("Sell", x, y, OnClickNull)
	result.Rectangle.Width = 50
	result.Rectangle.Height = 25
	result.TextSize = 20
	return result
}

func DrawMarketWindow(g *Game, win *Window) {
	scene := g.Scenes["HUD"]

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
		// edge := window.X + window.Width
		// sellButton := g.DrawSellButton(float32(edge-110), float32(y))
		// sellButton.Text = "Sell Some"
		// sellButton.Rectangle.Width = 100
		// g.DrawButton(sellButton)

		// if sellAllConfirm == productName {
		// 	sellAllButton := g.DrawSellButton(float32(edge-110), float32(y))
		// 	sellAllButton.Color = rl.Red
		// 	sellAllButton.Text = "Confirm"
		// 	sellAllButton.Rectangle.Width = 100
		// 	g.DrawButton(sellAllButton)
		// 	if g.WasButtonClicked(&sellAllButton) {
		// 		scene.Data["SellAllConfirm"] = ""
		// 		g.ScreenSkip = true
		// 		g.SellProduct(g.Run.Products[product])
		// 	}

		// } else {
		// 	sellAllButton := g.DrawSellButton(float32(edge-110), float32(y))
		// 	sellAllButton.Text = "Sell All"
		// 	sellAllButton.Rectangle.Width = 100
		// 	g.DrawButton(sellAllButton)
		// 	if g.WasButtonClicked(&sellAllButton) {
		// 		scene.Data["SellAllConfirm"] = productName
		// 		g.ScreenSkip = true
		// 	}
		// }

		// if g.WasButtonClicked(&sellButton) {
		// 	OpenSellWindow(g, g.Run.Products[productName])
		// }
		i += 1

	}

	closeButton := g.CloseButton(200+900-30, 60, OnClickOpenMarketWindow)
	g.DrawButton(closeButton)
	if g.WasButtonClicked(&closeButton) {
		g.ActivateWindow(scene.Windows, scene.Windows["Prices"])
	}
}

func DrawSellWindow(g *Game, win *Window) {
	scene := g.Scenes["HUD"]

	window := rl.NewRectangle(220, 50, 500, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	rl.DrawText("Sell", 225, 60, 30, rl.Black)

	closeButton := g.CloseButton(200+500-30, 60, OnClickOpenMarketWindow)
	g.DrawButton(closeButton)
	if g.WasButtonClicked(&closeButton) {
		g.ActivateWindow(scene.Windows, scene.Windows["Sell"])
		g.ActivateWindow(scene.Windows, scene.Windows["Prices"])
	}
}

func DrawGameOverWindow(g *Game, win *Window) {
	window := rl.NewRectangle(220, 50, 500, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	rl.DrawText("Game Over!", 225, 60, 30, rl.Black)
}
