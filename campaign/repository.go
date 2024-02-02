package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserId(userId int) ([]Campaign, error)
	FindById(id int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	SaveImage(campaignImages CampaignImage) (CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignid int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) MarkAllImagesAsNonPrimary(campaignid int) (bool, error) {
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignid).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error

	if err != nil {
		return campaign, err
	}
	return campaign, nil

}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	// err := r.db.Preload("CampaignImages").Find(&campaigns).Error
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByUserId(userId int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userId).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil

}

func (r *repository) FindById(id int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Where("id = ?", id).Preload("CampaignImages").Preload("User").Find(&campaign).Error

	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) SaveImage(campaignImages CampaignImage) (CampaignImage, error) {

	err := r.db.Create(&campaignImages).Error
	if err != nil {
		return campaignImages, err
	}
	return campaignImages, nil
}
