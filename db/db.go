package db

import (
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/pkg/errors"
	"github.com/tecnologer/pos/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var (
	m   sync.Mutex
	cnn *gorm.DB
)

func Connect() error {
	m.Lock()
	defer m.Unlock()

	if cnn != nil {
		return nil
	}

	c, err := gorm.Open(sqlite.Open("items.db"), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		return errors.Wrap(err, "error connecting to database")
	}

	cnn = c

	return nil
}

func Connection() *gorm.DB {
	if cnn == nil {
		_ = Connect()
	}

	return cnn
}

func Migrate() error {
	err := Connect()
	if err != nil {
		return errors.Wrap(err, "error connecting to database")
	}

	err = cnn.AutoMigrate(&models.Product{}, &models.Sale{}, &models.SaleDetails{})
	if err != nil {
		return errors.Wrap(err, "error migrating database")
	}

	return nil
}

func Begin() *gorm.DB {
	return Connection().Begin()
}

func Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}
