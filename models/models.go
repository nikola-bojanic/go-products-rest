package models

import "time"

type Product struct {
	ProductId        int       `json:"productId"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"shortDescription"`
	Description      string    `json:"description"`
	Price            float32   `json:"price"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Quantity         int       `json:"quantity"`
	Category         Category  `json:"category"`
}

type Category struct {
	Id        int       `json:"categoryId"`
	Name      string    `json:"categoryName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
