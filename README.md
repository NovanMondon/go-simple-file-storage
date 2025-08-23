# go-simple-file-storage

`go-simple-file-storage` implements a simple data storage with single file.


## Installation

```sh
go get -u github.com/NovanMondon/go-simple-file-storage/storage
```

## Usage

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
	storage := storage.New[*SampleData](
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

## License

`go-simple-file-storage` is licensed under the MIT License – see the [LICENSE](./LICENSE) file for details.

`go-simple-file-storage` depends on these - see [LICENSES](./LICENSES/) for details.
- https://github.com/gofrs/flock