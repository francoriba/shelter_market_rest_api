// app/models/responses_model.go

package models

// ErrorResponse defines the structure of an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse defines la estructura de una respuesta exitosa
type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DashboardResponse defines the structure of the response for the dashboard endpoint
type DashboardResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Orders  []OrderDashboard `json:"orders"`
}

// OrderDashboard defines the structure for each order's details in the dashboard response
type OrderDashboard struct {
	ID          uint               `json:"id"`
	Status      string             `json:"status"`
	TotalAmount float64            `json:"total_amount"`
	Items       []OrderItemDetails `json:"items"`
}

// OrderItemDetails defines the structure for each item's details in the order
type OrderItemDetails struct {
	OfferID  uint    `json:"offer_id"`
	Quantity int     `json:"quantity"`
	SubTotal float64 `json:"sub_total"`
}

// UpdateOrderStatusResponse defines the structure of the response for updating the order status
type UpdateOrderStatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// UpdateOrderStatusRequest defines the structure of the request to update the order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=preparing processing shipped delivered"`
}

// UserResponse represents the structure of the response for getting a user
type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GetAllUsersResponse defines the structure of the response for getting all users
type GetAllUsersResponse struct {
	Code    int            `json:"code"`
	Message []UserResponse `json:"message"`
}

// DeleteUserResponse defines the structure of the response for deleting a user
type DeleteUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
