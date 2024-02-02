package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	FindCampaigns(userId int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputparam CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, filelocation string) (CampaignImage, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, filelocation string) (CampaignImage, error) {
	campaign, err := s.repo.FindById(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("error not owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		_, err := s.repo.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
		isPrimary = 1
	}
	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID

	campaignImage.IsPrimary = isPrimary

	campaignImage.FileName = filelocation

	newCAmpaignImage, err := s.repo.SaveImage(campaignImage)

	if err != nil {
		return newCAmpaignImage, err
	}

	return newCAmpaignImage, nil

}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repo.FindById(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != input.User.ID {
		return campaign, errors.New("error not owner of the campaign")
	}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount

	updated, err := s.repo.Update(campaign)
	if err != nil {
		return updated, err
	}

	return updated, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	preslug := fmt.Sprintf("%s-%d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(preslug)
	newcampaign, err := s.repo.Save(campaign)
	if err != nil {
		return newcampaign, err
	}
	return newcampaign, nil
}

func (s *service) FindCampaigns(userId int) ([]Campaign, error) {
	if userId != 0 {
		campains, err := s.repo.FindByUserId(userId)
		if err != nil {
			return campains, err
		}
		return campains, nil
	}
	campaigs, err := s.repo.FindAll()
	if err != nil {
		return campaigs, err
	}
	return campaigs, nil
}

func (s *service) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
	var campaign Campaign
	campaign, err := s.repo.FindById(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
