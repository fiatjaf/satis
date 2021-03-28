package sparko

import lightning_interface "github.com/fiatjaf/satis/lightning"

type Client struct {
	URL   string
	Token string
}

func (sp *Client) Invoice(msat int64, desc string) (checkingId string) {

}

func (sp *Client) CheckInvoice(checkingId string) (result *lightning_interface.InvoiceResult) {

}

func (sp *Client) ListenInvoices() chan lightning_interface.InvoiceResult {

}

func (sp *Client) Pay(bolt11 string) (checkingId string) {

}

func (sp *Client) CheckPayment(checkingId string) (result *lightning_interface.PaymentResult) {

}

func (sp *Client) ListenPayments() chan lightning_interface.InvoiceResult {

}
