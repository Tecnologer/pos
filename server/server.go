package server

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/pos/data"
	"github.com/tecnologer/pos/db"
	"github.com/tecnologer/pos/models"
	"log"
	"net/http"
	"strconv"
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

func Run(port int) {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/checkout", checkoutHandler)
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
