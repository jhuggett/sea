package population

import (
	"log/slog"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/industry"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/timeline"
)

type Population struct {
	Persistent data.Population

	// Has Industries
	// And trade
	// That dictate supply

	IndustryLookup map[items.ItemType]*industry.Industry
}

func (p *Population) Industries() ([]industry.Industry, error) {
	if p.Persistent.Industries == nil {
		var industries []data.Industry
		err := db.Conn().Where("population_id = ?", p.Persistent.ID).Find(&industries).Error

		if err != nil {
			return nil, err
		}

		p.Persistent.Industries = industries
	}

	if p.IndustryLookup == nil {
		p.IndustryLookup = map[items.ItemType]*industry.Industry{}
	}

	inds := make([]industry.Industry, len(p.Persistent.Industries))
	for i, ind := range p.Persistent.Industries {
		inds[i] = industry.Industry{Persistent: ind}
		p.IndustryLookup[items.ItemType(ind.Product)] = &inds[i]
	}

	return inds, nil
}

func (p Population) Value(itemA items.Item) (float64, error) {
	timeInterval := timeline.Year

	if p.IndustryLookup == nil {
		_, err := p.Industries()
		if err != nil {
			slog.Error("Failed to get industries", "err", err)
			return 0.0, err
		}
	}
	speculativeAmount := 0.0

	// This is wrong.
	/*

		We need to operate off of uses.

		E.g.

		To value an item, we need to know what uses it can be put to.
		We need to know how much of those uses are need/wanted and how much
		of those uses are being fulfilled by industries / trade related to the population.

		Although novelty could also play a role.

	*/
	// industryForItem, ok := p.IndustryLookup[constructs.ItemType(itemA.Name)]
	// if ok {
	// 	speculativeAmount = float64(industryForItem.Persistent.Workers) * (float64(timeInterval) / (itemA.BaseTimeToProduce + itemA.WorkTimeToProduce))
	// }

	// slog.Info("x", "y", p.IndustryLookup, "bvc", industryForItem)

	const needWeight = 2

	// an item needs a list of uses. Uses are backed by demands.

	aAmountNeed := 0.0
	for use, perItem := range itemA.Uses {
		aAmountNeed += StandardDemands.Needs(use, p.Persistent.Size, timeInterval) * perItem
	}

	aAmountWanted := 0.0
	for use, perItem := range itemA.Uses {
		aAmountWanted += StandardDemands.Wants(use, p.Persistent.Size, timeInterval) * perItem
	}

	slog.Info("aAmountNeed", "aAmountNeed", aAmountNeed)
	slog.Info("aAmountWanted", "aAmountWanted", aAmountWanted)
	slog.Info("speculativeAmount", "speculativeAmount", speculativeAmount)

	/*

		Need to know how much is needed and how much is wanted.

		Need to keep track of surplus

		How much is project to be produced.

	*/

	// at the end normalize the ratio

	return aAmountWanted - speculativeAmount, nil
}

type Demand struct {
	Need float64
	Want float64
}

// Per person per day
type Demands map[items.Uses]Demand

func (d Demands) Needs(demand items.Uses, population uint, timeInterval uint64) float64 {
	return d[demand].Need * float64(population) * float64(timeInterval)
}

func (d Demands) Wants(demand items.Uses, population uint, timeInterval uint64) float64 {
	return d[demand].Want * float64(population) * float64(timeInterval)
}

var StandardDemands = Demands{
	items.EdibleGrain: Demand{
		Need: 0.5 / float64(timeline.Day),
		Want: 1 / float64(timeline.Day),
	},
	items.EdibleMeat: Demand{
		Need: 0.25 / float64(timeline.Day),
		Want: 2 / float64(timeline.Day),
	},
	items.Currency: Demand{
		Need: 0 / float64(timeline.Day),
		Want: 1.5 / float64(timeline.Day),
	},
}

/*

Pop of 1000

1000 grain per day
500 meat per day

already produces 500 grain per day

deficit of 500 grain per day
and 500 meat per day


so here the value of meat is the same as grain because the demand is the same?

*/

func Using(m data.Population) *Population {
	return &Population{Persistent: m}
}
