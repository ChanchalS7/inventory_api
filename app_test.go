package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App 

func TestMain(m *testing.M){
	err := a.Initialise(DBUser,DBPassword,"test_go")
	

	if err !=nil{
		log.Fatal("Error occurred while initialising the database")
	}
	createTable()
	m.Run()//run all other test withing package

}

func createTable(){
	createTableQuery:=`CREATE TABLE IF NOT EXISTS products(
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	quantity INT,
	price FLOAT(10,7),
	PRIMARY KEY(id)
	);`
	_,err:=a.DB.Exec(createTableQuery)
	if err!=nil{
		log.Fatal(err)
	}
}

func clearTable(){
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER TABLE products AUTO_INCREMENT=1")
	log.Println("Clear table")
}
func addProduct(name string, quantity int, price float64){
	// Prepare the insert query using placeholders
	query := "INSERT INTO products (name, quantity, price) VALUES (?, ?, ?)"

	_,err:= a.DB.Exec(query,name,quantity,price)
	if err!=nil{
		log.Println(err)
	}
}
func TestGetProduct(t *testing.T){
clearTable()
addProduct("keyboard",100,500)
request,_:=http.NewRequest("GET","/product/1",nil)
response :=sendRequest(request)
checkStatusCode(t,http.StatusOK,response.Code)



}
func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int){
 if expectedStatusCode != actualStatusCode{
	t.Errorf("Expected status : %v, Received:%v", expectedStatusCode,actualStatusCode)
 }
}
func sendRequest(request *http.Request) *httptest.ResponseRecorder{
	recorder:=httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}


func TestCreateProduct(t *testing.T) {
    clearTable()
    var product = []byte(`{"name":"chair","quantity":1,"price":100}`)
    req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
    req.Header.Set("Content-Type", "application/json")

    response := sendRequest(req)
    checkStatusCode(t, http.StatusCreated, response.Code)

    // if response.Code != http.StatusCreated {
    //     t.Errorf("Expected status code 201 but got %v. Response body: %s", response.Code, response.Body.String())
    // }

		//unmarshal json response
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&m)

	if m["name"] != "chair"{
		t.Errorf("Expected name:%v,Got :%v","chair",m["name"])
	}
	if m["quantity"]!=1.0{
		t.Errorf("Expected quantity:%v,Got:%v",1.0,m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
    clearTable()

    // Add a product to be deleted
    addProduct("connector", 10, 10)

    // Step 1: Check if the product exists
    req, _ := http.NewRequest("GET", "/product/1", nil)
    response := sendRequest(req)
    checkStatusCode(t, http.StatusOK, response.Code)

    // Step 2: Send DELETE request to delete the product
    req, _ = http.NewRequest("DELETE", "/product/1", nil)
    response = sendRequest(req)
    checkStatusCode(t, http.StatusOK, response.Code)  // Expect 200 OK for successful deletion

    // Step 3: Check if the product is gone
    req, _ = http.NewRequest("GET", "/product/1", nil)
    response = sendRequest(req)
    checkStatusCode(t, http.StatusNotFound, response.Code)  // Expect 404 Not Found after deletion
}

func TestUpdateProduct(t *testing.T) {
    clearTable()

    // Add the initial product
    addProduct("connector", 10, 10)

    // Step 1: Send PUT request to update the product
    var product = []byte(`{"name":"connector_updated","quantity":5,"price":20}`)
    req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
    req.Header.Set("Content-Type", "application/json")

    response := sendRequest(req)

    // Step 2: Check the status code of the update operation
    checkStatusCode(t, http.StatusOK, response.Code)

    // Step 3: Unmarshal response and verify the updated values
    var updatedValue map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &updatedValue)

    if updatedValue["id"] != float64(1) { // JSON unmarshalling converts numbers to float64 by default
        t.Errorf("Expected id: %v, Got: %v", 1, updatedValue["id"])
    }
    if updatedValue["name"] != "connector_updated" {
        t.Errorf("Expected name: %v, Got: %v", "connector_updated", updatedValue["name"])
    }
    if updatedValue["price"] != 20.0 {
        t.Errorf("Expected price: %v, Got: %v", 20.0, updatedValue["price"])
    }
    if updatedValue["quantity"] != 5.0 {
        t.Errorf("Expected quantity: %v, Got: %v", 5.0, updatedValue["quantity"])
    }
}
