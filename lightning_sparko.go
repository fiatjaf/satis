package main

import "github.com/fiatjaf/satis/lightning/sparko"

func initializeLightning() {
	lightning = &sparko.Client{
		URL:   s.SparkoURL,
		Token: s.SparkoToken,
	}
}
