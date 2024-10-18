package calc 

func MakeApr2021LBTT() *LBTT {
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