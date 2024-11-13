package main

import (
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

	rl.DrawRectangleRec(b.rect, backgroundColor)
	rl.DrawRectangleLinesEx(b.rect, 1, rl.Black)
	rl.DrawText(b.Unlock.Technology.Name, b.rect.ToInt32().X+5, b.rect.ToInt32().Y+5, 10, rl.Black)

}

func (b *UnlockButton) OnClick()                             {}
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

		rect := rl.NewRectangle(x+50+float32(i*160), y+45, 150, 300)
		button := g.NewUnlockButton(rect, unlock)

		temp = append(temp, &button)
		i += 1
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
}

func (r *Run) PackUnlocks() []UnlockSave {
	results := []UnlockSave{}

	for _, unlock := range r.Game.Unlocks {
		new := UnlockSave{
			TechnologyName: unlock.Technology.Name,
			Unlocked:       unlock.Unlocked,
		}
		results = append(results, new)

	}
	return results

}

func (r *Run) UnpackUnlocks(saved []UnlockSave) {
	for _, save := range saved {
		log.Printf("save %v", save)
		r.Game.Unlocks[save.TechnologyName].Unlocked = save.Unlocked
		r.Game.Technology[save.TechnologyName].Unlocked = save.Unlocked
	}
	// todo chain unlocks together
}
