package main

import (
	"flag"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/pos/data"
	"github.com/tecnologer/pos/db"
	"github.com/tecnologer/pos/models"
	"os"
)

var path = flag.String("path", "items.csv", "Path to csv file")

func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	products, err := readDataFromFile(*path)
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.Connect()
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.Migrate()
	if err != nil {
		logrus.Fatal(err)
	}

	pData := data.NewProductData(db.Connection())

	for _, product := range products {
		err = pData.Create(product)
		if err != nil {
			logrus.Error("create %s. Err: %v", product.Description, err)
		}
	}
}

func readDataFromFile(filePath string) (input []*models.Product, err error) {
	logrus.Debug("loading pav1 data csv")
	defer logrus.Debug("loading pav1 data csv completed")

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, errors.Errorf("pav1: file %s does not exist", filePath)
	}

	clientsFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "pav1: open file %s", filePath)
	}

	defer func(clientsFile *os.File) {
		err := clientsFile.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(clientsFile)

	err = gocsv.UnmarshalFile(clientsFile, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "pav1: marshal data")
	}

	return
}
