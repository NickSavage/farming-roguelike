package main

import (
	"math/rand"
)

func (g *Game) NewRandomEvent() Event {
	effects := []Effect{}
	// effects = append(effects, g.RandomPriceChange())
	result := Event{
		Name:    "asdf",
		Effects: effects,
	}
	return result
}

func (g *Game) RoundEndPriceChanges() Event {
	effects := []Effect{}
	productNames := g.GetProductNames()
	for _, product := range productNames {
		effect := g.RandomPriceChange(product)
		effects = append(effects, effect)
	}
	result := Event{
		Name:    "Price Changes",
		Effects: effects,
	}
	return result
}

func (g *Game) RandomPriceChange(product ProductType) Effect {
	baseRandom := rand.Float64()

	scaledRandom := baseRandom*0.2 - 0.1
	return Effect{
		IsPriceChange:   true,
		ProductImpacted: product,
		PriceChange:     float32(scaledRandom),
	}
}

// func (g *Game) CellTowerEvent() Effect {

// }
