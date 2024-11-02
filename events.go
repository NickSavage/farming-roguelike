package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func InitEvents() ([]Event, error) {
	log.Printf("init")

	events := []Event{}
	var data []EventJSON

	triggerFunctions := make(map[string]func(*Game))
	triggerFunctions["Nothing"] = BlankEventOnTrigger
	triggerFunctions["Cell Tower"] = CellTowerOnTrigger
	triggerFunctions["Land Clearage"] = LandClearageOnTrigger
	triggerFunctions["Hire Help"] = HireHelpOnTrigger

	canUseFunctions := make(map[string]func(*Game) bool)
	canUseFunctions["Nothing"] = BlankEventCanUse
	canUseFunctions["Cell Tower"] = BlankEventCanUse
	canUseFunctions["Land Clearage"] = BlankEventCanUse
	canUseFunctions["Hire Help"] = HireHelpCanUse

	file, err := os.Open("./assets/events.json")
	if err != nil {
		fmt.Println(err)
		return events, err
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)

	err = jsonDecoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		return events, err
	}

	var event Event

	// Iterate over each item in the initialData slice
	for _, item := range data {
		event = Event{
			Name:        item.Name,
			Description: item.Description,
			OnTrigger:   triggerFunctions[item.Name],
			CanUse:      canUseFunctions[item.Name],
		}
		events = append(events, event)

	}
	return events, nil
}

func (g *Game) PickEventChoices(choices int) []Event {
	var results []Event

	var events = make([]Event, len(g.Run.PossibleEvents))
	copy(events, g.Run.PossibleEvents)
	rand.Shuffle(len(events), func(i, j int) {
		events[i], events[j] = events[j], events[i]
	})

	for _, event := range events {
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
		if !event.CanUse(g) {
			continue
		}
		log.Printf("event before %v", event)
		event.Effects = g.RoundEndPriceChanges()
		log.Printf("event after %v", event.Effects)
		results = append(results, event)
		if len(results) >= choices {
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

func BlankEventOnTrigger(g *Game) {

}
func BlankEventCanUse(g *Game) bool {
	return true
}

// specific events

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
func CellTowerOnTrigger(g *Game) {
	for _, space := range g.Run.TechnologySpaces {
		if space.TechnologyType != CellTowerSpace {
			continue
		}
		g.PlaceTech(g.Technology["CellTower"], space)

	}
	g.Run.EventTracker.CellTowerTriggered = true

}
func CellTowerCanUse(g *Game) bool {
	if !g.Run.EventTracker.CellTowerTriggered {
		return true
	}
	return false
}

func HireHelpOnTrigger(g *Game) {
	g.Run.ActionsMaximum += 1
	g.Run.ActionsRemaining += 1
	g.Run.EventTracker.HireHelpTriggered = true
}

func HireHelpCanUse(g *Game) bool {
	if !g.Run.EventTracker.HireHelpTriggered {
		return true
	}
	return false
}
