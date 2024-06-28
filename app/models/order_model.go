// app/models/order_model.go

package models

import "gorm.io/gorm"

// Order model represents an order in the system
type Order struct {
	gorm.Model              // Embeds fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	Status      string      `json:"status" gorm:"not null"` // Status of the order (e.g., "processing", "completed")
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`     // Relation to order items
	TotalAmount float64     `json:"total_amount"`           // Total amount of the order
}

// OrderItem model represents an item in an order
type OrderItem struct {
	gorm.Model
	OrderID  uint    `json:"order_id" gorm:"not null"` // Foreign key to orders table
	OfferID  uint    `json:"offer_id" gorm:"not null"` // Foreign key to offers table
	Quantity int     `json:"quantity" gorm:"not null"` // Quantity of the item
	Offer    Offer   `gorm:"foreignKey:OfferID"`       // Relation to the offer
	SubTotal float64 `json:"sub_total"`                // Subtotal for the item (price * quantity)
}

// CheckoutRequest defines the structure of the request for the Checkout endpoint
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items" validate:"required"`
}

// CheckoutItem defines the structure of each item in the CheckoutRequest
type CheckoutItem struct {
	OfferID  uint `json:"offer_id" validate:"required"`
	Quantity int  `json:"quantity" validate:"required,min=1"`
}

// CheckoutResponse defines the structure of the response for the Checkout endpoint
type CheckoutResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	OrderID uint   `json:"order_id"`
}

// OrderResponse defines the structure of the response for the GetOrderStatus endpoint
type OrderResponse struct {
	Code    int                    `json:"code"`
	Message map[string]interface{} `json:"message"`
}
