package main

import(
	"errors"
)

type LBTT struct {
	bands   []Band
	isBuilt bool
}

type Band struct {
	start_range_inclusive, end_range_inclusive, consideration, rate float64
}


func (lbtt LBTT) Caclulate(house_price float64) (float64, error){
	if (!lbtt.isBuilt){
		return 0.0, errors.New("Cannot calculate house prices until the lbtt object is built. Please call .build() to validate the given bands.")
	}
	// Goes through each of the bands and determines the tax
	var total float64 = 0
	// min_band := 0
	for _, band := range lbtt.bands{
		band_end := band.end_range_inclusive
		// How far you are into each band * the rate. The max you can be into a band is the consideration
		band_diff := house_price - (band_end - band.consideration)
		if (band_diff >= 0){ // In or exceeding this band
			if band_diff > band.consideration{
				band_diff = band.consideration // Do not exceed the consideration of this band
			}
			total += band_diff * band.rate
		}
	}

	return total, nil

}