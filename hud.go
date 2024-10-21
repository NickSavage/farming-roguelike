package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
	"math"
)

func OnClickNull(g *Game) {}

func OnClickShopWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["ShopWindow"])
}

func OnClickTechWindowButton(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["TechWindow"])
}

func OnClickOpenEndRoundPage1Window(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["EndRound1"])
}

func OnClickOpenMarketWindow(g *Game) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["Prices"])
}

func OnClickEndRoundPageTwoButton(g *Game) {

}

func OnClickEndRoundConfirmButton(g *Game) {

}

func OpenSellWindow(g *Game, product *Product) {
	scene := g.Scenes["HUD"]
	g.ActivateWindow(scene.Windows, scene.Windows["Sell"])

}

func (g *Game) InitHUD() {
	scene := g.Scenes["HUD"]

	scene.Windows = make(map[string]*Window)
	scene.Windows["ShopWindow"] = &Window{
		Name:       "Shop Window",
		Display:    false,
		DrawWindow: DrawShopWindow,
	}
	scene.Windows["TechWindow"] = &Window{
		Name:       "Tech Window",
		Display:    false,
		DrawWindow: DrawTechnologyWindow,
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

}

func UpdateHUD(g *Game) {
	scene := g.Scenes["HUD"]
	for _, button := range scene.Buttons {
		if g.WasButtonClicked(&button) {
			button.OnClick(g)
		}
	}

}

func DrawHUD(g *Game) {
	scene := g.Scenes["HUD"]
	height := int32(150)
	sidebarWidth := int32(200)
	rl.DrawRectangle(0, g.screenHeight-height, g.screenWidth, height, rl.Black)
	rl.DrawRectangle(0, 0, sidebarWidth, g.screenHeight-height, rl.Black)

	DrawSidebar(g)
	g.DrawButtons(scene.Buttons)

	if g.Data["Message"].(string) != "" {
		rl.DrawText(g.Data["Message"].(string), 205, g.screenHeight-height+15, 20, rl.White)
		if g.Data["MessageCounter"].(int32) == g.Seconds {
			g.Data["Message"] = ""
			g.Data["MessageCounter"] = 0
		}

	}
	for _, window := range scene.Windows {
		if window.Display {
			window.DrawWindow(g, window)
		}
	}
}

func DrawSidebar(g *Game) {

	rl.DrawText(
		fmt.Sprintf("Actions: %v/%v", g.Run.RoundActionsRemaining, g.Run.RoundActions),
		30, 30, 20, rl.White,
	)
	rl.DrawText(fmt.Sprintf("Money: $%v", g.Run.Money), 30, 50, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Round: %v", g.Run.CurrentRound), 30, 70, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Season: %v", g.Run.CurrentSeason.String()), 30, 90, 20, rl.White)

	buttons := []*Button{}
	techButton := g.Button("Technology", 10, 150, OnClickTechWindowButton)
	buttons = append(buttons, &techButton)
	shopButton := g.Button("Shop", 10, 200, OnClickShopWindowButton)
	buttons = append(buttons, &shopButton)
	priceButton := g.Button("Market", 10, 250, OnClickOpenMarketWindow)
	buttons = append(buttons, &priceButton)
	viewEndRoundButton := g.Button("End Round", 10, 300, OnClickOpenEndRoundPage1Window)
	buttons = append(buttons, &viewEndRoundButton)

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

	rl.DrawText("Income", int32(windowRect.X+5), int32(windowRect.Y+5), 30, rl.Black)
	var totalEarned float32 = 0

	var x, y int32
	for i, tech := range g.Run.Technology {
		x = int32(windowRect.X + 10)
		y = int32(windowRect.Y + 50 + float32(i*30))
		index := tech.RoundHandlerIndex
		value := tech.RoundHandler[index].RoundEndValue(g, tech)
		totalEarned += value
		text := tech.RoundHandler[index].RoundEndText(g, tech)
		rl.DrawText(text, x, y, 20, rl.Black)
	}

	text := fmt.Sprintf("Total: $%v", totalEarned)
	rl.DrawText(text, x, y+30, 20, rl.Black)

	button := g.Button("Next Page", 500, 500, OnClickEndRoundPageTwoButton)

	g.DrawButton(button)
	if g.WasButtonClicked(&button) {
		g.ActivateWindow(g.Scenes["HUD"].Windows, g.Scenes["HUD"].Windows["EndRound2"])
	}
}

