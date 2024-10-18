package calc

import(
	"reflect"
)

type HousingTaxSubscriber struct{
	updateIfExceeded float64 
	update func(house_price float64)
}

type PublishableCalculator interface{
	CalculateAndObserve(house_price float64) 
	Subscribe(subscriber HousingTaxSubscriber)
	Unsubscribe(subscriber HousingTaxSubscriber)
}

type HousingTaxPublisher struct{
	calculator HousingTaxCalculator
	subscribers []HousingTaxSubscriber
}

func (housingTaxPublisher *HousingTaxPublisher) CalculateAndObserve(house_price float64) {
	tax, err := housingTaxPublisher.calculator.Calculate(house_price)
	if err != nil{
		return 
	}

	// Else notify all subscribers who care
	for _, subscriber := range housingTaxPublisher.subscribers{
		if tax > subscriber.updateIfExceeded{
			subscriber.update(tax)
		}
	}

}

func (housingTaxPublisher *HousingTaxPublisher) Subscribe(sub HousingTaxSubscriber){
	housingTaxPublisher.subscribers = append(housingTaxPublisher.subscribers, sub)
}
func (housingTaxPublisher *HousingTaxPublisher) Unsubscribe(subscriber HousingTaxSubscriber){
	var newSubsribers []HousingTaxSubscriber;
	subscribers := housingTaxPublisher.subscribers


	for _, sub := range subscribers {
		if !reflect.DeepEqual(sub, subscriber){
			newSubsribers= append(newSubsribers, sub)
		} 
	}

	housingTaxPublisher.subscribers = newSubsribers	
}



