package data

import (
	"fmt"

	"github.com/go-resty/resty/v2"
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

const url = "https://jsonplaceholder.typicode.com"

func FetchData(id string) (*Data, error) {
	data := &Data{}
	client := resty.New()
	_, err := client.R().SetResult(data).Get(fmt.Sprintf("%s/todos/%s", url, id))

	if err != nil {
		return nil, DataFetchError{Err: err}
	}
	return data, nil
}
