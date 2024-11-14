package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
)

type UnlockButton struct {
	g                *Game
	rect             rl.Rectangle
	Unlock           *Unlock
	SelectDirections engine.SelectDirections
	Selected         bool
}

func (g *Game) NewUnlockButton(rect rl.Rectangle, unlock *Unlock) UnlockButton {
	return UnlockButton{
		g:                g,
		rect:             rect,
		Unlock:           unlock,
		SelectDirections: engine.SelectDirections{},
	}
}
func (b *UnlockButton) Render() {

	mousePosition := rl.GetMousePosition()
	backgroundColor := rl.White

	if rl.CheckCollisionPointRec(mousePosition, b.rect) {
		backgroundColor = rl.LightGray
	}
	if b.IsSelected() {
		backgroundColor = rl.LightGray
	}
	if b.Unlock.OtherCost {
		if !b.Unlock.OtherCostFunction(b.g) {
			backgroundColor = rl.LightGray
		}
	}

	var desc string
	if b.Unlock.OtherCost {
		desc = b.Unlock.OtherCostDescriptionFunction(b.g)
	} else {
		desc = ""
	}

	rl.DrawRectangleRec(b.rect, backgroundColor)
	rl.DrawRectangleLinesEx(b.rect, 1, rl.Black)
	rl.DrawText(b.Unlock.Technology.Name, b.rect.ToInt32().X+5, b.rect.ToInt32().Y+5, 10, rl.Black)
	rl.DrawText(
		fmt.Sprintf("%v/%v", b.Unlock.RunSpentActions, b.Unlock.CostActions),
		b.rect.ToInt32().X+5,
		b.rect.ToInt32().Y+15,
		10,
		rl.Black,
	)
	rl.DrawText(
		desc,
		b.rect.ToInt32().X+5,
		b.rect.ToInt32().Y+55,
		10,
		rl.Black,
	)

}

func (b *UnlockButton) OnClick() {
	if b.g.Run.CanSpendAction(1) {
		b.Unlock.RunSpentActions += 1
		b.g.Run.SpendAction(1)
	}
	if b.Unlock.OtherCost {
		// to implement
		if !b.Unlock.OtherCostFunction(b.g) {
			return
		}
	}

	if b.Unlock.RunSpentActions >= b.Unlock.CostActions {
		b.Unlock.Unlocked = true
		b.Unlock.Technology.Unlocked = true

	}

}
func (b *UnlockButton) Rect() rl.Rectangle                   { return b.rect }
func (b *UnlockButton) Select()                              { b.Selected = true }
func (b *UnlockButton) Unselect()                            { b.Selected = false }
func (b *UnlockButton) IsSelected() bool                     { return b.Selected }
func (b *UnlockButton) Directions() *engine.SelectDirections { return &b.SelectDirections }

func (g *Game) InitTechUnlockWindow() {
	window := g.Scenes["Board"].Windows["UnlockWindow"]
	temp := make([]engine.UIComponent, 0)

	x := window.Rectangle.X
	y := window.Rectangle.Y
	i := 0
	for _, unlock := range g.Unlocks {
		if unlock.Unlocked {
			continue
		}

		rect := rl.NewRectangle(x+50+float32(i*160), y+45, 150, 150)
		button := g.NewUnlockButton(rect, unlock)

		temp = append(temp, &button)
		i += 1
		if i == 5 {
			y = y + 160
			i = 0
		}
	}

	components := make([]engine.UIComponent, 0)

	blank := engine.NewBlankComponent()
	blank.SelectDirections.Left = len(temp)
	blank.SelectDirections.Right = 1
	components = append(components, &blank)

	for _, button := range temp {
		components = append(components, button)
	}
	window.Components = components

}

func DrawTechUnlockWindow(gi engine.GameInterface, window *engine.Window) {

	// g := gi.(*Game)
	//	scene := g.Scenes["Board"]

	rl.DrawRectangleRec(window.Rectangle, rl.White)
	rl.DrawRectangleLinesEx(window.Rectangle, 5, rl.Black)
}

func (g *Game) InitUnlocks() {
	unlocks := make(map[string]*Unlock)
	otherCostFunctions := make(map[string]func(*Game) bool)
	otherCostDescriptionFunctions := make(map[string]func(*Game) string)

	otherCostFunctions["Chicken Egg Warmer"] = ChickenEggWarmerUnlockOtherCost
	otherCostDescriptionFunctions["Chicken Egg Warmer"] = ChickenEggWarmerUnlockOtherCostDescription

	otherCostFunctions["Cow Slaughterhouse"] = CowSlaughterhouseUnlockOtherCost
	otherCostDescriptionFunctions["Cow Slaughterhouse"] = CowSlaughterhouseUnlockOtherCostDescription
	for _, data := range g.UnlockBaseData {
		unlock := &Unlock{
			Technology:     g.Technology[data.TechnologyName],
			Unlocked:       false,
			CostActions:    data.CostActions,
			OtherCost:      data.OtherCost,
			DependencyMet:  false,
			DependencyName: data.Dependency,
		}
		if unlock.OtherCost {
			unlock.OtherCostFunction = otherCostFunctions[unlock.Technology.Name]
			unlock.OtherCostDescriptionFunction = otherCostDescriptionFunctions[unlock.Technology.Name]
			if unlock.OtherCostFunction == nil {
				log.Fatal("function not initialized for %v", unlock.Technology.Name)
			}

		}
		if unlock.DependencyName == "none" {
			unlock.DependencyMet = true
		}
		unlocks[data.TechnologyName] = unlock

	}
	g.Unlocks = unlocks
	g.UnpackUnlocks(g.UnlockSave)
}

func (g *Game) PackUnlocks() []UnlockSave {
	results := []UnlockSave{}

	for _, unlock := range g.Unlocks {
		new := UnlockSave{
			TechnologyName: unlock.Technology.Name,
			Unlocked:       unlock.Unlocked,
		}
		results = append(results, new)

	}
	return results

}

func (g *Game) UnpackUnlocks(saved []UnlockSave) {
	for _, save := range saved {
		log.Printf("save %v", save)
		g.Unlocks[save.TechnologyName].Unlocked = save.Unlocked
		g.Technology[save.TechnologyName].Unlocked = save.Unlocked
	}
	// todo chain unlocks together
}
