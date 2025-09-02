package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/NovanMondon/go-simple-file-storage/storage"
)

type Sample struct {
	Number int `json:"number"`
}

func main() {
	// initialize
	store := storage.NewTOMLStorage[Sample](
		"data/storage.txt",
		storage.WithLockPath("data/storage.lock"), // specify lockfile
	)

	//save
	err := store.Save(Sample{})
	if err != nil {
		log.Println("Error saving:", err)
		return
	}

	count := 100

	// note: Control parallel processing with waitgroup
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func() {
			defer wg.Done()

			// wait random length
			time.Sleep(time.Millisecond * time.Duration(1+rand.Intn(100)))

			// initialize
			storage := storage.NewTOMLStorage[Sample](
				"data/storage.txt",
				storage.WithLockPath("data/storage.lock"),
				storage.WithRetryInterval(1*time.Millisecond),
				storage.WithRetryMax(100),
			)

			// get lock
			file, err := storage.Open()
			if err != nil {
				log.Println("Error opening:", err)
				return
			}
			defer storage.Close() // release lock

			// read from locked file
			content, err := file.Read()
			if err != nil {
				log.Println("Error reading:", err)
				return
			}

			content.Number += i

			// write to locked file
			if err := file.Write(content); err != nil {
				log.Println("Error writing:", err)
				return
			}
			log.Println("Written:", content.Number)
		}()
	}
	wg.Wait()

	// load
	content, err := store.Load()
	if err != nil {
		log.Println("Error loading:", err)
		return
	}

	log.Println("Loaded:", content.Number)

	// cleanup
	store.Cleanup()
}
