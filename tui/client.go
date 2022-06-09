package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Data struct {
	UserId    int
	Id        int
	Title     string
	Completed bool
}

func getJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func GetData(id string) tea.Msg {
	data := &Data{}
	getJson(fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%s", id), data)
	fmt.Printf("%+v", data)
	return data
}
