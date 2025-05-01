package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendPostRequest(url string, body interface{}) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("❌ Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Error sending POST to", url, ":", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("📤 Sent to", url, "- Status:", resp.Status)
}
