package campaign

import (
	"bwastartup/user"
	"time"

	"github.com/leekchan/accounting"
)

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             user.User
}

// func (cam *Campaign) GoalAmountFormatIDR() string {

// 	return "formater.FormatMoney(cam.GoalAmount)"
// }

func (cam *Campaign) GoalAmountFormatIDR() string {
	formater := accounting.Accounting{
		Symbol:         "Rp",
		Precision:      2,
		Thousand:       ".",
		Decimal:        ",",
		Format:         "",
		FormatNegative: "",
		FormatZero:     "",
	}
	return formater.FormatMoney(cam.GoalAmount)
}
func (cam *Campaign) CurrentAmountFormatIDR() string {
	formater := accounting.Accounting{
		Symbol:         "Rp",
		Precision:      2,
		Thousand:       ".",
		Decimal:        ",",
		Format:         "",
		FormatNegative: "",
		FormatZero:     "",
	}
	return formater.FormatMoney(cam.CurrentAmount)
}
