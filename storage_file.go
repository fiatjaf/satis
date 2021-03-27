package main

import "github.com/fiatjaf/satis/storage/filedb"

func initializeStorage() {
	storage = filedb.FileDatabase{
		Path: s.FileDBPath,
	}

	if err := storage.Init(); err != nil {
		panic(err)
	}
}
