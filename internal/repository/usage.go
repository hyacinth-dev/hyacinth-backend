package repository

import (
	"context"
	"fmt"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"time"
)

type UsageRepository interface {
	GetUsage(ctx context.Context, userId string, vnetId string, timeRange string) (*[]v1.UsageData, error)
}

func NewUsageRepository(
	repository *Repository,
) UsageRepository {
	return &usageRepository{
		Repository: repository,
	}
}

type usageRepository struct {
	*Repository
}

func (r *usageRepository) GetUsage(ctx context.Context, userId string, vnetId string, timeRange string) (*[]v1.UsageData, error) {
	var result0 []v1.UsageData

	db := r.DB(ctx).Where("deleted_at IS NULL")

	fmt.Println("Fetching usage data with parameters:", userId, vnetId, timeRange)

	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	// 添加虚拟网络ID筛选条件
	if vnetId != "" {
		db = db.Where("vnet_id = ?", vnetId)
	}

	switch timeRange {
	case "24h":
		db = db.Where("created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY)")
		// 按小时分组
		db = db.Model(&model.Usage{}).
			Select("DATE_FORMAT(created_at, '%m-%d %H:00') as `date`, SUM(`usage`) as `usage`").
			Group("DATE_FORMAT(created_at, '%m-%d %H:00')").
			Order("`date` ASC")
	case "7d":
		println("Fetching usage for the last 7 days")
		db = db.Where("created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)")
		// 按天分组
		db = db.Model(&model.Usage{}).
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as `date`, SUM(`usage`) as `usage`").
			Group("DATE_FORMAT(created_at, '%Y-%m-%d')").
			Order("created_at ASC")
	case "30d":
		db = db.Where("created_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)")
		// 按天分组
		db = db.Model(&model.Usage{}).
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as `date`, SUM(`usage`) as `usage`").
			Group("DATE_FORMAT(created_at, '%Y-%m-%d')").
			Order("created_at ASC")
	case "month":
		db = db.Where("created_at >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)")
		// 按月分组
		db = db.Model(&model.Usage{}).
			Select("DATE_FORMAT(created_at, '%Y-%m') as `date`, SUM(`usage`) as `usage`").
			Group("DATE_FORMAT(created_at, '%Y-%m')").
			Order("created_at ASC")
	case "all":
		db = db.Where("created_at >= DATE_SUB(CURDATE(), INTERVAL 100 YEAR)")
		// 按年分组
		db = db.Model(&model.Usage{}).
			Select("DATE_FORMAT(created_at, '%Y') as `date`, SUM(`usage`) as `usage`").
			Group("DATE_FORMAT(created_at, '%Y')").
			Order("created_at ASC")
	default:
		return nil, nil
	}

	if err := db.Scan(&result0).Error; err != nil {
		return nil, err
	}

	results := []v1.UsageData{}

	if timeRange != "all" {
		today := time.Now()
		switch timeRange {
		case "24h":
			today = today.Add(-time.Hour * 24)
			j := 0
			for i := 0; i < 24; i++ {
				today = today.Add(time.Hour)
				if j < len(result0) && result0[j].Date == today.Format("01-02 15:00") {
					results = append(results, result0[j])
					j++
				} else {
					results = append(results, v1.UsageData{
						Date:  today.Format("01-02 15:00"),
						Usage: 0,
					})
				}
			}
		case "7d":
			today = today.AddDate(0, 0, -7)
			j := 0
			for i := 0; i < 7; i++ {
				today = today.AddDate(0, 0, 1)
				if j < len(result0) && result0[j].Date == today.Format("2006-01-02") {
					results = append(results, result0[j])
					j++
				} else {
					results = append(results, v1.UsageData{
						Date:  today.Format("2006-01-02"),
						Usage: 0,
					})
				}
			}
		case "30d":
			today = today.AddDate(0, 0, -30)
			j := 0
			for i := 0; i < 30; i++ {
				today = today.AddDate(0, 0, 1)
				if j < len(result0) && result0[j].Date == today.Format("2006-01-02") {
					results = append(results, result0[j])
					j++
				} else {
					results = append(results, v1.UsageData{
						Date:  today.Format("2006-01-02"),
						Usage: 0,
					})
				}
			}
		case "month":
			today = today.AddDate(0, -13, 0)
			j := 0
			for i := 0; i < 13; i++ {
				today = today.AddDate(0, 1, 0)
				if j < len(result0) && result0[j].Date == today.Format("2006-01") {
					results = append(results, result0[j])
					j++
				} else {
					results = append(results, v1.UsageData{
						Date:  today.Format("2006-01"),
						Usage: 0,
					})
				}
			}
			if results[len(results)-1].Date > time.Now().Format("2006-01") {
				results = results[:len(results)-1]
			}
		}
	}

	return &results, nil
}
