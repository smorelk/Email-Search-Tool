package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"runtime/pprof"
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
	payload["index"] = "enron_mail"
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

	tgzFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	gz, err := gzip.NewReader(tgzFile)
	if err != nil {
		log.Fatal(err)
	}

	tarArchive := tar.NewReader(gz)
	batch := make([]MailRecord, 0, batchSize)
	sent := 0

	for {
		h, err := tarArchive.Next()

		// End-of-file
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		switch h.Typeflag {
		// If regular file, parse it and add it to current batch.
		case tar.TypeReg:
			msg, err := mail.ReadMessage(tarArchive)
			if err != nil {
				// Malformed header errors may occur, skip those.
				continue
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
				Content: string(content)}
			batch = append(batch, mr)
			if len(batch) == batchSize {
				if err = sendBatch(batch); err != nil {
					log.Fatal(err)
				}
				// Clear current batch to make room for next.
				sent += len(batch)
				batch = nil
			}

		default:
			// Not a file, skip.
			continue
		}

	}

	// Batch might have remaining records at this point,
	// ensure all data is sent.
	if err = sendBatch(batch); err != nil {
		log.Fatal(err)
	}
	sent += len(batch)

	fmt.Println("Emails processed:", sent)
}
