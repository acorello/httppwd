package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type loggingFileServer struct {
	fs http.FileSystem
}

func (lfs *loggingFileServer) Open(name string) (http.File, error) {
	file, err := lfs.fs.Open(name)
	if err != nil {
		log.Printf("Failed to open file: %s, error: %v\n", name, err)
	} else {
		log.Printf("Serving file: %s\n", name)
	}
	return file, err
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}
	directory := flag.String("dir", pwd, "directory to serve; defaults to $PWD")
	port := flag.Int("port", 8080, "TCP port to listen on; defaults to 8080")
	flag.Parse()

	rootDir := validRootDir(directory)

	loggingFS := &loggingFileServer{fs: rootDir}
	http.Handle("/", http.FileServer(loggingFS))

	address := fmt.Sprintf(":%d", *port)
	log.Println("Listening on", address)
	log.Fatal(http.ListenAndServe(address, nil))
	// TODO: verify that symlinks are not served unless explicitly allowed with CLI option
}

func validRootDir(directory *string) http.Dir {
	info, err := os.Stat(*directory)
	if err != nil {
		log.Fatalf("Failed to access folder: %s, error: %v\n", *directory, err)
	}
	if !info.IsDir() {
		log.Fatalf("%s is not a directory\n", *directory)
	}
	rootDir := http.Dir(*directory)
	return rootDir
}
