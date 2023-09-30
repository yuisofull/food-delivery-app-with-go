package restaurantbusiness

import (
	"context"
	"food-delivery/common"
	restaurantmodel "food-delivery/module/restaurant/model"
)

type ListRestaurantStore interface {
	ListRestaurantWithCondition(
		context context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]restaurantmodel.Restaurant, error)
}

type listRestaurantBusiness struct {
	store ListRestaurantStore
}

func NewListRestaurantBusiness(store ListRestaurantStore) *listRestaurantBusiness {
	return &listRestaurantBusiness{store: store}
}

func (business *listRestaurantBusiness) ListRestaurant(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
) ([]restaurantmodel.Restaurant, error) {
	res, err := business.store.ListRestaurantWithCondition(context, filter, paging)
	if err != nil {
		return nil, err
	}
	return res, nil
}