package entities

import "time"

type Product struct {
	Id        	uint
	Name      	string
	CategoryId	uint
	Image		string
	Stock		int
	Description	string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	Category	Category
}