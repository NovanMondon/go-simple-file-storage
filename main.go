package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/NovanMondon/go-simple-file-storage/storage"
)

type SampleData struct {
	Number int `json:"number"`
}

func main() {
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(1+rand.Intn(100)))
			storage := storage.New[*SampleData](
				"data/storage.txt",
				storage.WithLockPath("data/storage.lock"),
				storage.WithCheckInterval(1*time.Millisecond),
			)
			storage.Save(&SampleData{Number: i})
			log.Println("Saved:", i)
		}()
	}
	wg.Wait()

	storage := storage.New[*SampleData](
		"data/storage.txt",
		storage.WithLockPath("data/storage.lock"),
		storage.WithCheckInterval(1*time.Millisecond),
	)
	data, _ := storage.Load()
	log.Println("Loaded:", data.Number)
}
