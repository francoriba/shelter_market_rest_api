// app/models/offer_model.go

package models

// Offer model represents an offer/product in the system
type Offer struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"type:varchar(100);not null" json:"name"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"`
	Category string  `gorm:"type:varchar(50);not null" json:"category"`
}

// OfferResponse defines the structure of the response for the GetOffers endpoint
type OfferResponse struct {
	Code    int     `json:"code"`
	Message []Offer `json:"message"`
}

// SuppliesResponse defines the structure of the response from the HPCPP /supplies endpoint
type SuppliesResponse struct {
	Food     map[string]int `json:"food"`
	Medicine map[string]int `json:"medicine"`
}
