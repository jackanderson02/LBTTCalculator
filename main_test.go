package main

import(
	"testing"
	"fmt"
)

func assertTaxEqual(house_price float64, expected_tax float64) (bool, string){
	lbtt := MakeApr2021LBTT()
	tax, _:= lbtt.Calculate(house_price)
	if tax != expected_tax{
		return false, fmt.Sprintf("Incorrect tax for houseprice %f. Got %f but expcted %f", house_price, tax, expected_tax);
	}
	return true, ""
}

func Test_Zero_Band(t *testing.T){

	valid, result := assertTaxEqual(144999.0, 0.0)
	if !valid{
		t.Error(result)
	}
}

func Test_First_Tax_Band(t *testing.T){
	valid, result := assertTaxEqual(150000.0, 100.0)
	if !valid{
		t.Error(result)
	}
}

func Test_First_And_Second_Bands(t *testing.T){
	valid, result := assertTaxEqual(300000.0, 4600.0)
	if !valid{
		t.Error(result)
	}
}

func Test_Max_Band(t *testing.T){
	valid, result := assertTaxEqual(760000.0, 49550.0)
	if !valid{
		t.Error(result)
	}
}

func Test_Exceed_Max_Band_And_Consideration(t *testing.T){
	valid, result := assertTaxEqual(830000.0, 57950.0)
	if !valid{
		t.Error(result)
	}
	valid, result = assertTaxEqual(875000.0, 63350.0)
	if !valid{
		t.Error(result)
	}
}

func Test_Build_Invalid_Ranges(t *testing.T){
	var err error;
	NewLBTT().
		WithBand(100000, 130000, 30000, 0.1).
		WithBand(99999, 129999, 30000, 0.2).
		WithBand(0, 100, 100, 0.02).
		Build(&err)
	if err == nil{
		t.Error("Expected an error when trying to build an LBTT object with invalid bands.")
	}
}


func Test_Build_Invalid_Considerations(t *testing.T){
	var err error;
	NewLBTT().
		WithBand(100000, 130000, 1000, 0.1).
		WithBand(99999, 129999, 2000, 0.2).
		WithBand(0, 100, 10, 0.02).
		Build(&err)
	if err == nil{
		t.Error("Expected an error when trying to build an LBTT object with invalid considerations.")
	}
}


