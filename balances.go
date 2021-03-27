package main

import "sync"

type balances map[string]int64

var (
	mu sync.RWMutex
	b  = make(balances)
)

func (b balances) get(account []byte) int64 {
	mu.Lock()
	msat, _ := b[string(account)]
	mu.Unlock()
	return msat
}

func (b balances) set(account []byte, msat int64) {
	mu.Lock()
	b[string(account)] = msat
	mu.Unlock()
}

func (b balances) del(account []byte) {
	mu.Lock()
	delete(b, string(account))
	mu.Unlock()
}

func (b balances) persist(accountsAffected ...[]byte) {
	for _, acct := range accountsAffected {
		store.SetBalance(string(acct), b.get(acct))
	}
}
