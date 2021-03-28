package lightning_interface

type InvoiceResult struct {
	CheckingId string
	Amount     int64
	Result     bool
}

type PaymentResult struct {
	CheckingId string
	Result     bool
}

type Lightning interface {
	Invoice(msat int64, desc string) (checkingId string)
	CheckInvoice(checkingId string) (result *InvoiceResult)
	ListenInvoices() chan InvoiceResult

	Pay(bolt11 string) (checkingId string)
	CheckPayment(checkingId string) (result *PaymentResult)
	ListenPayments() chan InvoiceResult
}
