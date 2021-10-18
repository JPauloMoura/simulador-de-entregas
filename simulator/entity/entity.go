package entity

// Order possui informações do pedido gerado
type Order struct {
	UUID        string `json:"order_id"`
	Destination string `json:"destination"`
}

// Destination possui informações da localização de onde o pedido deve ser entregue
type Destination struct {
	Order     string `json:"order_id"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
}
