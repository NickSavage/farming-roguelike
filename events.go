package main

import (
	"log"
	"math/rand"
)

func (g *Game) InitEvents() []Event {
	events := []Event{}
	log.Printf("init")
	events = append(events, g.BlankEvent())
	events = append(events, g.LandClearageEvent())
	events = append(events, g.CellTowerEvent())
	return events
}

func (g *Game) PickEventChoices(choices int) []Event {
	var results []Event

	for _ = range choices {

		for _, event := range g.Run.PossibleEvents {
			alreadyPicked := false
			for _, existingChoice := range results {
				if existingChoice.Name == event.Name {
					alreadyPicked = true
					break
				}
			}
			if alreadyPicked {
				continue
			}
			log.Printf("event before %v", event)
			event.Effects = g.RoundEndPriceChanges()
			log.Printf("event after %v", event.Effects)
			results = append(results, event)
			break
		}
	}
	return results
}

func (g *Game) RoundEndPriceChanges() []Effect {
	effects := []Effect{}
	productNames := g.GetProductNames()
	for _, product := range productNames {
		effect := g.RandomPriceChange(product)
		effects = append(effects, effect)
	}
	log.Printf("effects %v", effects)
	return effects
}

func (g *Game) ApplyEvent(event Event) {
	g.Run.EventChoices = []Event{}
	g.Run.Events = append(g.Run.Events, event)
	g.ApplyPriceChanges(event)
	event.OnTrigger(g)
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

// blank event

func (g *Game) BlankEvent() Event {
	result := Event{
		Name:        "Nothing",
		Description: "Add a new field",
		OnTrigger:   BlankEventOnTrigger,
	}
	return result

}

func BlankEventOnTrigger(g *Game) {

}

// specific events

func (g *Game) LandClearageEvent() Event {
	result := Event{
		Name:        "Land Clearage",
		Description: "Add a new field",
		OnTrigger:   LandClearageOnTrigger,
	}
	return result
}

func LandClearageOnTrigger(g *Game) {
	for _, space := range g.Run.TechnologySpaces {
		if space.Active {
			continue
		}
		space.Active = true
		break // only pick one
	}
	// g.Run.EventTracker.LandClearageTriggered = true
}

func (g *Game) CellTowerEvent() Event {
	effects := []Effect{
		{
			IsPriceChange: false,
			EventTrigger:  CellTowerOnTrigger,
		},
	}
	result := Event{
		Name:        "Cell Tower",
		Description: "A major telecom company has approached you about building a cell tower on your property.",
		OnTrigger:   CellTowerOnTrigger,

		Effects: effects,
	}
	return result
}

func CellTowerOnTrigger(g *Game) {
	for _, space := range g.Run.TechnologySpaces {
		if space.TechnologyType != CellTowerSpace {
			continue
		}
		g.PlaceTech(g.Technology["CellTower"], space)

	}

}
