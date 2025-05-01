package utils

import (
	"fmt"
	"net/http"
	"time"
)

var CurrentMaster string

func StartMasterMonitor() {
	go func() {
		for {
			time.Sleep(5 * time.Second)

			if IsMaster {
				continue // مش محتاج يراقب نفسه
			}

			resp, err := http.Get("http://" + CurrentMaster + "/ping")
			if err != nil || resp.StatusCode != 200 {
				fmt.Println("❌ Master is down! Trying to become master...")

				// تحويل الـ Node ده إلى Master مؤقتًا
				IsMaster = true
				fmt.Println("⚠️ [WARNING] Node promoted to TEMP MASTER")
			}
		}
	}()
}
