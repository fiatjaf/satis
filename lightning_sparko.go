package main

import (
	"github.com/fiatjaf/satis/lightning/sparko"
)

func initializeLightning() {
	lightning = &sparko.Client{
		Client: {
			SparkURL:   s.SparkoURL,
			SparkToken: s.SparkoToken,
		},
	}
}
