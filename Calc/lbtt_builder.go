package calc 

import (
	"errors"
)

func NewLBTT() *LBTT {
	return &LBTT{}
}

func (lbtt *LBTT) addBand(band Band) *LBTT {
	if !lbtt.isBuilt {
		lbtt.bands = append(lbtt.bands, band)
	}
	return lbtt
}

func (lbtt *LBTT) WithBand(start_range_inclusive, end_range_inclusive, consideration, rate float64) *LBTT {
	if !lbtt.isBuilt {
		lbtt.addBand(Band{start_range_inclusive: start_range_inclusive, end_range_inclusive: end_range_inclusive, consideration: consideration, rate: rate})
	}

	return lbtt
}

func (lbtt *LBTT) WithFinalBand(start_range_inclusive, consideration, rate float64) *LBTT {
	if !lbtt.isBuilt {
		lbtt.addBand(Band{start_range_inclusive: start_range_inclusive, end_range_inclusive: (start_range_inclusive + consideration), consideration: consideration, rate: rate})
	}
	return lbtt

}
func (lbtt *LBTT) bandsInOrder() error {
	bands := lbtt.bands

	var previousBand Band
	for i, calculatableBand := range bands {
		band, ok := calculatableBand.(Band)
		if !ok {
			return errors.New("One of your provided bands could not be type asserted to the concrete type Band.")
		}

		if i == 0 {
			// Do not need to validate the first band, just assign it as previousBand.
			previousBand = band
		} else {
			// Validate the current band using the previous band.
			err := band.CheckValidBand(previousBand)
			if err != nil {
				return err
			}

			// Update previousBand to the current band for the next iteration.
			previousBand = band
		}
	}

	return nil

}

func (lbtt *LBTT) Build(err *error) *LBTT {
	// Function is responsible for validating the given bands and ensuring that there is a final "boundless band"

	// assert bands are in order
	if err != nil {
		*err = lbtt.bandsInOrder()
	}

	if err != nil {
		lbtt.isBuilt = true
	}
	return lbtt

}