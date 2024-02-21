package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"strings"
	"time"
)

// Constants
const (
	user      = "admin"
	pass      = "Complexpass#123"
	batchSize = 1000 // Emails batch size
	api       = "http://localhost:4080/api/_bulkv2"
)

var mailDir = flag.String("maildir", ".", "Mail directory to index")
var cpuProf = flag.String("prof", "cpu.prof", "CPU profile output file")

// Regexp patterns
var msg = regexp.MustCompile("Message-ID:.*")
var to = regexp.MustCompile(`To:.*`)
var from = regexp.MustCompile(`From:.*`)
var date = regexp.MustCompile(`Date:.*`)
var subj = regexp.MustCompile(`Subject:.*`)
var endHeader = regexp.MustCompile(`\r\n\r\n`)

type MailRecord struct {
	MessageID string
	To        []string
	From      string
	Date      string
	Subject   string
	Content   string
}

// Parse file bytes into a MailRecord
func parse(data string) MailRecord {
	// If there is not message ID in data,
	// skip because there is no email data
	// to parse.
	if !msg.MatchString(data) {
		return MailRecord{}
	}

	cLoc := endHeader.FindStringIndex(data)
	msg_ := msg.FindString(data)
	to_ := strings.Split(to.FindString(data), ",")
	to_[0] = strings.Replace(to_[0], "To:", "", 1)
	to_[len(to_)-1] = strings.ReplaceAll(to_[len(to_)-1], ",", "")
	from_ := from.FindString(data)
	date_ := date.FindString(data)
	subj_ := subj.FindString(data)

	return MailRecord{
		MessageID: strings.Replace(msg_, "Message-ID:", "", 1),
		To:        to_,
		From:      strings.Replace(from_, "From:", "", 1),
		Date:      strings.Replace(date_, "Date:", "", 1),
		Subject:   strings.Replace(subj_, "Subject:", "", 1),
		Content:   data[cLoc[0]:],
	}
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
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: indexer -maildir <dir> -prof <cpu.prof>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set-up CPU profiling
	f, err := os.Create(*cpuProf)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	batch := make([]MailRecord, 0, batchSize)
	start := time.Now()

	fmt.Println("Indexing started...")
	fmt.Println("Collecting CPU profile...")

	err = filepath.WalkDir(*mailDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			fbytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			mr := parse(string(fbytes))
			batch = append(batch, mr)
			if len(batch) == batchSize {
				if err := sendBatch(batch); err != nil {
					log.Fatal(err)
				}

				batch = nil
				batch = make([]MailRecord, 0, batchSize)
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

	fmt.Println("Indexing Finished.")
	fmt.Println("CPU Profile output:", f.Name())
	fmt.Println("Time elapsed:", time.Since(start).Seconds(), "seconds")
}
