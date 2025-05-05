package models

type Product struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Sku         string `json:"sku" bson:"sku"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
	UpdatedAt   string `json:"updatedAt" bson:"updatedAt"`
}

// Set setters and getters
func (p *Product) GetID() string {
	return p.ID
}

func (p *Product) SetID(id string) {
	p.ID = id
}

func (p *Product) GetName() string {
	return p.Name
}

func (p *Product) SetName(name string) {
	p.Name = name
}

func (p *Product) GetDescription() string {
	return p.Description
}

func (p *Product) SetDescription(description string) {
	p.Description = description
}

func (p *Product) GetSku() string {
	return p.Sku
}

func (p *Product) SetSku(sku string) {
	p.Sku = sku
}

func (p *Product) GetCreatedAt() string {
	return p.CreatedAt
}
func (p *Product) SetCreatedAt(createdAt string) {
	p.CreatedAt = createdAt
}
func (p *Product) GetUpdatedAt() string {
	return p.UpdatedAt
}
func (p *Product) SetUpdatedAt(updatedAt string) {
	p.UpdatedAt = updatedAt
}
