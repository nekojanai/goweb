package index

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func Index() {
	args := os.Args
	if args[3] != "" {
		EnumerateDirectory(args[3])
	} else {
		fmt.Printf("Usage: go run indexer.go <path_to_index>")
	}
}

type DirectoryIndex struct {
	Quantity   uint64
	FullPath   string
	DirEntries []fs.DirEntry
}

func (index *DirectoryIndex) UpdateIndexFile() {
	buffer := new(bytes.Buffer)
	err := toml.NewEncoder(buffer).Encode(map[string]interface{}{
		"Quantity": index.Quantity,
	})
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(index.FullPath, os.O_CREATE, 0750)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write()

}

func EnumerateDirectory(path string) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", dirEntry)

}
