package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserId(userId int) ([]Transaction, error)
	GetByTransactionId(Id int) (Transaction, error)
	SaveTransaction(trans Transaction) (Transaction, error)
	Update(trans Transaction) (Transaction, error)
	FindAll() ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Update(trans Transaction) (Transaction, error) {

	err := r.db.Save(&trans).Error

	if err != nil {
		return trans, err
	}
	return trans, nil

}

func (r *repository) SaveTransaction(trans Transaction) (Transaction, error) {

	err := r.db.Create(&trans).Error
	// err := r.db.Omit("PaymentURL").Create(&trans).Error
	if err != nil {
		return trans, err
	}
	return trans, nil
}

func (r *repository) GetByTransactionId(Id int) (Transaction, error) {
	var transaction Transaction

	err := r.db.Where("id = ?", Id).Order("id desc").Find(&transaction).Error
	// err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		fmt.Println(err, "error repo get bu transid")
		return transaction, err
	}
	return transaction, nil
}

func (r *repository) GetByUserId(userId int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userId).Order("id desc").Find(&transactions).Error
	// err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		fmt.Println(err, "error repo get bu userid")
		return transactions, err
	}
	return transactions, nil
}

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}
	return transactions, nil

}

func (r *repository) FindAll() ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("User").Preload("Campaign").Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
