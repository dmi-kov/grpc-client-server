package api

import (
	"encoding/json"
	"fmt"
	"io"
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

	var (
		chunkSize       = 1024
		buffer          = make([]byte, chunkSize)
		needSendHeaders = true
	)

	for {
		if !needSendHeaders {
			bytesHeaders = []byte{}
			needSendHeaders = false
		}

		n, err := resp.Body.Read(buffer)
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
			Headers:  bytesHeaders,
		})
		if err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}
	}

	return nil
}
