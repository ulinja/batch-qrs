package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"io"
	"os"
	"runtime"
	"sync"
)

type job struct {
	index   int
	payload string
}

type result struct {
	index int
	data  string
	err   error
}

func main() {
	var args = os.Args[1:]
	if len(args) < 1 {
		errorExit("at least one argument is required")
	}

	numWorkers := runtime.NumCPU()
	jobs := make(chan job, len(args))
	results := make(chan result, len(args))

	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	for i, arg := range args {
		jobs <- job{index: i, payload: arg}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	orderedResults := make([]string, len(args))
	for r := range results {
		if r.err != nil {
			errorExit(r.err.Error())
		}
		orderedResults[r.index] = r.data
	}

	for _, data := range orderedResults {
		fmt.Println(data)
	}
}

func worker(jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		bytes := generateQrCode(j.payload)
		results <- result{
			index: j.index,
			data:  base64Encode(bytes),
			err:   nil,
		}
	}
}

func errorExit(message string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	os.Exit(1)
}

func generateQrCode(payload string) []byte {
	qrc, err := qrcode.NewWith(
		payload,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium),
	)
	if err != nil {
		errorExit(err.Error())
	}

	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(
		wr,
		standard.WithQRWidth(10),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithBorderWidth(20),
	)
	if err = qrc.Save(w2); err != nil {
		errorExit(err.Error())
	}

	return buf.Bytes()
}

func base64Encode(payload []byte) string {
	return base64.StdEncoding.EncodeToString(payload)
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
