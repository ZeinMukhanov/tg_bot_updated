package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func getUpdates(token string, offset int) ([]update, error) {
	url := apiURL + token + "/getUpdates"
	response, err := http.Get(url + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		OK     bool     `json:"ok"`
		Result []update `json:"result"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, errors.New("getUpdates request failed")
	}

	return result.Result, nil
}

func sendMessage(token string, chatID int, text string) error {
	url := apiURL + token + "/sendMessage"

	requestBody, err := json.Marshal(sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return nil
}
