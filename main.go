package main

import (
	"net/http"
	"sync"
)

type FlightTask struct {
	Destination string
	IATACode string
	TrackPrice float64
}


func main() {
	// _ = os.Setenv("TEQUILA_API_KEY", "YOUR_API_KEY")
	tasks, err := listFlightTasks("./data/flight-task.csv")
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	c := &http.Client{}
	amazonNotifier := AmazonSmsNotifier{
		client: snsClient,
	}

	for _, task := range tasks {
		go scheduleFlightTask(wg, c, amazonNotifier, task)
	}
	wg.Wait()
}
