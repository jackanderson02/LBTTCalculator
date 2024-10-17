package main

import (
	"errors"
)

type AggeregatedCalculator struct {
	IsBuilt                       bool
	containsFTBCalc               bool
	containsDwellingCalc          bool
	containsLinkedTransactionCalc bool
	calulators                    []HousingTaxCalculator
	finalCalculator               HousingTaxCalculator
}

func NewAggregatedCalculator() *AggeregatedCalculator {
	return &AggeregatedCalculator{
		containsFTBCalc:               false,
		containsDwellingCalc:          false,
		containsLinkedTransactionCalc: false,
	}
}

func (aggCalc *AggeregatedCalculator) addCalc(calc HousingTaxCalculator) *AggeregatedCalculator {
	aggCalc.calulators = append(aggCalc.calulators, calc)
	return aggCalc
}

func (aggCalc *AggeregatedCalculator) WithAdditionDwellingCalculator(dwelling_price float64) *AggeregatedCalculator {
	aggCalc.addCalc(lbttWithAdditionalDwelling{
		additionalDwelling: dwelling_price,
	})
	aggCalc.containsDwellingCalc = true

	return aggCalc
}

func (aggCalc *AggeregatedCalculator) WithFTBCalculator(ftbNilRate float64) *AggeregatedCalculator {
	aggCalc.addCalc(lbttWithFirstTimeBuyersRelief{
		ftbNilRateBand: ftbNilRate,
	})
	aggCalc.containsFTBCalc = true

	return aggCalc
}

func (aggCalc *AggeregatedCalculator) Build(err *error) *AggeregatedCalculator {
	// Order of the calulators does not matter I dont think. Just cannot have additional dwellings and a FTB calculator in same one
	// This is extensible for not allowing ftb relief for other conditions such as linked transactions

	if aggCalc.containsFTBCalc {
		// Check for additional dwellings calc
		if aggCalc.containsDwellingCalc || aggCalc.containsLinkedTransactionCalc {
			*err = errors.New("This calculator aggregation is not allowed. Transactions to which ADS and linked transaction fees apply are not eligible to FTB relief.")
		}

	} else {
		// Compose the calculator as a nested structure
		var nestedCalculator HousingTaxCalculator
		calculators := aggCalc.calulators

		for i := len(calculators) - 1; i >= 0; i-- {
			if i == len(calculators)-1 {
				// Use the last element as the base calculator
				nestedCalculator = calculators[i]
			} else {
				// Wrap the current nestedCalculator in a new one based on the type of the current calculator
				switch calc := calculators[i].(type) {
				case lbttWithFirstTimeBuyersRelief:
					nestedCalculator = lbttWithFirstTimeBuyersRelief{
						lbttCalculator: nestedCalculator,
						ftbNilRateBand: calc.ftbNilRateBand,
					}
				case lbttWithAdditionalDwelling:
					nestedCalculator = lbttWithAdditionalDwelling{
						lbttCalculator:     nestedCalculator,
						additionalDwelling: calc.additionalDwelling,
					}
				}
			}
		}
		aggCalc.IsBuilt = true
		aggCalc.finalCalculator = nestedCalculator
	}

	return aggCalc

}

func (aggCalc *AggeregatedCalculator) Calculate(house_price float64) (float64, error){
	return aggCalc.finalCalculator.Calculate(house_price)
}
