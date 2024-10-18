package calc 

import (
	"fmt"
	"testing"
)

func assertCalculationCorrect(calculator HousingTaxCalculator, house_price float64, expected_tax float64) (bool, string) {
	tax, _ := calculator.Calculate(house_price)
	if tax != expected_tax {
		return false, fmt.Sprintf("Incorrect tax for houseprice %f. Got %f but expcted %f", house_price, tax, expected_tax)
	}

	return true, ""

}

func assertApril2021Tax(house_price float64, expected_tax float64) (bool, string) {
	fact := Apr2021Factory{}
	lbtt := fact.createCalculator()
	return assertCalculationCorrect(lbtt, house_price, expected_tax)
}

func Test_Zero_Band(t *testing.T) {

	valid, result := assertApril2021Tax(144999.0, 0.0)
	if !valid {
		t.Error(result)
	}
}

func Test_First_Tax_Band(t *testing.T) {
	valid, result := assertApril2021Tax(150000.0, 100.0)
	if !valid {
		t.Error(result)
	}
}

func Test_First_And_Second_Bands(t *testing.T) {
	valid, result := assertApril2021Tax(300000.0, 4600.0)
	if !valid {
		t.Error(result)
	}
}

func Test_Max_Band(t *testing.T) {
	valid, result := assertApril2021Tax(760000.0, 49550.0)
	if !valid {
		t.Error(result)
	}
}

func Test_Exceed_Max_Band_And_Consideration(t *testing.T) {
	valid, result := assertApril2021Tax(830000.0, 57950.0)
	if !valid {
		t.Error(result)
	}
	valid, result = assertApril2021Tax(875000.0, 63350.0)
	if !valid {
		t.Error(result)
	}
}

func Test_Additional_Dwelling_Supplement(t *testing.T) {
	fact := Apr2021Factory{}
	lbtt := fact.createCalculator()
	calc := lbttWithAdditionalDwelling{
		lbttCalculator:    lbtt, 
		additionalDwelling: 200000,
	}

	valid, result := assertCalculationCorrect(calc, 250000, 14100)
	if !valid {
		t.Error(result)
	}
}

func Test_Invalid_Aggregation_Of_Calculators(t *testing.T) {
	var err error;
	NewAggregatedCalculator().WithFTBCalculator(175000).WithAdditionDwellingCalculator(200000).Build(&err)
	if err == nil{
		t.Error("Should not be able to make an aggregated calculator with FTB where ADS applies.")
	}
}

func Test_Valid_Aggregration_Of_Calculators(t *testing.T){
	var err error;
	NewAggregatedCalculator().WithFTBCalculator(175000).Build(&err)
	if err != nil{
		t.Error("Creation of aggregate calculator with just a FTB calculator should not error.")
	}
	

}
func Test_First_Time_Buyers(t *testing.T) {
	fact := Apr2021Factory{}
	lbtt := fact.createCalculator()
	calc := lbttWithFirstTimeBuyersRelief{
		lbttCalculator: lbtt,
		ftbNilRateBand: 175000,
	}

	valid, result := assertCalculationCorrect(calc, 175000, 0)
	if !valid {
		t.Error(result)
	}

	valid, result = assertCalculationCorrect(calc, 180000, 700)
	if !valid {
		t.Error(result)
	}
}

func Test_Build_Invalid_Ranges(t *testing.T) {
	var err error
	NewLBTT().
		WithBand(100000, 130000, 30000, 0.1).
		WithBand(99999, 129999, 30000, 0.2).
		WithBand(0, 100, 100, 0.02).
		Build(&err)
	if err == nil {
		t.Error("Expected an error when trying to build an LBTT object with invalid bands.")
	}
}

func Test_Build_Invalid_Considerations(t *testing.T) {
	var err error
	NewLBTT().
		WithBand(100000, 130000, 1000, 0.1).
		WithBand(99999, 129999, 2000, 0.2).
		WithBand(0, 100, 10, 0.02).
		Build(&err)
	if err == nil {
		t.Error("Expected an error when trying to build an LBTT object with invalid considerations.")
	}
}
