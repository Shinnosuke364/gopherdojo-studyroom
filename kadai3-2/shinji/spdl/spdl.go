package spdl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

	// parallel download
	for i := 0; i < procs; i++ {

		// make range
		r := makeRange(i, split, procs, filesize)

		// download
		errChan := make(chan error, 1)

		go func() {
			err := dlpart(r, url, filepath)
			errChan <- err
		}()

		err := <-errChan
		if err != nil {
			return err
		}
	}

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

func dlpart(r Range, url string, filepath string) error {
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

	// copy
	if _, err := io.Copy(output, res.Body); err != nil {
		return err
	}

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
