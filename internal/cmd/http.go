package cmd

import (
	"log"
	"encoding/json"
	"net/http"
	"net/url"
	"bytes"
	"time"
	"io"
)

func request(
	method string,
	path string,
	headersPtr *Headers,
	bodyPtr *[]byte,
) (*http.Request, error) {
	config, err := GetConfig();
	if err != nil {
		log.Fatal("Could not load user configuration", err)
	}

	var body io.Reader = nil

	if bodyPtr != nil {
		body = bytes.NewBuffer(*bodyPtr)
	}

	req, reqErr := http.NewRequest(method, config.JiraURL + path, body)

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

func get_ticket(ticketId string) (*Ticket, error) {
	req, reqErr := request("GET", "/rest/api/2/issue/" + url.QueryEscape(ticketId), nil, nil)

	if reqErr != nil {
		return nil, reqErr
	}

	client := &http.Client{ Timeout: 30 * time.Second }
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

func get_issue_transitions(ticketId string) (*TicketTransitions, error) {
	req, reqErr := request("GET", "/rest/api/2/issue/" + url.QueryEscape(ticketId) + "/transitions", nil, nil)

	if reqErr != nil {
		return nil, reqErr
	}

	client := &http.Client{ Timeout: 30 * time.Second }
	resp, resErr := client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	parsed := JiraIssueTransitions{}
	unmarshalError := json.Unmarshal(body, &parsed)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	out := &TicketTransitions{
		TicketID:		ticketId,
		Transitions:	parsed.Transitions,
	}

	return out, nil
}

func get_matching_users_search(partialUserName string) (*[]JiraUser, error) {
	req, reqErr := request("GET", "/rest/api/2/user/search?username=" + url.QueryEscape(partialUserName), nil, nil)

	if reqErr != nil {
		return nil, reqErr
	}

	client := &http.Client{ Timeout: 30 * time.Second }
	resp, resErr := client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	parsed := &[]JiraUser{}
	unmarshalError := json.Unmarshal(body, parsed)

	return parsed, unmarshalError
}

func post_ticket_comment(ticketId string, comment NewComment) (bool, error) {
	bodyData, marshalError := json.Marshal(comment)
	if marshalError != nil {
		return false, marshalError
	}

	req, reqErr := request(
		"POST",
		"/rest/api/2/issue/" + url.QueryEscape(ticketId) + "/comment",
		&Headers{
			"Content-Type": "application/json",
		},
		&bodyData,
	)
	if reqErr != nil {
		return false, reqErr
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, resErr := client.Do(req)
	if resErr != nil {
		return false, resErr
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	return body != nil, readErr
}

func post_ticket_transition(ticketId string, transitionId string) (bool, error) {
	bodyData := []byte(`{"transition": {"id": "`+transitionId+`"}}`)

	req, reqErr := request(
		"POST",
		"/rest/api/2/issue/" + url.QueryEscape(ticketId) + "/transitions",
		&Headers{
			"Content-Type": "application/json",
		},
		&bodyData,
	)
	if reqErr != nil {
		return false, reqErr
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, resErr := client.Do(req)
	if resErr != nil {
		return false, resErr
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	return body != nil, readErr
}

func put_ticket_fields(ticketId string, bodyData []byte) (bool, error) {
	req, reqErr := request(
		"PUT",
		"/rest/api/2/issue/" + url.QueryEscape(ticketId),
		&Headers{
			"Content-Type": "application/json",
		},
		&bodyData,
	)
	if reqErr != nil {
		return false, reqErr
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, resErr := client.Do(req)
	if resErr != nil {
		return false, resErr
	}

	defer resp.Body.Close()

	_, readErr := io.ReadAll(resp.Body)
	return readErr == nil, readErr
}
