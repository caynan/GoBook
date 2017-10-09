package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string, 10)

	if os.Args[1] == "-top1M" {
		println("Getting top 1 Million websites...")
		go fetchTop1Million(ch)

		for {
			if val, status := <-ch; status {
				fmt.Println(val)
			} else {
				break
			}
		}

	} else {
		for _, url := range os.Args[1:] {
			go fetch(url, ch) // start a go routine
		}

		for range os.Args[1:] {
			fmt.Println(<-ch) // receive from channel ch
		}
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		// ch <- fmt.Sprint(err) // send to channel ch
		// log.Print(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs \t %7d bytes \t %s", secs, nbytes, url)
}

func fetchTop1Million(ch chan<- string) {
	start := time.Now()
	url := "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	zipfile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	fmt.Println("opening zipfile buffer...\n")

	files, err := zip.NewReader(bytes.NewReader(zipfile), int64(len(zipfile)))
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	csvFile, err := files.File[0].Open()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	reader := csv.NewReader(csvFile)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			ch <- fmt.Sprint(err)
			return
		}
		// Start fetch process
		url := fmt.Sprintf("http://%s", line[1])
		go fetch(url, ch)
	}
	defer close(ch)

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs \t %s", secs, url)
}
