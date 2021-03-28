package main

import (
	"log"
	"strconv"
	"strings"

	decodepay "github.com/fiatjaf/ln-decodepay"
	storage_interface "github.com/fiatjaf/satis/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/tidwall/redcon"
)

type Settings struct {
	Port string `envconfig:"PORT" default: "9763"`

	FileDBPath string `envconfig:"FILEDB_PATH" default:"satis.filedb"`
}

var (
	s         Settings
	store     storage_interface.Storage
	lightning lightning_interface.Lightning

	ps redcon.PubSub
)

func main() {
	// read env
	if err := envconfig.Process("", &s); err != nil {
		panic(err)
	}

	initializeStorage()
	initializeLightning()

	go log.Printf("started server at %s", s.Port)

	err := redcon.ListenAndServe(s.Port,
		func(conn redcon.Conn, cmd redcon.Command) {
			command := strings.ToLower(string(cmd.Args[0]))

			if minParams, ok := minParamCount[command]; !ok {
				conn.WriteError("ERR unknown command '" + command + "'")
				return
			} else if (minParams + 1) != len(cmd.Args) {
				conn.WriteError("ERR wrong number of arguments for '" + command + "' command, needed " + strconv.Itoa(minParams))
				return
			}

			switch command {
			case "ping":
				conn.WriteString("PONG")
				return
			case "quit":
				conn.WriteString("OK")
				conn.Close()
				return
			case "set":
				if msat, err := getMsat(cmd.Args[2]); err != nil {
					conn.WriteError("ERR value is not an integer or out of range")
					return
				} else {
					b.set(cmd.Args[1], msat)
					conn.WriteString("OK")

					go b.persist(cmd.Args[1])
				}
			case "get":
				val := b.get(cmd.Args[1])
				conn.WriteBulk([]byte(strconv.FormatInt(val, 10)))
			case "del":
				b.del(cmd.Args[1])
				conn.WriteString("OK")
				go b.persist(cmd.Args[1])
			case "incr":
				if msat, err := getMsat(cmd.Args[2]); err != nil {
					conn.WriteError("ERR value is not an integer or out of range")
					return
				} else {
					curr := b.get(cmd.Args[1])
					b.set(cmd.Args[1], curr+msat)
					conn.WriteString("OK")
					go b.persist(cmd.Args[1])
				}
			case "decr":
				if msat, err := getMsat(cmd.Args[2]); err != nil {
					conn.WriteError("ERR value is not an integer or out of range")
					return
				} else {
					curr := b.get(cmd.Args[1])
					next := curr - msat
					if next < 0 {
						conn.WriteError("ERR balance would go below 0")
						return
					}

					b.set(cmd.Args[1], next)
					conn.WriteString("OK")
					go b.persist(cmd.Args[1])
				}
			case "transfer":
				if msat, err := getMsat(cmd.Args[2]); err != nil {
					conn.WriteError("ERR value is not an integer or out of range")
					return
				} else {
					currA := b.get(cmd.Args[1])
					nextA := currA - msat
					if nextA < 0 {
						conn.WriteError("ERR balance would go below 0")
						return
					}

					currB := b.get(cmd.Args[3])
					nextB := currB + msat

					b.set(cmd.Args[1], nextA)
					b.set(cmd.Args[3], nextB)

					conn.WriteString("OK")
					go b.persist(cmd.Args[1], cmd.Args[3])
				}
			case "pay":
				account := string(cmd.Args[1])
				bolt11 := string(cmd.Args[2])
				inv, err := decodepay.Decodepay(bolt11)
				if err != nil {
					conn.WriteError("ERR invoice is invalid")
					return
				}

				next := b.get(cmd.Args[1]) - inv.MSatoshi
				if next < 0 {
					conn.WriteError("ERR balance would go below 0")
					return
				}
				b.set(cmd.Args[1], next)

				go func() {
					checkingId := lightning.Pay(bolt11)
					store.SavePendingPayment(account, checkingId, msat)
				}()

				conn.WriteString("OK")
				go b.persist(cmd.Args[1])
			case "invoice":
				if msat, err := getMsat(cmd.Args[2]); err != nil {
					conn.WriteError("ERR value is not an integer or out of range")
					return
				} else {
					account := string(cmd.Args[1])
					desc := string(cmd.Args[3])
					bolt11, checkingId, err := lightning.Invoice(msat, desc)

					if err != nil {
						conn.WriteError("ERR failed: '" + err.Error() + "'")
					} else {
						store.SavePendingInvoice(account, checkingId)
						conn.WriteString(bolt11)
					}
				}
			}
		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