func DrawEndRoundWindowPage2(g *Game, win *Window) {

	windowRect := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(windowRect, rl.White)
	rl.DrawRectangleLinesEx(windowRect, 5, rl.Black)

	rl.DrawText("Investments", int32(windowRect.X+5), int32(windowRect.Y+5), 30, rl.Black)

	var actions float32 = float32(g.Run.RoundActions)

	var x, y int32
	for i, tech := range g.Run.Technology {
		x = int32(windowRect.X + 10)
		y = int32(windowRect.Y + 50 + float32(i*30))
		nextSeason := tech.RoundHandler[tech.RoundHandlerIndex]

		actions -= nextSeason.CostActions
		text := fmt.Sprintf(
			"%s: -%v actions -$%v money",
			tech.Name,
			nextSeason.CostActions,
			nextSeason.CostMoney,
		)
		rl.DrawText(text, x, y, 20, rl.Red)

	}
	text := fmt.Sprintf("Actions next season: %v", actions)
	rl.DrawText(text, x, y+30, 20, rl.Red)
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

func DrawNextEventWindow(g *Game, win *Window) {

	window := rl.NewRectangle(220, 50, 900, 500)
	rl.DrawRectangleRec(window, rl.White)
	rl.DrawRectangleLinesEx(window, 5, rl.Black)

	button := g.Button("Confirm", 500, 500, OnClickConfirmNextEvent)

	event := g.Run.Events[g.Run.CurrentRound]
	g.DrawButton(button)
	rl.DrawText(event.Name, 225, 60, 30, rl.Black)

	for _, effect := range event.Effects {
		if effect.IsPriceChange {
			newPrice := g.Run.Products[effect.ProductImpacted].Price * float32(1+effect.PriceChange)
			newPrice = float32(math.Round(float64(newPrice*100))) / 100

			displayChange := math.Round(float64(effect.PriceChange*100*100)) / 100
			text := fmt.Sprintf("Price of %v is now %v (%v%%)", effect.ProductImpacted, newPrice, displayChange)
			rl.DrawText(text, 225, 95, 20, rl.Black)
		}
	}

	if g.WasButtonClicked(&button) {
		button.OnClick(g)
	}

}

func OnClickConfirmNextEvent(g *Game) {
	g.ProcessNextEvent()
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
	products := g.GetProductNames()
	rl.DrawText("Products", int32(window.X+10), int32(window.Y+50), 25, rl.Black)
	rl.DrawText("Inventory", int32(window.X+10)+columnOffset, int32(window.Y+50), 25, rl.Black)
	rl.DrawText("Spot Price", int32(window.X+10)+columnOffset*2, int32(window.Y+50), 25, rl.Black)
	for _, productName := range products {
		x = int32(window.X + 10)
		y = int32(window.Y + 80 + float32(i*30))
		rl.DrawText(productName, x, y, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%v", g.Run.Products[productName].Quantity), x+columnOffset, y, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%v", g.Run.Products[productName].Price), x+columnOffset*2, y, 20, rl.Black)

		edge := window.X + window.Width
		sellButton := g.DrawSellButton(float32(edge-110), float32(y))
		sellButton.Text = "Sell Some"
		sellButton.Rectangle.Width = 100
		g.DrawButton(sellButton)

		sellAllButton := g.DrawSellButton(float32(edge-220), float32(y))
		sellAllButton.Text = "Sell All"
		sellAllButton.Rectangle.Width = 100
		g.DrawButton(sellAllButton)

		if g.WasButtonClicked(&sellButton) {
			OpenSellWindow(g, g.Run.Products[productName])
		}
		if g.WasButtonClicked(&sellAllButton) {
			OpenSellWindow(g, g.Run.Products[productName])
		}

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
