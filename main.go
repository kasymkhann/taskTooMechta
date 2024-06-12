package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Item Структура для хранения данных из JSON
type Item struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	var numWorkers int
	flag.IntVar(&numWorkers, "workers", 4, "number of workers")
	flag.Parse()

	if numWorkers <= 0 {
		log.Fatalf("The number of workers must be greater than 0")
	}

	items, err := readJSONFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	totalSum := processItems(items, numWorkers) // Запускаем обработку элементов

	fmt.Printf("Total amount: %d\n", totalSum)
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

// worker обрабатывает часть элементов,
// суммирует значения полей A и B и отправляет результат в канал results.
func worker(items []Item, results chan<- int, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	sum := 0
	for _, item := range items {
		sum += item.A + item.B
	}
	log.Printf("Worker %d processed %d items\n", workerID, len(items))
	results <- sum
}

// processItems разбивает элементы на части и распределяет их между горутинами для параллельной обработки.
func processItems(items []Item, numWorkers int) int {
	chunkSize := (len(items) + numWorkers - 1) / numWorkers
	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(items) {
			end = len(items)
		}

		wg.Add(1)
		go worker(items[start:end], results, &wg, i)
	}

	wg.Wait()
	close(results)
	return collectResults(results)
}

// collectResults собирает результаты из канала и возвращает их общую сумму.
func collectResults(results chan int) int {
	totalSum := 0
	for result := range results {
		totalSum += result
	}
	return totalSum
}
