package models

import "time"

type Sale struct {
	ID            uint       `json:"id"             gorm:"primarykey"`
	Items         []*Product `json:"items"          gorm:"many2many:product_checkouts" validate:"required,dive,required"`
	GrandTotal    float32    `json:"grand_total"                                       validate:"gte=0"`
	PaymentMethod string     `json:"payment_method"                                    validate:"oneof=cash card"`
	User          string     `json:"user"                                              validate:"oneof=carlitos chio tecnologer"`
	TotalPaid     float32    `json:"total_paid"                                        validate:"gte=0"`
	Date          time.Time  `json:"date"           gorm:"default:CURRENT_TIMESTAMP;"`
}

type SaleDetails struct {
	ID        uint     `json:"id"         gorm:"primarykey"`
	SaleID    uint     `json:"sale_id"`
	Sale      *Sale    `json:"sale"       gorm:"foreignKey:SaleID"`
	ProductID uint     `json:"product_id"`
	Product   *Product `json:"product"    gorm:"foreignKey:ProductID"`
	Qty       uint     `json:"qty"                                    validate:"gte=0"`
	Price     float32  `json:"price"                                  validate:"gte=0"`
}
