package main

import (
	"os"
	"reflect"
	"testing"
)

func TestReadJSONFile(t *testing.T) {
	tempFile, cleanup := createTempFile(t, []byte(`[
		{"a": 1, "b": 3},
		{"a": 5, "b": -9},
		{"a": -2, "b": 4}
	]`))
	defer cleanup()

	items, err := readJSONFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Error reading JSON file: %v", err)
	}

	expected := []Item{
		{A: 1, B: 3},
		{A: 5, B: -9},
		{A: -2, B: 4},
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Wrong result. Expected: %v, received: %v", expected, items)
	}
}

func TestProcessItems(t *testing.T) {
	testItems := []Item{
		{A: 1, B: 3},
		{A: 5, B: -9},
		{A: -2, B: 4},
	}

	result := processItems(testItems, 2)

	expected := 2
	if result != expected {
		t.Errorf("Wrong result. Expected: %d, received: %d", expected, result)
	}
}

// Вспомогательная функция для создания временного файла с тестовыми данными
func createTempFile(t *testing.T, data []byte) (*os.File, func()) {
	tempFile, err := os.CreateTemp("", "test_data_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	return tempFile, func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}
}
