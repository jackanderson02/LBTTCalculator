package main

import(
	"testing"
)

func Test_Zero_Band(t *testing.T){

	lbtt := MakeApr2021LBTT()
	tax, _:= lbtt.Calculate(144999)
	if (tax != 0){
		t.Errorf("Non zero tax for given price outside of minimum band.")
	}

}

func Test_First_Tax_Band(t *testing.T){
	lbtt := MakeApr2021LBTT()
	tax, _:= lbtt.Calculate(150000)
	if (tax != 100){
		t.Errorf("Incorrect tax for house price %f. Got %f but expected %f", 150000.0, tax, 100.0);
	}
}

func Test_First_And_Second_Bands(t *testing.T){
	lbtt := MakeApr2021LBTT()
	tax, _ := lbtt.Calculate(300000)
	if (tax != 4600){
		t.Errorf("Incorrect tax for house price %f. Got %f but expected %f", 300000.0, tax, 4600.0);
	}

}

func Test_Max_Band(t *testing.T){
	lbtt := MakeApr2021LBTT()
	tax, _ := lbtt.Calculate(760000)
	if (tax != 49550){
		t.Errorf("Incorrect tax for house price %f. Got %f but expected %f", 760000.0, tax, 49550.0);
	}
}

func Test_Exceed_Max_Band_And_Consideration(t *testing.T){
	lbtt := MakeApr2021LBTT()
	tax, _:= lbtt.Calculate(830000)
	if (tax != 57950){
		t.Errorf("Incorrect tax for house price %f. Got %f but expected %f", 830000.0, tax, 57950.0);
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


