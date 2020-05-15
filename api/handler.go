package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Handler struct {
}

func (s *Handler) CallURL(req *URLMessage, stream API_CallURLServer) error {
	if req.Url == "" {
		return fmt.Errorf("URL is required")
	}

	_, err := url.ParseRequestURI(req.Url)
	if err != nil {
		return fmt.Errorf("URL is not valid: %v", err)
	}

	resp, err := http.Get(req.Url)
	if err != nil {
		return fmt.Errorf("fail HTTP call to %v", req.Url)
	}
	defer resp.Body.Close()

	bytesHeaders, err := json.Marshal(resp.Header)
	if err != nil {
		return fmt.Errorf("fail marshal headers: %v", err)
	}

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("fail read body: %v", err)
	}

	var (
		chunkSize = 1024
		buffer    = make([]byte, chunkSize)
	)

	reader := io.MultiReader(
		bytes.NewReader(bytesHeaders),
		bytes.NewReader(bytesBody),
	)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return fmt.Errorf("failed raeding from body to buffer: %v", err)
		}

		log.Printf("sent chunk in %v bytes \n", n)

		err = stream.Send(&ResponseMessage{
			Response: buffer[:n],
		})
		if err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}
	}

	return nil
}
