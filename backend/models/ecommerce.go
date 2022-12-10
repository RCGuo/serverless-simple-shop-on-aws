package models

type Product struct {
	ProductId   string  `dynamodbav:"productId"   json:"productId"`
	Name        string  `dynamodbav:"name"        json:"name"`
	Description string  `dynamodbav:"description" json:"description"`
	ImageFile   string  `dynamodbav:"imageFile"   json:"imageFile"`
	Category    string  `dynamodbav:"category"    json:"category"`
	Price       float64 `dynamodbav:"price"       json:"price"`
	Company     string  `dynamodbav:"company"     json:"company"`
	Rating      float64 `dynamodbav:"rating"      json:"rating"`
	Topic       string  `dynamodbav:"topic"       json:"topic"`
}

type Products []Product

type ProductKey struct {
	ProductId string `dynamodbav:"productId"   json:"productId"`
}

type ProductKeys []ProductKey

type UpdateProductFavorite struct {
	UserId    string `dynamodbav:"userId"    json:"userId"`
	ProductId string `dynamodbav:"productId" json:"productId"`
	Favorite  bool   `dynamodbav:"favorite"  json:"favorite"`
}

type UpdateProductFavorites []ProductFavorite

type ProductFavorite struct {
	UserId    string `dynamodbav:"userId"    json:"userId"`
	ProductId string `dynamodbav:"productId" json:"productId"`
}

type ProductFavorites []ProductFavorite

type ProductUpdate struct {
	Name        string  `dynamodbav:":name"        json:"name"`
	Description string  `dynamodbav:":description" json:"description"`
	ImageFile   string  `dynamodbav:":imageFile"   json:"imageFile"`
	Category    string  `dynamodbav:":category"    json:"category"`
	Price       float64 `dynamodbav:":price"       json:"price"`
}

type Order struct {
	PaymentIntentId string        `dynamodbav:"paymentIntentId" json:"paymentIntentId"`
	UserId          string        `dynamodbav:"userId"          json:"userId"`
	OrderDate       string        `dynamodbav:"orderDate"       json:"orderDate"`
	OrderId         string        `dynamodbav:"orderId"         json:"orderId"`
	Items           CartItems     `dynamodbav:"items"           json:"items"`
	Email           string        `dynamodbav:"email"           json:"email"`
	Address         string        `dynamodbav:"address"         json:"address"`
	PaymentMethod   PaymentMethod `dynamodbav:"paymentMethod"   json:"paymentMethod"`
	PaymentStatus   string        `dynamodbav:"paymentStatus"   json:"paymentStatus"`
	ShippingFee     string        `dynamodbav:"shippingFee"     json:"shippingFee"`
	Subtotal        string        `dynamodbav:"subtotal"        json:"subtotal"`
	Total           string        `dynamodbav:"total"           json:"total"`
}

type Orders []Order

type PaymentMethod struct {
	Type  string `dynamodbav:"type"  json:"type"`
	Brand string `dynamodbav:"brand" json:"brand"`
}

type CartItem struct {
	UserId         string  `dynamodbav:"userId"         json:"userId"`
	Quantity       int     `dynamodbav:"quantity"       json:"quantity"`
	Price          float64 `dynamodbav:"price"          json:"price"`
	ProductId      string  `dynamodbav:"productId"      json:"productId"`
	ExpirationTime int64   `dynamodbav:"expirationTime" json:"expirationTime"`
}

type CartItems []CartItem

type CartItemUpdate struct {
	ProductId string `dynamodbav:"productId" json:"productId"`
	Quantity  int    `dynamodbav:"quantity"  json:"quantity"`
}

type CartItemDelete struct {
	ProductId string `dynamodbav:"productId" json:"productId"`
}