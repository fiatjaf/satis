package filedb

import (
	"fmt"
	"strconv"

	"github.com/tidwall/buntdb"
)

type FileDatabase struct {
	Path string

	db *buntdb.DB
}

func (fd *FileDatabase) Init() error {
	db, err := buntdb.Open(fd.Path)
	if err != nil {
		return err
	}
	fd.db = db
	return nil
}

func (fd *FileDatabase) SetBalance(account string, msat int64) error {
	return fd.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("b:"+account, strconv.FormatInt(msat, 10), nil)
		return err
	})
}

func (fd *FileDatabase) GetBalances() (balances map[string]int64, err error) {
	balances = make(map[string]int64)

	err = fd.db.View(func(tx *buntdb.Tx) error {
		var err error

		tx.AscendRange("", "b:", "c", func(key, value string) bool {
			account := key[2:]
			msatoshi, errW := strconv.ParseInt(value, 10, 64)
			if errW != nil {
				// exit with this error
				err = fmt.Errorf(
					"there's an invalid balance in the database: %s: %s, %w",
					key, value, errW,
				)
				return false
			} else {
				balances[account] = msatoshi
			}
			return true
		})

		return err
	})

	return balances, err
}

func (fd *FileDatabase) SavePendingPayment(account string, checkingId string) error {
	return fd.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("p:"+account+":"+checkingId, "", nil)
		return err
	})
}

func (fd *FileDatabase) SavePendingInvoice(account string, checkingId string) error {
	return fd.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("i:"+checkingId, account, nil)
		return err
	})
}

func (fd *FileDatabase) GetPendingLightning() (
	payments map[string][]string, // {[account]: [checkingId, ...]}
	invoices map[string][]string, // {[account]: [checkingId, ...]}
	err error,
) {
	payments = make(map[string][]string)
	invoices = make(map[string][]string)

	err = fd.db.View(func(tx *buntdb.Tx) error {
		var err error

		tx.AscendRange("", "p:", "q", func(key, value string) bool {
			checkingId := key[2:]
			account := value

			pays, ok := payments[account]
			if !ok {
				pays = make([]string, 0, 1)
			}
			pays = append(pays, checkingId)

			payments[account] = pays
			return true
		})

		tx.AscendRange("", "i:", "j", func(key, value string) bool {
			checkingId := key[2:]
			account := value

			invs, ok := invoices[account]
			if !ok {
				invs = make([]string, 0, 1)
			}
			invs = append(invs, checkingId)

			invoices[account] = invs
			return true
		})

		return err
	})

	return payments, invoices, err
}
