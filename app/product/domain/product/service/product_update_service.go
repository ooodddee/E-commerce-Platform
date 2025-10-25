package service

import (
	"context"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/product/common/model/entity"
	productrepo "github.com/Vigor-Team/youthcamp-2025-mall-be/app/product/domain/product/repository"
)

type ProductUpdateService struct{}

var productUpdateService ProductUpdateService

func GetProductUpdateService() *ProductUpdateService {
	return &productUpdateService
}

func (s *ProductUpdateService) AddProduct(ctx context.Context, product *entity.ProductEntity) error {
	err := productrepo.GetFactory().GetProductRepository().AddProduct(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductUpdateService) UpdateProduct(ctx context.Context, origin, target *entity.ProductEntity) error {
	err := productrepo.GetFactory().GetProductRepository().UpdateProduct(ctx, origin, target)
	if err != nil {
		return err
	}
	return nil
}
