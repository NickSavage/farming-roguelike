package main

import (
	"math/rand"
)

func (g *Game) NewRandomEvent() Event {
	effects := []Effect{}
	effects = append(effects, g.RandomPriceChange())
	result := Event{
		Name:    "asdf",
		Effects: effects,
	}
	return result
}

func (g *Game) RandomPriceChange() Effect {
	var productNames []string
	for key, _ := range g.Run.Products {
		productNames = append(productNames, key)
	}
	index := rand.Intn(len(productNames))
	baseRandom := rand.Float64()

	scaledRandom := baseRandom*0.2 - 0.1
	return Effect{
		IsPriceChange:   true,
		ProductImpacted: productNames[index],
		PriceChange:     float32(scaledRandom),
	}
}
