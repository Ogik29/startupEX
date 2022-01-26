package payment

import (
	"bwastartup/transaction"
	"bwastartup/user"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction transaction.Transaction, user user.User) (string, error)
}

func ServiceBaru() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction transaction.Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
    midclient.ServerKey = "SB-Mid-server-99OWVKV-wpWhm1dTH_vk3qgw"
    midclient.ClientKey = "SB-Mid-client-_1A8PXa42FUNXacp"
    midclient.APIEnvType = midtrans.Sandbox

    snapGateway := midtrans.SnapGateway {
        Client: midclient,
    }

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},

		TransactionDetails: midtrans.TransactionDetails{
			OrderID: strconv.Itoa(transaction.ID), //fungsi "strconv.Itoa()" adalah untuk mengubah data int menjadi string
			GrossAmt: int64(transaction.Amount), // fungsi "int64()" adalah untuk mengubah data int biasa menjadi int64
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}
	return snapTokenResp.RedirectURL, nil
}