package data

import (
	"github.com/pkg/errors"
	"github.com/tecnologer/pos/models"
	"gorm.io/gorm"
)

type CheckoutData struct {
	cnn *gorm.DB
}

func NewCheckoutData(cnn *gorm.DB) *CheckoutData {
	return &CheckoutData{cnn: cnn}
}

func (cd *CheckoutData) Checkout(chk *models.Sale) (err error) {
	if err := Validate(chk); err != nil {
		return errors.Wrap(err, "checkout: error validating checkout")
	}

	items := append([]*models.Product{}, chk.Items...)
	chk.Items = nil

	err = cd.cnn.Save(chk).Error
	if err != nil {
		return errors.Wrap(err, "checkout: error saving checkout")
	}

	pd := NewProductData(cd.cnn)

	for _, product := range items {
		err = pd.Sell(chk.ID, product)
		if err != nil {
			return errors.Wrap(err, "checkout: error selling product")
		}
	}

	return nil
}
