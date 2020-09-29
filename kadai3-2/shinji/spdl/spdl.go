package spdl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type Range struct {
	first int
	last  int
}

func Download(filepath string, url string, procs int) error {
	filesize, err := getFileSize(url)
	if err != nil {
		return err
	}

	split := filesize / procs

	var wg sync.WaitGroup

	// parallel download
	for i := 0; i < procs; i++ {
		wg.Add(1)

		// make range
		r := makeRange(i, split, procs, filesize)

		// download
		go dlpart(r, url, filepath, &wg)
	}

	wg.Wait()

	return nil
}

func getFileSize(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return 0, err
	}

	return size, err
}

func makeRange(i, split, procs, filesize int) Range {
	var r Range

	r.first = i * split

	if i == procs-1 {
		r.last = filesize
	} else {
		r.last = r.first + split - 1
	}

	return r
}

func dlpart(r Range, url, filepath string, wg *sync.WaitGroup) error {

	var limit int = r.last - r.first + 1
	bar := pb.Full.Start(limit)

	//
	defer wg.Done()

	// make response
	res, err := makeResponse(r, url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// open file
	output, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer output.Close()

	barReader := bar.NewProxyReader(res.Body)
	io.Copy(output, barReader)
	bar.Finish()

	return nil
}

func makeResponse(r Range, url string) (*http.Response, error) {
	// create get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set download ranges
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", r.first, r.last))

	// make response
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
