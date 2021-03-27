package lightning_interface

type Lightinng interface {
	Invoice(msat int64, desc string) (checkingId string)
	Pay(bolt11 string) (checkingId string)

	ListenInvoices()
	ListenPayments()
}
