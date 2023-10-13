package service

import (
	"github.com/miha3009/market/products/internal/model"
	"github.com/miha3009/market/products/internal/repository"
	"github.com/miha3009/market/products/pkg/handlers"
)

type ProductService interface {
	SelectById(id int) (model.Product, error)
	Select(offset, limit int) (model.SelectResponse, error)
	Create(value model.Product) (model.Product, error)
	Delete(id int) error
	Update(value model.Product) error
}

type ProductServiceImpl struct {
	repo repository.ProductRepository
	inv  InventoryService
}

func NewProductService(ctx handlers.Context) ProductService {
	return &ProductServiceImpl{
		repo: repository.NewProductRepository(ctx.DB),
		inv:  NewInvetroryService(ctx.Inventory, ctx.Cache, ctx.Logger),
	}
}

func (s *ProductServiceImpl) SelectById(id int) (model.Product, error) {
	product, err := s.repo.SelectById(id)
	if err != nil {
		return product, err
	}

	product.Avaliable = s.inv.Check(id)
	return product, nil
}

func (s *ProductServiceImpl) Select(offset, limit int) (model.SelectResponse, error) {
	data, count, err := s.repo.Select(offset, limit)
	if err != nil {
		return model.SelectResponse{}, err
	}

	ids := make([]int, len(data))
	for i := range data {
		ids[i] = data[i].Id
	}
	avaliable := s.inv.CheckRange(ids)
	for i := range data {
		data[i].Avaliable = avaliable[i]
	}
	return model.SelectResponse{Products: data, Count: count}, err
}

func (s *ProductServiceImpl) Create(value model.Product) (model.Product, error) {
	id, err := s.repo.Create(value)
	value.Id = id
	return value, err
}

func (s *ProductServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductServiceImpl) Update(value model.Product) error {
	return s.repo.Update(value)
}
