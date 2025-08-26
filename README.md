# go-simple-file-storage

`go-simple-file-storage` implements a simple data storage with single file.


## Installation

```sh
go get -u github.com/NovanMondon/go-simple-file-storage/storage
```

## Usage

### simple use

```go
package main

import (
	"log"

	"github.com/NovanMondon/go-simple-file-storage/storage"
)

type SampleData struct {
	Number int `json:"number"`
}

func main() {
	storage := storage.NewJSONStorage[*SampleData](
		"data/storage.txt",
	)

	err := storage.Save(&SampleData{Number: 100})
	if err != nil {
		log.Println("Error saving:", err)
		return
	}

	data, err := storage.Load()
	if err != nil {
		log.Println("Error loading:", err)
		return
	}

	log.Println("Loaded:", data.Number)
}

```

### complex use

```go
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
				storage.WithRetryMax(100),
			)

			file, err := storage.Open()
			if err != nil {
				log.Println("Error opening:", err)
				return
			}
			defer storage.Close()

			content, err := file.Read()
			if err != nil {
				log.Println("Error reading:", err)
				return
			}
			content.Number += i
			if err := file.Write(content); err != nil {
				log.Println("Error writing:", err)
				return
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

```

## License

`go-simple-file-storage` is licensed under the MIT License – see the [LICENSE](./LICENSE) file for details.

`go-simple-file-storage` depends on these - see [LICENSES](./LICENSES/) for details.
- https://github.com/gofrs/flock