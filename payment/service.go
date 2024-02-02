package payment

import (
	"bwastartup/user"
	"fmt"
	"log"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
}

func NewPaymentService() *service {
	return &service{}
}

func (s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-rfoEAZy_b7QibK3ZZTM4-Jsp"
	midclient.ClientKey = "SB-Mid-client-S1AvvSdTk2QBf7D9"
	midclient.APIEnvType = midtrans.Sandbox

	// var snapGateway midtrans.SnapGateway
	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
	}
	log.Println("GetToken:")
	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		fmt.Println(snapTokenResp, "madang", err)
		return "error getting token", err
	}
	return snapTokenResp.RedirectURL, nil
}
