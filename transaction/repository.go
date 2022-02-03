package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
	GetByID(ID int) (Transaction, error) // untuk mengambil data transaksi berdasarkan transaction id
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
}

func RepositoryBaru(db *gorm.DB) *repository {
	return &repository{db}
}


// get campaign transactions endpoint
func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction
    
	// Fungsi 'Order("id desc")' adalah untuk mengurutkan data dari id yang paling besar terlebih dahulu
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}


// get user transactions endpoint
func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction
    
	// Fungsi dari 'Preload("Campaign.CampaignImages")' agar transaction dapat memiliki relasi dengan campaign_images melalui campaign sembari mengload campaign
	// Fungsi dari "campaign_images.is_primary = 1" untuk membatasi agar data yg dipanggil hanya gambar yg is_primarynya 1
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}


// get transaction's data from transaction id (nitification payment midtrans)
func (r *repository) GetByID(ID int) (Transaction, error) {
	var transaction Transaction
    
	err := r.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}


// user create transaction endpoin
func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
