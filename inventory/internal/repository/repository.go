package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/miha3009/market/inventory/internal/model"
)

type InventoryRepository interface {
	Avaliable(ids []int32) ([]bool, error)
	Reserve(req []model.ReserveRequest) (bool, error)
	CancelReserve(req []model.ReserveRequest) error
}

type InventoryRepositoryImpl struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &InventoryRepositoryImpl{
		db: db,
	}
}

func (r *InventoryRepositoryImpl) Avaliable(ids []int32) ([]bool, error) {
	strIds := make([]string, len(ids))
	for i := range ids {
		strIds[i] = strconv.Itoa(int(ids[i]))
	}
	params := "{" + strings.Join(strIds, ",") + "}"
	rows, err := r.db.Query("SELECT id FROM inventory WHERE id = ANY($1::int[]) AND count > reserved", params)
	defer rows.Close()

	data := make(map[int32]bool)
	for rows.Next() {
		var id int32
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		data[id] = true
	}

	res := make([]bool, len(ids))
	for i := range ids {
		_, ok := data[ids[i]]
		res[i] = ok
	}

	return res, err
}

func (r *InventoryRepositoryImpl) Reserve(req []model.ReserveRequest) (bool, error) {
	tx, err := r.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	for i := range req {
		res, err := tx.Exec("UPDATE inventory SET reserved = reserved + $1 WHERE product_id = $2 AND count-reserved >= $1", req[i].Count, req[i].ProductId)
		if err != nil {
			return false, err
		}

		rows, err := res.RowsAffected()
		if rows != 1 {
			return false, err
		}
	}

	return true, tx.Commit()
}

func (r *InventoryRepositoryImpl) CancelReserve(req []model.ReserveRequest) error {
	tx, err := r.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := range req {
		_, err := tx.Exec("UPDATE inventory SET reserved = reserved - $1 WHERE product_id = $2", req[i].Count, req[i].ProductId)
		if err != nil {
			return err
		}

	}

	return tx.Commit()
}
