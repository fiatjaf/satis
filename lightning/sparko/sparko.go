package sparko

import (
	"errors"

	lg "github.com/fiatjaf/lightningd-gjson-rpc"
	decodepay "github.com/fiatjaf/ln-decodepay"
	lightning_interface "github.com/fiatjaf/satis/lightning"
	"github.com/lucsky/cuid"
)

type Client struct {
	lg.Client
}

func (sp *Client) Invoice(msat int64, desc string) (bolt11, checkingId string, err error) {
	checkingId = "sis:" + cuid.Slug()

	resp, err := sp.Call("invoice", msat, checkingId, desc)
	if err != nil {
		return "", "", err
	}

	bolt11 = resp.Get("bolt11").String()
	if bolt11 == "" {
		return "", "", errors.New("bolt11 is blank")
	}

	return checkingId, bolt11, nil
}

func (sp *Client) CheckInvoice(checkingId string) (result *lightning_interface.InvoiceResult) {
	resp, err := sp.Call("listinvoices", checkingId)
	if err != nil {
		return nil
	}

	result = &lightning_interface.InvoiceResult{
		CheckingId: checkingId,
		Amount:     resp.Get("invoices.0.msatoshi").Int(),
	}

	if resp.Get("invoices.0.status").String() == "paid" {
		result.Result = true
	}

	return result
}

func (sp *Client) ListenInvoices() chan lightning_interface.InvoiceResult {
	c := make(chan lightning_interface.InvoiceResult)
	return c
}

func (sp *Client) Pay(bolt11 string) (checkingId string) {
	inv, _ := decodepay.Decodepay(bolt11)
	checkingId = inv.PaymentHash
	go sp.Call("pay", bolt11)
	return checkingId
}

func (sp *Client) CheckPayment(checkingId string) (result *lightning_interface.PaymentResult) {
	result = &lightning_interface.PaymentResult{
		CheckingId: checkingId,
	}

	resp, err := sp.Call("listpays", checkingId)
	if err != nil {
		result.Result = false
	} else {
		switch resp.Get("pays.0.status").String() {
		case "pending":
			return nil
		case "failed":
			result.Result = false
		case "complete":
			result.Result = true
		}
	}

	return result
}

func (sp *Client) ListenPayments() chan lightning_interface.PaymentResult {
	c := make(chan lightning_interface.PaymentResult)
	return c
}
