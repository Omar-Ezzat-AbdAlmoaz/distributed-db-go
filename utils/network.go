package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	//"log"

	"net/http"
	//"os/exec"
)

func SendPostRequest(url string, body interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling body:", err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request to", url, ":", err)
		return err
	}
	defer resp.Body.Close()

	// اطبع الاستجابة
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response from", url, ":", string(respBody))

	if resp.StatusCode >= 400 {
		fmt.Println("Request failed with status:", resp.Status)
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return nil
}

type Notification struct {
	Message string `json:"message"`
}

func BroadcastNotification(message string) {
	body := Notification{Message: message}

	for _, node := range append([]NodeConfig{{Address: CurrentMaster}}, OtherNodes...) {
		go func(n NodeConfig) {
			err := SendPostRequest("http://"+n.Address+"/notify", body)
			if err != nil {
				fmt.Println("Failed to notify", n.Address, ":", err)
			}
		}(node)
	}
}
