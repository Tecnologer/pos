package models

type Product struct {
	ID          uint    `json:"id"          gorm:"primarykey" csv:"id"`
	Description string  `json:"description"                   csv:"description" validate:"required"`
	Qty         uint    `json:"qty"                           csv:"qty"         validate:"gte=0"`
	Price       float32 `json:"price"                         csv:"price"       validate:"gte=0"`
}

func (p *Product) IsInit() bool {
	return p != nil && p.ID != 0
}

type Products []*Product

func (p *Products) Len() int {
	if p == nil {
		return 0
	}

	return len(*p)
}

func (p *Products) Less(i, j int) bool {
	if p.Len() == 0 {
		return false
	}

	return (*p)[i].Description < (*p)[j].Description
}

func (p *Products) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}
