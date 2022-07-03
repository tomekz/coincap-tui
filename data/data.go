package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Data struct {
	UserId    int
	Id        int
	Title     string
	Completed bool
}

type DataFetchError struct {
	Err error
}

func (e DataFetchError) Error() string {
	return e.Err.Error()
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

const url = "https://jsonplaceholder.typicode.com"

func FetchData(id string) (*Data, error) {
	data := &Data{}
	error := getJson(fmt.Sprintf("%s/todos/%s", url, id), data)
	if error != nil {
		return nil, DataFetchError{Err: error}
	}
	return data, nil
}
