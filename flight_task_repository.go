package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func listFlightTasks(file string) ([]FlightTask, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	tasks := make([]FlightTask, 0) // remove header
	for i, row := range data {
		if i == 0 && len(row) < 3 {
			continue
		}

		trackPrice, errParse := strconv.ParseFloat(row[2], 64)
		if errParse != nil {
			continue
		}
		tasks = append(tasks, FlightTask{
			Destination: row[0],
			IATACode:    row[1],
			TrackPrice:  trackPrice,
		})
	}
	return tasks, nil
}
