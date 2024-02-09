package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

// Constants
const (
	user      = "admin"
	pass      = "Complexpass#123"
	batchSize = 1000 // Emails batch size
	api       = "http://localhost:4080/api/_bulkv2"
)

type MailRecord struct {
	To      []string
	From    string
	Date    string
	Subject string
	Cc      []string
	Content string
}

func sendBatch(batch []MailRecord) error {
	payload := make(map[string]any)
	payload["index"] = "mails"
	payload["records"] = batch

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// If not OK, something went wrong.
	if resp.StatusCode != 200 {
		return err
	}

	return nil
}

func main() {
	if len(os.Args[1:]) == 0 {
		panic("usage: indexer <email directory>")
	}

	// Set-up CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	batch := make([]MailRecord, 0, batchSize)
	sent := 0
	start := time.Now()

	fmt.Println("Indexing started...")
	fmt.Println("Collection CPU profile...")

	err = filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			msg, err := mail.ReadMessage(f)
			if err != nil {
				// Malformed header errors may occur, skip those.
				return nil
			}
			content, err := io.ReadAll(msg.Body)
			if err != nil {
				log.Fatal(err)
			}

			mr := MailRecord{
				To:      msg.Header["To"],
				From:    msg.Header.Get("From"),
				Date:    msg.Header.Get("Date"),
				Subject: msg.Header.Get("Subject"),
				Cc:      msg.Header["Cc"],
				Content: string(content),
			}

			batch = append(batch, mr)
			if len(batch) == batchSize {
				if err = sendBatch(batch); err != nil {
					return err
				}
				// Clear current batch to make room for next.
				sent += len(batch)
				batch = nil
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Batch might have remaining records at this point,
	// ensure all data is sent.
	if err = sendBatch(batch); err != nil {
		log.Fatal(err)
	}
	sent += len(batch)

	fmt.Println("Emails processed:", sent)
	fmt.Println("CPU Profile output:", f.Name())
	fmt.Println("Time elapsed:", time.Since(start).Seconds(), "seconds")
}
