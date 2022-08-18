package barkserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/c1emon/barkbridge/src/log"
)

type Message struct {
	Title             string `json:"title"`
	Body              string `json:"body"`
	DeviceKey         string `json:"device_key"`
	Category          string `json:"category"`
	Badge             int    `json:"badge,omitempty"`
	Sound             string `json:"sound,omitempty"`
	Icon              string `json:"icon,omitempty"`
	Group             string `json:"group,omitempty"`
	URL               string `json:"url,omitempty"`
	Level             string `json:"level,omitempty"`
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty"`
	Copy              string `json:"copy,omitempty"`
	IsArchive         string `json:"isArchive,omitempty"`
}

func Push(server string, message Message) {
	log.Info()
	msg, err := json.Marshal(message)
	if err != nil {
		print("")
	}

	client := &http.Client{}
	// Create request
	req, err := http.NewRequest("POST", server, bytes.NewBuffer(msg))
	if err != nil {
		fmt.Println("Failure : ", err)
	}
	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}
