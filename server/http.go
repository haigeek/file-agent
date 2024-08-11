package server

import (
	"fileagent/config"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func StartHttp(port int64) {
	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	r.HandleFunc("/print", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.RequestURI: %s\n", r.RequestURI)
		fmt.Printf("r.URL.Path: %s\n", r.URL.Path)
		fmt.Printf("r.URL.Host: %s\n", r.URL.Host)
		fmt.Printf("r.URL.Hostname(): %s\n", r.URL.Hostname())
		fmt.Printf("r.Method: %s\n", r.Method)
		fmt.Printf("r.URL.Scheme: %s\n", r.URL.Scheme)
		fmt.Printf("gohttpd pid: %d\n", os.Getpid())

		fmt.Printf("\nHeaders: \n")
		for name, values := range r.Header {
			if len(values) == 1 {
				fmt.Printf("%s: %v\n", name, values[0])
				continue
			}

			fmt.Println(name)
			for i := 0; i < len(values); i++ {
				fmt.Printf("  - #%d: %s\n", i, values[i])
			}
		}

		fmt.Printf("\nPayload: \n")
		defer r.Body.Close()
		bs, _ := io.ReadAll(r.Body)
		fmt.Println(string(bs))
		fmt.Fprintln(w, "ok")
	})

	r.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "r.RequestURI: %s\n", r.RequestURI)
		fmt.Fprintf(w, "r.URL.Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "r.URL.Host: %s\n", r.URL.Host)
		fmt.Fprintf(w, "r.URL.Hostname(): %s\n", r.URL.Hostname())
		fmt.Fprintf(w, "r.Method: %s\n", r.Method)
		fmt.Fprintf(w, "r.URL.Scheme: %s\n", r.URL.Scheme)
		fmt.Fprintf(w, "gohttpd pid: %d\n", os.Getpid())

		fmt.Fprintf(w, "\nHeaders: \n")
		for name, values := range r.Header {
			if len(values) == 1 {
				fmt.Fprintf(w, "%s: %v\n", name, values[0])
				continue
			}

			fmt.Fprintln(w, name)
			for i := 0; i < len(values); i++ {
				fmt.Fprintf(w, "  - #%d: %s\n", i, values[i])
			}
		}

		fmt.Fprintf(w, "\nPayload: \n")
		defer r.Body.Close()
		bs, _ := io.ReadAll(r.Body)
		fmt.Fprintln(w, string(bs))
	})

	r.HandleFunc("/read", getFile)

	r.HandleFunc("/update", updateFile)

	n := negroni.New()
	n.UseHandler(r)

	log.Println("listening http on", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), n))
}

func getFile(w http.ResponseWriter, r *http.Request) {
	auth := config.Conf.Auth
	if !checkAuth(auth, w, r) {
		return
	}
	// get file param from url
	fileKey := r.URL.Query().Get("file")
	//check from config files
	if fileKey == "" {
		http.Error(w, "file param is required", http.StatusBadRequest)
		return
	}
	files := config.Conf.Files
	filePath, exists := files[fileKey]
	if exists {
		//read file
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var content []byte
		content, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.Write(content)
	} else {
		http.Error(w, "file is not exist", http.StatusBadRequest)
	}
}

func updateFile(w http.ResponseWriter, r *http.Request) {
	auth := config.Conf.Auth
	if !checkAuth(auth, w, r) {
		return
	}
	// get file param from url
	fileKey := r.URL.Query().Get("file")
	//check from config files
	if fileKey == "" {
		http.Error(w, "file param is required", http.StatusBadRequest)
		return
	}
	files := config.Conf.Files
	filePath, exists := files[fileKey]
	if exists {
		// read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		// replace file content
		err = replaceFileContent(filePath, body)
		if err != nil {
			http.Error(w, "Failed to replace file content", http.StatusInternalServerError)
			return
		}

		//read file
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var content []byte
		content, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.Write(content)
	}
}

func replaceFileContent(filename string, content []byte) error {
	// read file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// replace file content
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func checkAuth(authstr string, w http.ResponseWriter, r *http.Request) bool {
	if authstr == "" {
		return true
	}

	auth := r.Header.Get("FA-Token")
	if strings.TrimSpace(auth) == "" {
		http.Error(w, "FA-Token is blank", http.StatusUnauthorized)
		return false
	}

	if auth != authstr {
		http.Error(w, "FA-Token invalid", http.StatusUnauthorized)
		return false
	}
	return true
}
