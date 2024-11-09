package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func (r *Run) InitEvents() ([]Event, error) {
	log.Printf("init")

	events := []Event{}
	var data []EventJSON

	triggerFunctions := make(map[string]func(*Game))
	triggerFunctions["Nothing"] = BlankEventOnTrigger
	triggerFunctions["Cell Tower"] = CellTowerOnTrigger
	triggerFunctions["Land Clearage"] = LandClearageOnTrigger
	triggerFunctions["Hire Help"] = HireHelpOnTrigger
	triggerFunctions["Flood"] = FloodOnTrigger
	r.triggerFunctions = triggerFunctions

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
			Repeatable:  item.Repeatable,
			OnTrigger:   triggerFunctions[item.Name],
			Severity:    item.Severity,
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
		if !CheckCanUseEvent(g, event) {
			continue
		}
		log.Printf("event before %v", event)
		event.Effects = g.RoundEndPriceChanges(event)
		log.Printf("event after %v", event.Effects)
		results = append(results, event)
		if len(results) >= choices {
			break
		}
	}
	return results
}

func (g *Game) RoundEndPriceChanges(event Event) []Effect {
	effects := []Effect{}
	productNames := g.GetProductNames()
	for _, product := range productNames {
		effect := g.RandomPriceChange(product, event.Severity)
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

	g.Run.Money -= event.CostMoney

	if !event.Repeatable {
		log.Printf("save event")
		g.Run.EventTracker[event.Name] = true
	}
	log.Printf("tracker %v", g.Run.EventTracker)
}

func (g *Game) ApplyPriceChanges(event Event) {
	for _, effect := range event.Effects {
		if effect.IsPriceChange {
			current := g.Run.Products[effect.ProductImpacted].Price
			g.Run.Products[effect.ProductImpacted].Price = current * (1 + effect.PriceChange)
		}
	}
}

func (g *Game) RandomPriceChange(product ProductType, severity float32) Effect {
	adjustment := severity / 10
	baseRandom := rand.Float64() + float64(adjustment)

	var scaledRandom float64
	if product == Solar {
		scaledRandom = baseRandom*0.3 - 0.2

	} else {
		scaledRandom = baseRandom*0.2 - 0.1

	}
	log.Printf("base %v adjust %v final %v", baseRandom, adjustment, scaledRandom)
	return Effect{
		IsPriceChange:   true,
		ProductImpacted: product,
		PriceChange:     float32(scaledRandom),
	}
}

func CheckCanUseEvent(g *Game, event Event) bool {

	if result, ok := g.Run.EventTracker[event.Name]; ok {
		if result {
			return false
		}
		return true
	} else {
		if !g.Run.CanSpendMoney(event.CostMoney) {
			return false

		}
		return true
	}
}

// specific events

func BlankEventOnTrigger(g *Game) {

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
func CellTowerOnTrigger(g *Game) {
	for _, space := range g.Run.TechnologySpaces {
		if space.TechnologyType != CellTowerSpace {
			continue
		}
		g.PlaceTech(g.Technology["Cell Tower"], space)

	}

}

func HireHelpOnTrigger(g *Game) {
	g.Run.ActionsMaximum += 1
	g.Run.ActionsRemaining += 1
}

func FloodOnTrigger(g *Game) {
	for _, space := range g.Run.TechnologySpaces {
		if !space.IsFilled {
			continue
		}
		if space.TechnologyType != BuildingSpace {
			continue
		}
		g.RemoveTech(space.Technology)
	}
}
