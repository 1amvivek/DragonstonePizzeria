package main

type Product struct {
	Id    int    `json: "id"`
	Name  string `json: "name"`
	Price int    `json: "price"`
}

type Cart struct {
	SerialNumber string    `json: "SerialNumber"`
	Products     []Product `json: "products"`
	Clock        int       `json: "clock"`
}
