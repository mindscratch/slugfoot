package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/post", PostOnly(HandlePost))
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type handler func(w http.ResponseWriter, r *http.Request)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Accept %s %s from %s\n", r.Method, r.URL, r.RemoteAddr)
	io.WriteString(w, "hello world\n")
}

func PostOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	//get the multipart reader for the request.
	reader, err := r.MultipartReader()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saving files to %s\n", dir)
	//copy each part to destination.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		//if part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			continue
		}
		dst, err := os.Create(path.Join(dir, part.FileName()))
		defer dst.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("created %s\n", dst.Name())

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//display success message.
	io.WriteString(w, "upload success")
}
