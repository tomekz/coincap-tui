package tui

import (
	"encoding/json"
	"net/http"
	"time"
)

type Data struct {
	UserId    int
	Id        int
	Title     string
	Completed bool
}

func GetJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
