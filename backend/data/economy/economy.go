package economy

import (
	"fmt"
	"math"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/population"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/timeline"
	"gorm.io/gorm"
)

type Economy struct {
	Persistent data.Economy
}

func (e *Economy) Create() (uint, error) {
	err := db.Conn().Create(&e.Persistent).Error

	if err != nil {
		return 0, err
	}

	return e.Persistent.ID, nil
}

func (e *Economy) IndustriesByProduct() map[string]data.Industry {
	industries := map[string]data.Industry{}

	for _, industry := range e.Persistent.Industries {
		industries[industry.Product] = industry
	}

	return industries
}

func (e *Economy) InitializeMarkets() error {
	markets := []data.Market{}

	for _, industry := range e.Persistent.Industries {
		market := data.Market{
			ProductName: industry.Product,

			Surplus:                 0,
			HistoricalSupplyAverage: 0,
			HistoricalDemandAverage: 0,
			HistoricalSurplus:       0,

			EconomyID: e.Persistent.ID,

			SampleCount: 0,
		}

		markets = append(markets, market)
	}

	return db.Conn().Create(&markets).Error
}

type Workforce struct {
	Total uint

	Industries map[string]uint
}

func (e *Economy) Workforce() Workforce {
	workforce := Workforce{
		Industries: map[string]uint{},
	}

	for _, population := range e.Persistent.Populations {
		workforce.Total += population.Size
	}

	totalShareOfWorkers := 0.0
	for _, industry := range e.Persistent.Industries {
		totalShareOfWorkers += industry.ShareOfWorkers
	}

	for _, industry := range e.Persistent.Industries {
		workforce.Industries[industry.Product] = uint(math.Floor(float64(workforce.Total) * (industry.ShareOfWorkers / totalShareOfWorkers)))
	}

	return workforce
}

func (e *Economy) DemandsPerDay() (map[items.Uses]population.Demand, error) {
	demands := map[items.Uses]population.Demand{}

	for _, populationData := range e.Persistent.Populations {
		demandsForPopulation, err := population.Using(populationData).DemandsPerDay()
		if err != nil {
			return nil, err
		}

		for product, demand := range demandsForPopulation {
			demands[product] = population.Demand{
				Need: demands[product].Need + demand.Need,
				Want: demands[product].Want + demand.Want,
			}
		}
	}

	return demands, nil
}

func (e *Economy) Tick(elapsedTime timeline.Tick) error {

	industriesByProduct := e.IndustriesByProduct()

	workforce := e.Workforce()

	amountOfNeedsMet := map[items.Uses]float64{}

	demandPerDay, err := e.DemandsPerDay()
	if err != nil {
		return fmt.Errorf("failed to get demands per day: %w", err)
	}

	for _, marketData := range e.Persistent.Markets {
		var produced float64
		var needed float64

		industryForProduct, ok := industriesByProduct[marketData.ProductName]
		if ok {

			productConstruct := items.LookupItem(marketData.ProductName)

			produced =
				float64(workforce.Industries[marketData.ProductName]) *
					(float64(elapsedTime) / productConstruct.BaseTimeToProduce)

			needPerUse := map[items.Uses]float64{}

			for use, percent := range productConstruct.Uses {
				needPerUse[use] = demandPerDay[use].Need * float64(elapsedTime/timeline.Day) * percent * float64(workforce.Total)
			}

			// needed needs to take into consideration the amount of the needs already met

			for _, need := range needPerUse {
				needed += need
			}

			marketData.SampleCount++

			delta := produced - needed

			percentOfNeedsMet := 1.0

			if delta+marketData.Surplus < 0 {
				percentOfNeedsMet = (delta + marketData.Surplus + needed) / needed
			}

			for use, _ := range productConstruct.Uses {
				amountOfNeedsMet[use] += needPerUse[use] * percentOfNeedsMet
			}

			marketData.Surplus = math.Max(0, marketData.Surplus+delta)

			// rolling averages
			marketData.HistoricalSupplyAverage = (marketData.HistoricalSupplyAverage*(float64(marketData.SampleCount)-1) + produced) / float64(marketData.SampleCount)
			marketData.HistoricalDemandAverage = (marketData.HistoricalDemandAverage*(float64(marketData.SampleCount)-1) + needed) / float64(marketData.SampleCount)
			marketData.HistoricalSurplus = (marketData.HistoricalSurplus*(float64(marketData.SampleCount)-1) + marketData.Surplus) / float64(marketData.SampleCount)

			if marketData.Surplus <= marketData.HistoricalSupplyAverage {
				industryForProduct.ShareOfWorkers += 1
			} else if marketData.Surplus > (marketData.HistoricalDemandAverage * 2.5) {
				industryForProduct.ShareOfWorkers = math.Max(0, industryForProduct.ShareOfWorkers-1)
			}

			err := db.Conn().Save(&industryForProduct).Error
			if err != nil {
				return err
			}

			err = db.Conn().Save(&marketData).Error
			if err != nil {
				return err
			}
		}
	}

	totalNeeds := 0.0
	totalNeedsMet := 0.0
	for use, met := range amountOfNeedsMet {
		need := demandPerDay[use].Need * float64(elapsedTime/timeline.Day) * float64(e.Workforce().Total)
		totalNeeds += need
		totalNeedsMet += math.Min(met/need, need)
	}

	needsMetPercent := totalNeedsMet / totalNeeds

	if needsMetPercent < 1 {
		for _, populationData := range e.Persistent.Populations {
			populationData.Size = uint(math.Max(0, float64(populationData.Size)-float64(populationData.Size)*.1*needsMetPercent))
			err := db.Conn().Save(&populationData).Error
			if err != nil {
				return err
			}
		}
	} else { // to simulate population growth
		// for _, populationData := range e.Persistent.Populations {
		// 	populationData.Size = uint(math.Min(1000000, float64(populationData.Size)+float64(populationData.Size)*.1*(needsMetPercent-1)))
		// 	err := db.Conn().Save(&populationData).Error
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	return nil
}

func (e *Economy) Market(product items.ItemType) (data.Market, error) {
	for _, market := range e.Persistent.Markets {
		if market.ProductName == string(product) {
			return market, nil
		}
	}
	return data.Market{}, gorm.ErrRecordNotFound
}

func (e *Economy) Hydrate(scopes ...func(*gorm.DB) *gorm.DB) error {
	return db.Conn().Scopes(scopes...).First(&e.Persistent, e.Persistent.ID).Error
}

func WithMarkets(db *gorm.DB) *gorm.DB {
	return db.Preload("Markets")
}

func WithIndustries(db *gorm.DB) *gorm.DB {
	return db.Preload("Industries")
}

func WithPopulations(db *gorm.DB) *gorm.DB {
	return db.Preload("Populations")
}
