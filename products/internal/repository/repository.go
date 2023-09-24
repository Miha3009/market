package repository

import (
	"database/sql"
	"products/internal/model"
)

type ProductRepository interface {
	SelectById(id int) (model.Product, error)
	Select(offset, limit int) ([]model.Product, int, error)
	Create(value model.Product) (int, error)
	Delete(id int) error
	Update(value model.Product) error
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

func (r *ProductRepositoryImpl) SelectById(id int) (model.Product, error) {
	var res model.Product
	err := r.db.QueryRow("SELECT id, name, description, price FROM products WHERE id=$1", id).Scan(&res.Id, &res.Name, &res.Description, &res.Price)
	return res, err
}

func (r *ProductRepositoryImpl) Select(offset, limit int) ([]model.Product, int, error) {
	var res []model.Product
	rows, err := r.db.Query("SELECT id, name, description, price FROM products ORDER BY id OFFSET $1 LIMIT $2", offset, limit)
	defer rows.Close()

	for rows.Next() {
		var temp model.Product
		err = rows.Scan(&temp.Id, &temp.Name, &temp.Description, &temp.Price)
		if err != nil {
			return res, 0, err
		}
		res = append(res, temp)
	}

	count := 0
	err = r.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	return res, count, err
}

func (r *ProductRepositoryImpl) Create(value model.Product) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id", value.Name, value.Description, value.Price).Scan(&id)
	return id, err
}

func (r *ProductRepositoryImpl) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}

func (r *ProductRepositoryImpl) Update(value model.Product) error {
	_, err := r.db.Exec("UPDATE products SET name=$1, description=$2, price=$3 WHERE id=$4", value.Name, value.Description, value.Price, value.Id)
	return err

}
