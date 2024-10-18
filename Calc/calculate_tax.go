package calc 

import (
	"errors"
	"fmt"
)

type LBTT struct {
	bands   []CalculatableBand
	isBuilt bool
}

type Band struct {
	start_range_inclusive, end_range_inclusive, consideration, rate float64
}

type CalculatableBand interface {
	CalculateTaxInBand(house_price float64) float64
	CheckValidBand(previousBand CalculatableBand) error
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

func (currentBand Band) CheckValidBand(previousBand CalculatableBand) error {
	concretePreviousBand, ok := previousBand.(Band)
	if ok {
		// do not need to validate in this case
		previous_max := concretePreviousBand.end_range_inclusive

		if currentBand.start_range_inclusive+1 != previous_max {
			return fmt.Errorf("Found non contiguous bands. Band with starting range %f does not coincide with previous band with ending range %f", currentBand.end_range_inclusive, previous_max)
		}

		if currentBand.consideration != currentBand.end_range_inclusive-previous_max || currentBand.end_range_inclusive == 0 {
			return fmt.Errorf("Found invalid conideration given. Considerations should not exceed the width of the band itself. Given consideration %f should not exceed %f", currentBand.consideration, currentBand.end_range_inclusive-previous_max)
		}

	}
	return nil
}

func (band Band) CalculateTaxInBand(house_price float64) float64 {
	band_end := band.end_range_inclusive
	// How far you are into each band * the rate. The max you can be into a band is the consideration
	band_diff := house_price - (band_end - band.consideration)
	if band_diff >= 0 { // In or exceeding this band
		if band_diff > band.consideration {
			band_diff = band.consideration // Do not exceed the consideration of this band
		}
		return band_diff * band.rate
	} else {
		return 0.0
	}
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
