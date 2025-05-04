// Generate User Model
package models

import "time"

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
