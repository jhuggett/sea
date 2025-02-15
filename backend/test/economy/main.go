package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/economy"
	"github.com/jhuggett/sea/data/industry"
	"github.com/jhuggett/sea/data/population"
	"github.com/jhuggett/sea/data/port"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/log"
	"github.com/jhuggett/sea/timeline"

	"github.com/guptarohit/asciigraph"
)

func do() (uint, string) {
	log.UnderTest()
	db.SetupInMemDB()

	// port
	portModel := port.Port{
		Persistent: data.Port{},
	}

	_, err := portModel.Create()
	if err != nil {
		panic(err)
	}

	// economy
	economyModel := economy.Economy{
		Persistent: data.Economy{
			PortID: portModel.Persistent.ID,
		},
	}
	_, err = economyModel.Create()
	if err != nil {
		panic(err)
	}

	// town population
	townPopulationModel := population.Population{
		Persistent: data.Population{
			EconomyID: economyModel.Persistent.ID,
			Size:      1000,
		},
	}
	_, err = townPopulationModel.Create()
	if err != nil {
		panic(err)
	}

	// grain industry
	grainIndustryModel := industry.Industry{
		Persistent: data.Industry{
			EconomyID:      economyModel.Persistent.ID,
			ShareOfWorkers: 0,
			Product:        string(items.Grain),
		},
	}
	_, err = grainIndustryModel.Create()
	if err != nil {
		panic(err)
	}

	// fish industry
	fishIndustryModel := industry.Industry{
		Persistent: data.Industry{
			EconomyID:      economyModel.Persistent.ID,
			ShareOfWorkers: 0,
			Product:        string(items.Fish),
		},
	}
	_, err = fishIndustryModel.Create()
	if err != nil {
		panic(err)
	}

	// wood industry
	woodIndustryModel := industry.Industry{
		Persistent: data.Industry{
			EconomyID:      economyModel.Persistent.ID,
			ShareOfWorkers: 0,
			Product:        string(items.Wood),
		},
	}
	_, err = woodIndustryModel.Create()
	if err != nil {
		panic(err)
	}

	// Test

	err = economyModel.Hydrate(economy.WithIndustries, economy.WithPopulations, economy.WithMarkets)
	if err != nil {
		panic(err)
	}

	err = economyModel.InitializeMarkets()
	if err != nil {
		panic(err)
	}

	err = economyModel.Hydrate(economy.WithIndustries, economy.WithPopulations, economy.WithMarkets)
	if err != nil {
		panic(err)
	}

	populationDataOverTime := map[string][]struct {
		size uint
	}{}

	marketDataOverTime := map[string]struct {
		surplus                 []float64
		workerShare             []float64
		HistoricalSupplyAverage []float64
		HistoricalDemandAverage []float64
		HistoricalSurplus       []float64
	}{}

	normalizedRelativeValuesOverTime := map[string][]float64{}

	for i := 0; i < 1000; i++ {
		err = economyModel.Tick(7 * timeline.Day)
		if err != nil {
			panic(err)
		}

		err = economyModel.Hydrate(economy.WithIndustries, economy.WithPopulations, economy.WithMarkets)
		if err != nil {
			panic(err)
		}

		totalSurplus := 0.0

		for _, market := range economyModel.Persistent.Markets {
			marketDataOverTime[market.ProductName] = struct {
				surplus                 []float64
				workerShare             []float64
				HistoricalSupplyAverage []float64
				HistoricalDemandAverage []float64
				HistoricalSurplus       []float64
			}{
				surplus:                 append(marketDataOverTime[market.ProductName].surplus, market.Surplus),
				workerShare:             marketDataOverTime[market.ProductName].workerShare,
				HistoricalSupplyAverage: append(marketDataOverTime[market.ProductName].HistoricalSupplyAverage, market.HistoricalSupplyAverage),
				HistoricalDemandAverage: append(marketDataOverTime[market.ProductName].HistoricalDemandAverage, market.HistoricalDemandAverage),
				HistoricalSurplus:       append(marketDataOverTime[market.ProductName].HistoricalSurplus, market.HistoricalSurplus),
			}

			totalSurplus += market.Surplus
		}

		for _, industry := range economyModel.Persistent.Industries {
			marketDataOverTime[industry.Product] = struct {
				surplus                 []float64
				workerShare             []float64
				HistoricalSupplyAverage []float64
				HistoricalDemandAverage []float64
				HistoricalSurplus       []float64
			}{
				surplus:                 marketDataOverTime[industry.Product].surplus,
				HistoricalSupplyAverage: marketDataOverTime[industry.Product].HistoricalSupplyAverage,
				HistoricalDemandAverage: marketDataOverTime[industry.Product].HistoricalDemandAverage,
				HistoricalSurplus:       marketDataOverTime[industry.Product].HistoricalSurplus,
				workerShare:             append(marketDataOverTime[industry.Product].workerShare, float64(industry.ShareOfWorkers)),
			}
		}

		for _, market := range economyModel.Persistent.Markets {
			normalizedRelativeValuesOverTime[market.ProductName] = append(normalizedRelativeValuesOverTime[market.ProductName], market.HistoricalSupplyAverage/market.HistoricalDemandAverage)
		}

		for _, population := range economyModel.Persistent.Populations {
			populationDataOverTime[fmt.Sprint(population.ID)] = append(populationDataOverTime[fmt.Sprint(population.ID)], struct {
				size uint
			}{
				size: population.Size,
			})
		}
	}

	maxGraphWidth := 40
	maxGraphHeight := 10

	for productName, marketData := range marketDataOverTime {
		fmt.Println()
		fmt.Println(strings.Repeat(" ", maxGraphWidth/2-len(productName)/2) + strings.ToUpper(productName))
		// graph := asciigraph.Plot(marketData.surplus)
		graph := asciigraph.PlotMany(
			[][]float64{marketData.surplus, marketData.HistoricalSurplus},
			asciigraph.SeriesColors(asciigraph.Blue, asciigraph.Gold),
			asciigraph.Width(maxGraphWidth),
			asciigraph.Height(maxGraphHeight),
			asciigraph.SeriesLegends("Surplus", "Rolling Surplus Average"),
			asciigraph.Caption("Surplus"),
		)
		fmt.Println(graph)
		graph = asciigraph.Plot(marketData.workerShare,
			asciigraph.Width(maxGraphWidth),
			asciigraph.Height(maxGraphHeight),
			asciigraph.Caption("Worker Share"),
		)
		fmt.Println(graph)
		graph = asciigraph.PlotMany(
			[][]float64{marketData.HistoricalDemandAverage, marketData.HistoricalSupplyAverage},
			asciigraph.SeriesColors(asciigraph.Red, asciigraph.Green),
			asciigraph.Width(maxGraphWidth),
			asciigraph.Height(maxGraphHeight),
			asciigraph.SeriesLegends("Demand", "Supply"),
			asciigraph.Caption("Rolling Demand Average vs Rolling Supply Average"),
		)
		fmt.Println(graph)
		fmt.Println(strings.Repeat("-", maxGraphWidth))
	}

	productNames := []string{}
	relativeValues := [][]float64{}

	for productName, values := range normalizedRelativeValuesOverTime {
		productNames = append(productNames, productName)
		relativeValues = append(relativeValues, values)
	}

	graph := asciigraph.PlotMany(relativeValues,
		asciigraph.SeriesLegends(productNames...),
		asciigraph.SeriesColors(asciigraph.Red, asciigraph.Green, asciigraph.Blue),
		asciigraph.Width(maxGraphWidth),
		asciigraph.Height(maxGraphHeight),
		asciigraph.Caption("Relative value of each product"),
	)

	fmt.Println(graph)

	populationNames := []string{}
	populationsOverTime := [][]float64{}

	for populationName, values := range populationDataOverTime {
		populationNames = append(populationNames, populationName)
		populationSizes := []float64{}
		for _, value := range values {
			populationSizes = append(populationSizes, float64(value.size)/100)
		}
		populationsOverTime = append(populationsOverTime, populationSizes)
	}

	graph = asciigraph.PlotMany(populationsOverTime,
		asciigraph.SeriesLegends(populationNames...),
		asciigraph.SeriesColors(asciigraph.Default, asciigraph.Green, asciigraph.Blue),
		asciigraph.Width(maxGraphWidth),
		asciigraph.Height(maxGraphHeight),
		asciigraph.Caption("Population size"),
	)

	fmt.Println(graph)

	for _, population := range economyModel.Persistent.Populations {
		slog.Info("Population", "population.size", population.Size)
	}

	return economyModel.Persistent.Populations[0].Size, graph
}

func main() {
	run()
}
