package main 

import (
	"fmt"
)

func NewLBTT() *LBTT {
	return &LBTT{}
}

func (lbtt *LBTT) addBand(band Band) *LBTT {
	lbtt.bands = append(lbtt.bands, band)
	return lbtt
}

func (lbtt *LBTT) WithBand(start_range_inclusive, end_range_inclusive, consideration, rate float64) *LBTT {
	return lbtt.addBand(Band{start_range_inclusive: start_range_inclusive, end_range_inclusive: end_range_inclusive, consideration: consideration, rate: rate})
}

func (lbtt *LBTT) WithFinalBand(start_range_inclusive, consideration, rate float64) *LBTT {
	return lbtt.addBand(Band{start_range_inclusive: start_range_inclusive, end_range_inclusive: (start_range_inclusive + consideration), consideration: consideration, rate: rate})
}

func (lbtt *LBTT) bandsInOrder() error {
	bands := lbtt.bands

	previous_max := 0.0
	for i, band := range bands {

		// New min should be 1 greater than previous max
		if i == 0 {
			// do not need to validate in this case
			previous_max = band.end_range_inclusive
			continue
		}

		if band.start_range_inclusive+1 != previous_max {
			return fmt.Errorf("Found non contiguous bands. Band with starting range %f does not coincide with previous band with ending range %f", band.end_range_inclusive, previous_max)
		}

		if band.consideration != band.end_range_inclusive - previous_max || band.end_range_inclusive == 0{
			return fmt.Errorf("Found invalid conideration given. Considerations should not exceed the width of the band itself. Given consideration %f should not exceed %f", band.consideration, band.end_range_inclusive - previous_max)
		}

		previous_max = band.end_range_inclusive
	}

	return nil

}

func (lbtt *LBTT) Build(err *error) *LBTT{
	// Function is responsible for validating the given bands and ensuring that there is a final "boundless band"

	// assert bands are in order
	*err = lbtt.bandsInOrder()
	// do we need to check for final band?

	if err != nil{
		lbtt.isBuilt = true
	}
	return lbtt

}

func MakeApr2021LBTT() (*LBTT) {
	var err error
	apr2021LBTT := NewLBTT().
		WithBand(0, 145000, 145000, 0).
		WithBand(145001, 250000, 105000, 0.02).
		WithBand(250001, 325000, 75000, 0.05).
		WithBand(325001, 750000, 425000, 0.1).
		WithFinalBand(750000, 125000, 0.12).
		Build(&err)
	return apr2021LBTT
}