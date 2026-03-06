package cmd

import (
	"log"
	"encoding/json"
	"net/http"
	"time"
	"io"
)

func request(
	method string,
	path string,
	headersPtr *Headers,
) (*http.Request, error) {
	config, err := GetConfig();
	if err != nil {
		log.Fatal("Could not load user configuration", err)
	}

	req, reqErr := http.NewRequest(method, config.JiraURL + path, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	if headersPtr != nil {
		headers := *headersPtr;

		if _, ok := headers["Authorization"]; !ok {
			req.Header.Add("Authorization", "Bearer " + config.Token)
		}

		if _, ok := headers["Accept"]; !ok {
			req.Header.Add("Accept", "application/json");
		}

		for k, v := range headers {
			if len(k) > 0 && len(v) > 0 {
				req.Header.Add(k, v)
			}
		}
	} else {
		req.Header.Add("Authorization", "Bearer " + config.Token)
		req.Header.Add("Accept", "application/json");
	}

	return req, nil;
}

func getTicket(ticketId string) (*Ticket, error) {
	req, reqErr := request("GET", "/rest/api/2/issue/" + ticketId, nil)

	if reqErr != nil {
		return nil, reqErr
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, resErr := client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	parsed := &Ticket{}
	unmarshalError := json.Unmarshal(body, parsed)

	return parsed, unmarshalError
}
