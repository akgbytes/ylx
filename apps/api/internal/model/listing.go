package model

import "time"

type Listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}
