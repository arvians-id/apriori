package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/model"
	"github.com/go-redis/redis/v8"
	"strings"
)

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) ProductFindAllByAdmin(ctx context.Context) ([]*model.Product, error) {
	products, err := r.ProductService.FindAllByAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Product
	for _, product := range products {
		result = append(result, &model.Product{
			IdProduct:   product.IdProduct,
			Code:        product.Code,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Category:    product.Category,
			IsEmpty:     product.IsEmpty,
			Mass:        product.Mass,
			Image:       product.Image,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) ProductFindAllSimilarCategory(ctx context.Context, code string) ([]*model.Product, error) {
	products, err := r.ProductService.FindAllBySimilarCategory(ctx, code)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}
		return nil, err
	}

	return products, nil
}

func (r *queryResolver) ProductFindAllByUser(ctx context.Context, search *string, category *string) ([]*model.Product, error) {
	if search == nil {
		search = new(string)
	}

	if category == nil {
		category = new(string)
	}

	searchQuery := strings.ToLower(*search)
	categoryQuery := strings.ToLower(*category)
	products, err := r.ProductService.FindAll(ctx, searchQuery, categoryQuery)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *queryResolver) ProductFindAllRecommendation(ctx context.Context, code string) ([]*model.ProductRecommendation, error) {
	products, err := r.ProductService.FindAllRecommendation(ctx, code)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}
		return nil, err
	}

	return products, nil
}

func (r *queryResolver) ProductFindByCode(ctx context.Context, code string) (*model.Product, error) {
	key := fmt.Sprintf("product-%s", code)
	productCache, err := r.CacheService.Get(ctx, key)
	if err == redis.Nil {
		product, err := r.ProductService.FindByCode(ctx, code)
		if err != nil {
			if err.Error() == response.ErrorNotFound {
				return nil, errors.New(response.ResponseErrorNotFound)
			}
			return nil, err
		}

		err = r.CacheService.Set(ctx, key, product)
		if err != nil {
			return nil, err
		}

		return product, nil
	} else if err != nil {
		return nil, err
	}

	var productCacheResponse model.Product
	err = json.Unmarshal(productCache, &productCacheResponse)
	if err != nil {
		return nil, err
	}

	return &productCacheResponse, nil
}

func (r *queryResolver) AuthToken(ctx context.Context) (bool, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	_, isExist := ginContext.Get("id_user")
	if !isExist {
		return false, errors.New("unauthorized")
	}

	return true, nil
}

func (r *queryResolver) UserFindAll(ctx context.Context) ([]*model.User, error) {
	users, err := r.UserService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *queryResolver) UserProfile(ctx context.Context) (*model.User, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	id, isExist := ginContext.Get("id_user")
	if !isExist {
		return nil, errors.New("unauthorized")
	}

	user, err := r.UserService.FindById(ctx, int(id.(float64)))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) UserFindByID(ctx context.Context, id int) (*model.User, error) {
	user, err := r.UserService.FindById(ctx, id)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return user, nil
}

func (r *queryResolver) CategoryFindAll(ctx context.Context) ([]*model.Category, error) {
	categoriesCache, err := r.CacheService.Get(ctx, "categories")
	if err == redis.Nil {
		categories, err := r.CategoryService.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		err = r.CacheService.Set(ctx, "categories", categories)
		if err != nil {
			return nil, err
		}

		return categories, nil
	} else if err != nil {
		return nil, err
	}

	var categoryCacheResponses []*model.Category
	err = json.Unmarshal(categoriesCache, &categoryCacheResponses)
	if err != nil {
		return nil, err
	}

	return categoryCacheResponses, nil
}

func (r *queryResolver) CategoryFindByID(ctx context.Context, id int) (*model.Category, error) {
	category, err := r.CategoryService.FindById(ctx, id)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return category, nil
}

func (r *queryResolver) TransactionFindAll(ctx context.Context) ([]*model.Transaction, error) {
	transactionCache, err := r.CacheService.Get(ctx, "all-transaction")
	if err == redis.Nil {
		transaction, err := r.TransactionService.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		err = r.CacheService.Set(ctx, "all-transaction", transaction)
		if err != nil {
			return nil, err
		}

		return transaction, nil
	} else if err != nil {
		return nil, err
	}

	var transactionCacheResponses []*model.Transaction
	err = json.Unmarshal(transactionCache, &transactionCacheResponses)
	if err != nil {
		return nil, err
	}

	return transactionCacheResponses, nil
}

func (r *queryResolver) TransactionFindByNoTransaction(ctx context.Context, numberTransaction string) (*model.Transaction, error) {
	transactions, err := r.TransactionService.FindByNoTransaction(ctx, numberTransaction)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return transactions, nil
}

func (r *queryResolver) PaymentFindAll(ctx context.Context) ([]*model.Payment, error) {
	payments, err := r.PaymentService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *queryResolver) PaymentFindByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	payment, err := r.PaymentService.FindByOrderId(ctx, orderID)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return payment, nil
}

func (r *queryResolver) CommentFindAllRatingByProductCode(ctx context.Context, productCode string) ([]*model.RatingFromComment, error) {
	comments, err := r.CommentService.FindAllRatingByProductCode(ctx, productCode)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return comments, nil
}

func (r *queryResolver) CommentFindAllByProductCode(ctx context.Context, productID string, tags string, ratings string) ([]*model.Comment, error) {
	comments, err := r.CommentService.FindAllByProductCode(ctx, productID, ratings, tags)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return comments, nil
}

func (r *queryResolver) CommentFindByUserOrderID(ctx context.Context, userOrderID int) (*model.Comment, error) {
	comment, err := r.CommentService.FindByUserOrderId(ctx, userOrderID)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return comment, nil
}

func (r *queryResolver) CommentFindByID(ctx context.Context, id int) (*model.Comment, error) {
	comment, err := r.CommentService.FindById(ctx, id)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return comment, nil
}

func (r *queryResolver) UserOrderFindAll(ctx context.Context) ([]*model.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) UserOrderFindAllByUserID(ctx context.Context) ([]*model.UserOrder, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) UserOrderFindAllByID(ctx context.Context, orderID string) ([]*model.UserOrder, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) UserOrderFindByID(ctx context.Context, id int) (*model.UserOrder, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) RajaOngkirFindAll(ctx context.Context, place string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) NotificationFindAll(ctx context.Context) ([]*model.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) NotificationFindAllByUserID(ctx context.Context) ([]*model.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) AprioriFindAll(ctx context.Context) ([]*model.Apriori, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) AprioriFindAllByCode(ctx context.Context, code string) ([]*model.Apriori, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) AprioriFindAllByActive(ctx context.Context) ([]*model.Apriori, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) AprioriFindByCodeAndID(ctx context.Context, code string, id int) (*model.ProductRecommendation, error) {
	//TODO implement me
	panic("implement me")
}
