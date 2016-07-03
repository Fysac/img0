package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const ListenOn = ":8000"
const LoggingEnabled = true

const ReencodeJPEG = true

const ImageURL = "/i/"
const UploadURL = "/upload"

const HTMLDir = "./html/"
const UploadDir = "./img/"

func upload(w http.ResponseWriter, r *http.Request) {
	const MaxRequestSize = 10 << 20
	const MaxInMemory = 5 << 20

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
		r.ParseMultipartForm(MaxInMemory)

		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}
		defer file.Close()

		// Validate image type and determine extension
		var ext string
		switch handler.Header.Get("Content-Type") {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			fmt.Fprintf(w, "Image must be JPEG, PNG, or GIF.")
			return
		}

		// Create unique, human-friendly filename
		h := sha256.New()
		h.Write([]byte(handler.Filename + time.Now().String()))
		name := strings.Replace(base64.StdEncoding.EncodeToString(h.Sum(nil)), "+", "", -1)
		name = strings.Replace(name, "/", "", -1)[:12] + ext

		// Write image to local directory
		f, err := os.OpenFile(UploadDir+name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Something went wrong on the backend.")
			return
		}
		defer f.Close()

		if ext == ".jpg" && ReencodeJPEG {
			// Re-encode jpeg to remove metadata
			img, _, _ := image.Decode(file)
			jpeg.Encode(f, img, nil)
		} else {
			io.Copy(f, file)
		}

		http.Redirect(w, r, ImageURL+name, 302)
	}
}

func hideDirListing(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != ImageURL {
			handler.ServeHTTP(w, r)
		}
	}
}

func logHTTP(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if LoggingEnabled {
			log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}
		handler.ServeHTTP(w, r)
	}
}

func main() {
	os.Mkdir(UploadDir, 0700)

	http.Handle("/", http.FileServer(http.Dir(HTMLDir)))
	http.Handle(ImageURL, hideDirListing(http.StripPrefix(ImageURL, http.FileServer(http.Dir(UploadDir)))))
	http.HandleFunc(UploadURL, upload)
	http.ListenAndServe(ListenOn, logHTTP(http.DefaultServeMux))
}
