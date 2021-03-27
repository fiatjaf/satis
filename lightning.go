package main

func makeInvoice(account string, msat int64, desc string) (bolt11 string, err error) {
	checkingId := lightning.Invoice(msat, desc)
	store.SavePendingInvoice(account, checkingId)

	return "", nil
}

func payInvoice(account string, bolt11 string) {
	checkingId := lightning.Pay(bolt11)
	store.SavePendingPayment(account, checkingId)
}
