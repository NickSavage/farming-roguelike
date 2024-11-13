package main

func (g *Game) PackProductStats() []ProductStat {
	results := []ProductStat{}
	for _, value := range g.ProductStats {
		results = append(results, *value)
	}
	return results
}

func (g *Game) UnpackProductStats(input []ProductStat) map[ProductType]*ProductStat {
	results := make(map[ProductType]*ProductStat)

	for _, stat := range input {
		stat.CurrentRunEarned = 0
		stat.CurrentRunProduction = 0
		results[ProductType(stat.ProductType)] = &stat
	}
	return results

}

func (g *Game) RecordProductStat(product *Product, quantity float32, earned float32) {

	if current, exists := g.ProductStats[product.Type]; !exists {
		g.ProductStats[product.Type] = &ProductStat{
			ProductType:          product.Type,
			MaxProduction:        quantity,
			TotalProduction:      quantity,
			CurrentRunProduction: quantity,
			MaxEarned:            earned,
			TotalEarned:          earned,
			CurrentRunEarned:     earned,
		}
	} else {
		current.CurrentRunProduction += quantity
		current.TotalProduction += quantity
		if current.CurrentRunProduction > current.MaxProduction {
			current.MaxProduction = current.CurrentRunProduction
		}
		current.CurrentRunEarned += earned
		current.TotalEarned += earned
		if current.CurrentRunEarned > current.MaxEarned {
			current.MaxEarned = current.CurrentRunEarned
		}
	}
}
