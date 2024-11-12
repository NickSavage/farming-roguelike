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

func (g *Game) InitTechUnlockWindow() {

}

func DrawTechUnlockWindow(gi engine.GameInterface) {

}

func UpdateTechUnlockWindow(gi engine.GameInterface) {

}

func (g *Game) InitUnlocks() {
	unlocks := make(map[string]*Unlock)
	otherCostFunctions := make(map[string]func(*Game) bool)

	for _, data := range g.UnlockBaseData {
		unlock := &Unlock{
			Technology:  g.Technology[data.TechnologyName],
			Unlocked:    false,
			CostActions: data.CostActions,
			OtherCost:   data.OtherCost,
		}
		if unlock.OtherCost {
			unlock.OtherCostFunction = otherCostFunctions[unlock.Technology.Name]
			if unlock.OtherCostFunction == nil {
				log.Fatal("function not initialized for %v", unlock.Technology.Name)
			}

		}
		unlocks[data.TechnologyName] = unlock

	}
	g.Unlocks = unlocks
}
