package main

import "github.com/fiatjaf/satis/storage/filedb"

func initializeStorage() {
	store = &filedb.FileDatabase{
		Path: s.FileDBPath,
	}

	if err := store.Init(); err != nil {
		panic(err)
	}
}
