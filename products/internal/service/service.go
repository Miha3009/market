package service

import (
	"products/internal/model"
	"products/internal/repository"
	"products/pkg/handlers"
)

type ProductService interface {
	SelectById(id int) (model.Product, error)
	Select(offset, limit int) (model.SelectResponse, error)
	Create(value model.Product) (model.CreateResponse, error)
	Delete(id int) error
	Update(value model.Product) error
}

type ProductServiceImpl struct {
	repo repository.ProductRepository
}

func NewProductService(ctx handlers.Context) ProductService {
	return &ProductServiceImpl{
		repo: repository.NewProductRepository(ctx.DB),
	}
}

func (s *ProductServiceImpl) SelectById(id int) (model.Product, error) {
	return s.repo.SelectById(id)
}

func (s *ProductServiceImpl) Select(offset, limit int) (model.SelectResponse, error) {
	data, count, err := s.repo.Select(offset, limit)
	return model.SelectResponse{Products: data, Count: count}, err
}

func (s *ProductServiceImpl) Create(value model.Product) (model.CreateResponse, error) {
	id, err := s.repo.Create(value)
	return model.CreateResponse{Id: id}, err
}

func (s *ProductServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductServiceImpl) Update(value model.Product) error {
	return s.repo.Update(value)
}
