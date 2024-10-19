package main

import (
	"database/sql"
	"errors"

	"log"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Get all products
func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price FROM products"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []product
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error in row iteration: %v", err)
		return nil, err
	}

	return products, nil
}

// Get a single product by ID
func (p *product) getProduct(db *sql.DB) error {
	query := "SELECT name, quantity, price FROM products WHERE id = ?"
	row := db.QueryRow(query, p.ID)

	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		log.Printf("Error fetching product with ID %d: %v", p.ID, err)
		return err
	}

	return nil
}

//Create product
// Create product
func (p *product) createProduct(db *sql.DB) error {
	// Prepare the insert query using placeholders
	query := "INSERT INTO products (name, quantity, price) VALUES (?, ?, ?)"

	// Use a prepared statement to execute the query
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)
	if err != nil {
		return err
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Assign the ID to the product
	p.ID = int(id)
	return nil
}

//update product method
func (p *product) updateProduct(db *sql.DB) error{
	query:="UPDATE products SET name=?,quantity=?, price=? WHERE id=?"
	result,err:=db.Exec(query,p.Name,p.Quantity,p.Price,p.ID)
	if err!=nil{
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err!=nil{
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No such row exist")
	}
	
	return nil
}


//DELETE PRODUCT
func (p *product) deleteProduct(db *sql.DB) error{
	query:= "DELETE FROM products WHERE id=?"
	_,err:=db.Exec(query,p.ID)
	return err
}