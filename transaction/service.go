package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"log"

	"strconv"

	"errors"
	"fmt"
)

type service struct {
	repo               Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(UserID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(trans TransactionNotificationInput) error
	GetAllTransaction() ([]Transaction, error)
}

func NewServiceTransaction(r Repository, campaignRepository campaign.Repository, payservice payment.Service) *service {
	return &service{r, campaignRepository, payservice}
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	trans_id, _ := strconv.Atoi(input.OrderID)
	fmt.Println(trans_id)
	transaction, err := s.repo.GetByTransactionId(trans_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedtransaction, err := s.repo.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindById(updatedtransaction.CampaignID)
	if err != nil {
		return err
	}

	if transaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedtransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	var trans Transaction
	log.Println("data", input)
	trans.Amount = input.Amount
	trans.User = input.User
	trans.CampaignID = input.CampaignID
	trans.Status = "pending"

	trans.Code = "muamama"

	fmt.Println(trans, "madang trans")
	newtrans, err := s.repo.SaveTransaction(trans)
	if err != nil {
		return newtrans, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newtrans.ID,
		Amount: newtrans.Amount,
	}

	url, err := s.paymentService.GetPaymentUrl(paymentTransaction, input.User)
	if err != nil {
		fmt.Println("error url")
		return newtrans, err
	}
	newtrans.PaymentURL = url
	newtranss, err := s.repo.Update(newtrans)
	if err != nil {
		return newtranss, err
	}
	return newtranss, nil
}

func (s *service) GetTransactionByUserID(UserID int) ([]Transaction, error) {
	var transactions []Transaction
	fmt.Println("gettransbyuserid", UserID)
	transactions, err := s.repo.GetByUserId(UserID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	var transactions []Transaction
	campaign, err := s.campaignRepository.FindById(input.ID)
	if err != nil {
		return transactions, err
	}

	if campaign.UserID != input.User.ID {
		fmt.Println("error not owner check")
		fmt.Println("error not owner check", campaign.UserID, input.User.ID)
		return transactions, errors.New("not an owner of campaign")
	}

	transactions, err = s.repo.GetByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) GetAllTransaction() ([]Transaction, error) {
	var transactions []Transaction
	transactions, err := s.repo.FindAll()
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
