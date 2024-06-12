package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Item struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {

}

// readJSONFile читает JSON файл и преобразует его содержимое в срез структур Item.
// Возвращает срез items и ошибку, если что-то пошло не так.
func readJSONFile(filename string) ([]Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var items []Item
	err = json.Unmarshal(byteValue, &items)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return items, nil
}
