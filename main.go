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
	store := storage.NewTOMLStorage[Sample](
		"data/storage.txt",
		storage.WithLockPath("data/storage.lock"),
		storage.WithCheckInterval(1*time.Millisecond),
	)
	err := store.Save(Sample{})
	if err != nil {
		log.Println("Error saving:", err)
		return
	}

	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(1+rand.Intn(100)))
			storage := storage.NewTOMLStorage[Sample](
				"data/storage.txt",
				storage.WithLockPath("data/storage.lock"),
				storage.WithCheckInterval(1*time.Millisecond),
			)

			file, err := storage.Open()
			if err != nil {
				log.Println("Error opening:", err)
			}
			defer storage.Close()

			content, err := file.Read()
			if err != nil {
				log.Println("Error reading:", err)
			}
			content.Number += i
			if err := file.Write(content); err != nil {
				log.Println("Error writing:", err)
			}
			log.Println("Written:", content.Number)
		}()
	}
	wg.Wait()

	content, err := store.Load()
	if err != nil {
		log.Println("Error loading:", err)
		return
	}

	log.Println("Loaded:", content.Number)
}
