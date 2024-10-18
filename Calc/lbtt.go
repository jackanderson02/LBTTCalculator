package calc 

import (
	"errors"
)

type LBTT struct {
	bands   []CalculatableBand
	isBuilt bool
}
type HousingTaxCalculator interface {
	Calculate(house_price float64) (float64, error)
}

type lbttWithAdditionalDwelling struct {
	lbttCalculator     HousingTaxCalculator
	additionalDwelling float64
}

type lbttWithFirstTimeBuyersRelief struct {
	lbttCalculator HousingTaxCalculator
	ftbNilRateBand float64
}

func (lbtt LBTT) Calculate(house_price float64) (float64, error) {
	if !lbtt.isBuilt {
		return 0.0, errors.New("Cannot calculate house prices until the lbtt object is built. Please call .build() to validate the given bands.")
	}
	// Goes through each of the bands and determines the tax
	var total float64 = 0
	// min_band := 0
	for _, band := range lbtt.bands {
		total += band.CalculateTaxInBand(house_price)
	}

	return total, nil
}

func (lbttWithAD lbttWithAdditionalDwelling) Calculate(house_price float64) (float64, error) {
	dwelling := 0.06 * lbttWithAD.additionalDwelling

	tax, err := lbttWithAD.lbttCalculator.Calculate(house_price)
	if err != nil {
		return 0.0, err
	}
	return dwelling + tax, nil
}

func (lbttWithFTB lbttWithFirstTimeBuyersRelief) Calculate(house_price float64) (float64, error) {
	if house_price > lbttWithFTB.ftbNilRateBand {
		return lbttWithFTB.lbttCalculator.Calculate(house_price)
	} else {
		return 0.0, nil
	}
}
