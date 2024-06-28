// pkg/utils/cron_job.go

package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/robfig/cron/v3"
)

const (
	PriceFruits      = 2
	PriceMeat        = 4
	PriceVegetables  = 1
	PriceWater       = 1
	PriceAnalgesics  = 5
	PriceAntibiotics = 9
	PriceBandages    = 4

	CategoryFood     = "food"
	CategoryDrink    = "drink"
	CategoryMedicine = "medicine"

	SuppliesURL = "http://192.168.0.57:8011/supplies?id=latest"
)

func FetchAndStoreSupplies() {
	resp, err := http.Get(SuppliesURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch supplies: %v", err)
		return
	}
	defer resp.Body.Close()

	var supplies models.SuppliesResponse
	if err := json.NewDecoder(resp.Body).Decode(&supplies); err != nil {
		log.Printf("Failed to decode supplies response: %v", err)
		return
	}

	log.Println("Successfully fetched supplies")
	db := database.GetDB()

	offers := []models.Offer{
		{Name: "fruits", Quantity: supplies.Food["fruits"] / 5, Price: PriceFruits, Category: CategoryFood},
		{Name: "meat", Quantity: supplies.Food["meat"] / 5, Price: PriceMeat, Category: CategoryFood},
		{Name: "vegetables", Quantity: supplies.Food["vegetables"] / 5, Price: PriceVegetables, Category: CategoryFood},
		{Name: "water", Quantity: supplies.Food["water"] / 5, Price: PriceWater, Category: CategoryDrink},
		{Name: "analgesics", Quantity: supplies.Medicine["analgesics"] / 5, Price: PriceAnalgesics, Category: CategoryMedicine},
		{Name: "antibiotics", Quantity: supplies.Medicine["antibiotics"] / 5, Price: PriceAntibiotics, Category: CategoryMedicine},
		{Name: "bandages", Quantity: supplies.Medicine["bandages"] / 5, Price: PriceBandages, Category: CategoryMedicine},
	}

	for _, offer := range offers {
		db.Where(models.Offer{Name: offer.Name}).Assign(offer).FirstOrCreate(&offer)
	}
}

func StartCronJob() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", FetchAndStoreSupplies)
	if err != nil {
		log.Fatalf("Error starting cron job: %v", err)
	}
	c.Start()
}
