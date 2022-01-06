package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
}

func RepositoryBaru(db *gorm.DB) *repository {
	return &repository{db}
}

// campaign transactions endpoint
func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction
    
	// Fungsi "Order("id desc")" adalah untuk mengurutkan data dari id yang paling besar terlebih dahulu
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}