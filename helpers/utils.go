package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

func SlackHook(slackURL string, text string) error {
	bodyRequest, err := json.Marshal(map[string]interface{}{
		"text": text,
	})
	req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(bodyRequest))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	return nil
}

func SubStringBodyResponse(obj string, limit int) string {
	if len(obj) > limit {
		return obj[0:limit]
	}
	return obj
}

var debugCounter uint64

func nextDebugId() string {
	return fmt.Sprintf("%d", atomic.AddUint64(&debugCounter, 1))
}
func BuildFileName() string {
	return time.Now().Format("20060102150405") + "_" + nextDebugId()
}
