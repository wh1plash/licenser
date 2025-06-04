package appchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type App struct {
	Name  string    `json:"name"`
	Until time.Time `json:"until"`
}

func Validate(name string) {
	until, err := CheckName("http://localhost:9080", name)
	if err != nil {
		log.Fatalf("❌ App check failed: %v", err)
	}
	if time.Now().After(until) {
		log.Fatal("expired license")
	}
}

func CheckName(serverURL, appName string) (time.Time, error) {
	body := map[string]string{
		"app": appName,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodGet, serverURL+"/app", bytes.NewBuffer(data))
	if err != nil {
		return time.Time{}, fmt.Errorf("error to build get request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("get request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("shutting down")
		//return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var app App
	if err := json.Unmarshal(respBody, &app); err != nil {
		return time.Time{}, err
	}

	return app.Until, nil
}
