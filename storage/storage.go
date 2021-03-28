package storage_interface

type Storage interface {
	Init() error

	SetBalance(account string, msat int64) error
	GetBalances() (balances map[string]int64, err error)

	SavePendingPayment(account string, checkingId string, msat int64) error
	SavePendingInvoice(account string, checkingId string) error
	GetPendingLightning(
		payments map[string][]string,
		invoices map[string][]string,
		err error,
	)
}
