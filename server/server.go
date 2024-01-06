package server

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/pos/data"
	"github.com/tecnologer/pos/db"
	"github.com/tecnologer/pos/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	// Handle creation logic here
	fmt.Fprintf(w, "Create endpoint hit")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchDes := r.URL.Query().Get("s")
	includeSoldOutStr := r.URL.Query().Get("include_sold_out")
	user := r.URL.Query().Get("user")

	err := data.ValidateUser(user)
	if err != nil {
		logrus.Infof("unauthorized request by %s", user)
		sendErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	includeSoldOut, _ := strconv.ParseBool(includeSoldOutStr)

	products, err := data.NewProductData(db.Connection()).Search(searchDes, includeSoldOut)
	if err != nil {
		logrus.WithError(err).Error("search error")
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(products)
	if err != nil {
		logrus.WithError(err).Error("search marshal error")
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Infof(
		"result %d with search s=%s and sold_out=%s by %s",
		products.Len(),
		searchDes,
		includeSoldOutStr,
		user,
	)

	// Handle search logic here with searchDes variable
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(result)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
	}
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	var chk *models.Sale
	// parse the request body
	err := json.NewDecoder(r.Body).Decode(&chk)
	if err != nil {
		logrus.Error("checkout decode error", err)
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Infof("checkout request by %s for a total %f with %s", chk.User, chk.TotalPaid, chk.PaymentMethod)

	cnn := db.Begin()

	//chk.User = strings.ToLower(chk.User)

	err = data.NewCheckoutData(cnn).Checkout(chk)
	if err != nil {
		logrus.Error("checkout error", err)
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		_ = db.Rollback(cnn)
		return
	}

	_ = db.Commit(cnn)

	// Handle checkout logic here with data variable
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write([]byte(`{"message": "Sale successful"}`))
	if err != nil {
		logrus.Error("checkout write result error", err)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./pages/index.html")
}

// itemHandler handler func of /itemHandler endpoint, returns the quantity of a product by id from db
func itemHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the url
	idStr := r.URL.Path[len("/item/"):]

	// convert the id to int
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logrus.Error("itemHandler id conversion error", err)
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := r.URL.Query().Get("user")

	err = data.ValidateUser(user)
	if err != nil {
		logrus.Infof("check_qty: unauthorized request by %s", user)
		sendErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// get the product from db
	product, err := data.NewProductData(db.Connection()).Get(uint(id))
	if err != nil {
		logrus.Error("itemHandler get product error", err)
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse product to json and return it
	result, err := json.Marshal(product)
	if err != nil {
		logrus.Error("itemHandler marshal error", err)
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Infof("check_qty: result %d with id=%d by %s", product.ID, id, user)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(result)
	if err != nil {
		logrus.Error("itemHandler write result error", err)
	}
}

func newItemHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/create.html")
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	// get user from url
	user := r.URL.Query().Get("user")

	err := data.ValidateUser(user)
	if err != nil {
		logrus.Infof("create: unauthorized request by %s", user)
		sendErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// parse the request body
	var productsCsv string
	err = json.NewDecoder(r.Body).Decode(&productsCsv)
	if err != nil {
		logrus.Error("create decode error", err)
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cnn := db.Begin()

	rows := strings.Split(productsCsv, "\n")
	for _, row := range rows {
		cols := strings.Split(row, ",")
		if len(cols) != 3 {
			logrus.Error("create invalid row", row)
			sendErrorResponse(w, "invalid row", http.StatusBadRequest)
			return
		}

		qty, err := strconv.ParseUint(cols[1], 10, 32)
		if err != nil {
			logrus.Error("create invalid qty", cols[1])
			sendErrorResponse(w, "invalid qty", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(cols[2], 64)
		if err != nil {
			logrus.Error("create invalid price", cols[2])
			sendErrorResponse(w, "invalid price", http.StatusBadRequest)
			return
		}

		product := &models.Product{
			Description: cols[0],
			Qty:         uint(qty),
			Price:       float32(price),
		}

		err = data.NewProductData(cnn).Create(product)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error("create error", err)
			sendErrorResponse(w, err.Error(), http.StatusBadRequest)
			_ = db.Rollback(cnn)
			return
		}

	}

	_ = db.Commit(cnn)
}

func Run(port int) {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/checkout", checkoutHandler)
	http.HandleFunc("/item/", itemHandler)
	http.HandleFunc("/newitem", newItemHandler)
	http.HandleFunc("/create", createItemHandler)
	http.HandleFunc("/", serveRoot)

	addr := fmt.Sprintf(":%d", port)

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func sendErrorResponse(w http.ResponseWriter, errMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := map[string]string{"error": errMsg}
	json.NewEncoder(w).Encode(errorResponse)
}
