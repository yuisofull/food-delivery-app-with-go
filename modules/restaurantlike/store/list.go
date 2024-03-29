package restaurantlikestorage

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/yuisofull/food-delivery-app-with-go/common"
	restaurantlikemodel "github.com/yuisofull/food-delivery-app-with-go/modules/restaurantlike/model"
)

var timeLayout = "2006-01-02T15:04:05.99999"

func (s *sqlStore) GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)

	type sqlData struct {
		RestaurantId int `gorm:"column:restaurant_id"`
		LikedCount   int `gorm:"column:count"`
	}

	var listLike []sqlData

	if err := s.db.Table(restaurantlikemodel.Like{}.TableName()).Select("restaurant_id, count(restaurant_id) as count").
		Where("restaurant_id in (?)", ids).
		Group("restaurant_id").
		Find(&listLike).
		Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for _, item := range listLike {
		result[item.RestaurantId] = item.LikedCount
	}

	return result, nil
}

func (s *sqlStore) GetUsersLikeRestaurant(ctx context.Context,
	condition map[string]interface{},
	filter *restaurantlikemodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]common.SimpleUser, error) {
	var result []restaurantlikemodel.Like
	db := s.db.Table(restaurantlikemodel.Like{}.TableName()).Where(condition)

	if v := filter; v != nil {
		if v.RestaurantId > 0 {
			db = db.Where("restaurant_id = ?", v.RestaurantId)
		}
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	//for i := range moreKeys {
	//	db = db.Preload(moreKeys[i])
	//}
	db = db.Preload("User")

	if v := paging.FakeCursor; v != "" {
		db = db.Where("created_at < ?", base58.Decode(v))
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db = db.Offset(int(offset))

	}

	if err := db.
		Limit(int(paging.Limit)).
		Order("created_at desc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(result))

	for i := range result {
		//result[i].User = &common.SimpleUser{}
		result[i].User.CreatedAt = result[i].CreatedAt
		result[i].User.UpdatedAt = nil
		users[i] = *result[i].User // when len of two arr not equal (mean: user may is null) => system crashing
	}
	if len(result) > 0 {
		cursorStr := base58.Encode([]byte(fmt.Sprintf("%v", result[len(result)-1].CreatedAt.Format(timeLayout))))
		paging.NextCursor = cursorStr
	}

	return users, nil
}
