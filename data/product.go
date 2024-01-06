package data

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/pos/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"sort"
)

type ProductData struct {
	cnn *gorm.DB
}

func NewProductData(cnn *gorm.DB) *ProductData {
	return &ProductData{cnn: cnn}
}

func (pd *ProductData) Create(product *models.Product) error {
	if err := Validate(product); err != nil {
		return errors.Wrap(err, "product: error validating data")
	}

	product.Description = title(product.Description)

	existingProduct, err := pd.GetByDescription(product.Description)
	if err != nil {
		return errors.Wrap(err, "product: error getting product by description")
	}

	if existingProduct.IsInit() {
		product.ID = existingProduct.ID
	}

	return pd.cnn.Save(product).Error
}

func (pd *ProductData) GetByDescription(description string) (*models.Product, error) {
	var product models.Product

	err := pd.cnn.Where("description = ?", description).First(&product).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return nil, errors.Wrap(err, "product: error getting product by description")
	}

	return &product, nil
}

func (pd *ProductData) Search(searchValue string, includeSoldOut bool) (*models.Products, error) {
	tx := pd.cnn
	if searchValue != "" {
		tx = tx.Where("description LIKE ? OR ID = ?", fmt.Sprintf("%%%s%%", searchValue), searchValue)
	}

	if !includeSoldOut {
		tx = tx.Where("qty > 0")
	}

	var products *models.Products

	err := tx.Find(&products).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "product: error searching products")
	}

	sort.Sort(products)

	return products, nil
}

func (pd *ProductData) Sell(checkoutID uint, product *models.Product) error {
	// Update the product qty
	tx := pd.cnn.Model(&models.Product{}).
		Where("id = ?", product.ID).
		Update("qty", gorm.Expr("qty - ?", product.Qty))
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "product: error selling product")
	}

	// register the sale details
	saleDetails := &models.SaleDetails{
		SaleID:    checkoutID,
		ProductID: product.ID,
		Qty:       product.Qty,
		Price:     product.Price,
	}

	err := pd.cnn.Save(saleDetails).Error
	if err != nil {
		return errors.Wrap(err, "product: error selling product")
	}

	return nil
}

func (pd *ProductData) Get(id uint) (*models.Product, error) {
	var product *models.Product

	err := pd.cnn.First(&product, id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "product.get: error getting product")
	}

	return product, nil
}

func title(s string) string {
	return cases.Title(language.Spanish).String(s)
}
