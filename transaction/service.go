package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
	"strconv"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service     
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
    ProccesPayment(input TransactionNotificationInput) error
}

func ServiceBaru(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}


// get campaign transactions endpoint
func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
    
	// dapatkan data campaign terlebih dahulu melalui "campaignRepository campaign.Repository"
    campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	
	// Authorization = berfungsi agar suatu user tidak bisa melihat data campaign's transaction milik orang lain
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}


// get user transactions endpoint
func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}


// user create transaction endpoint
func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = ""

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}
    
	// memanggil sistem midtrans (payment service)
	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

// notification payment midtrans
func (s *service) ProccesPayment(input TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID) // fungsi "strconv.Atoi()" adalah untuk mengubah data string menjadi int

	transaction, err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}
    
	// if use payment type credit card if capture & accept
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" { // if use other payment type (settlement not capture & accept)
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	// cara update data campaign setelah ada transaksi yang masuk pada campaign tsb
	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}
    
	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil

}