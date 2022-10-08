package barkserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
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
	Id                string `json:"-"`
}

type BarkResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

func Push(server string, message Message) bool {
	msg, err := json.Marshal(message)
	if err != nil {
		logrus.WithField("message", fmt.Sprintf("%+v", message)).Warn(err)
		return false
	}

	logrus.WithFields(logrus.Fields{
		"bark endpont":  server,
		"device key":    message.DeviceKey,
		"message title": message.Title,
	}).Debug("send message")

	client := &http.Client{}
	// Create request
	req, err := http.NewRequest("POST", server, bytes.NewBuffer(msg))
	if err != nil {
		logrus.Warn(fmt.Sprintf("failed new request: %s", err))
		return false
	}
	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		logrus.Warn(fmt.Sprintf("failed do request: %s", err))
		return false
	}

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Warn(fmt.Sprintf("failed read response body: %s", err))
		return false
	}
	var bresp BarkResp
	err = json.Unmarshal(respBody, &bresp)
	if err != nil {
		logrus.Warn(fmt.Sprintf("failed unmarshal response: %s\n%s", err, string(respBody)))
		return false
	}

	if bresp.Code != 200 {
		logrus.Warn(fmt.Sprintf("bad response code %d:\n%s", bresp.Code, bresp.Message))
		return false
	}
	logrus.Info("send success")
	return true
}
