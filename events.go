package main

import (
	"math/rand"
)

func (g *Game) InitEvents() []Event {
	events := []Event{}
	events = append(events, g.LandClearageEvent())
	events = append(events, g.CellTowerEvent())

	return events
}

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

func (g *Game) ApplyEvent(event Event) {
	g.Run.EventChoices = []Event{}
	g.Run.Events = append(g.Run.Events, event)
	g.ApplyPriceChanges(event)
}

func (g *Game) ApplyPriceChanges(event Event) {
	for _, effect := range event.Effects {
		if effect.IsPriceChange {
			current := g.Run.Products[effect.ProductImpacted].Price
			g.Run.Products[effect.ProductImpacted].Price = current * (1 + effect.PriceChange)
		}
	}
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

func (g *Game) LandClearageEvent() Event {
	effects := []Effect{
		{
			IsPriceChange: false,
			EventTrigger:  LandClearageTrigger,
		},
	}
	result := Event{
		Name:    "Land Clearage",
		Effects: effects,
	}
	return result
}

func LandClearageTrigger(g *Game) {
	g.Run.EventTracker.LandClearageTriggered = true
}

func (g *Game) CellTowerEvent() Event {
	effects := []Effect{
		{
			IsPriceChange: false,
			EventTrigger:  LandClearageTrigger,
		},
	}
	result := Event{
		Name:    "Cell Tower",
		Effects: effects,
	}
	return result
}
