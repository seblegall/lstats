package lstats

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/gosuri/uiprogress"
)

//FailureLog represents the response send by the tested url when status code is >= 300.
type FailureLog struct {
	statusCode int
	url        string
	response   string
}

//LoadStats represents a load test configuration and the load test result.
type LoadStats struct {
	Reqs                []*http.Request
	workers             int
	CallsResponseTime   []float64
	AverageResponseTime float64
	FailureCount        int
	FailuresLogs        []FailureLog
}

//NewLoadStats init a new load test.
func NewLoadStats(reqs []*http.Request, workers int) LoadStats {
	return LoadStats{
		Reqs:    reqs,
		workers: workers,
	}
}

//Launch actualy execute the load test using workers as set in the LoadStats
//This function print a progress bar to follow the test current status.
func (ls *LoadStats) Launch() {

	fmt.Println("Starting load test...")
	//init progress bar
	uiprogress.Start()
	bar := uiprogress.AddBar(len(ls.Reqs)).AppendCompleted().PrependElapsed()

	//init chan and mutex
	wg := new(sync.WaitGroup)
	in := make(chan *http.Request, 2*ls.workers)
	mu := &sync.Mutex{}

	for i := 0; i < ls.workers; i++ {
		wg.Add(1)
		go func(tests *LoadStats, bar *uiprogress.Bar) {
			defer wg.Done()
			for req := range in {
				respTime, status := doRequest(req, bar)
				mu.Lock()
				//Aggregate responses time.
				tests.CallsResponseTime = append(tests.CallsResponseTime, respTime.Seconds())
				if status > 300 {
					tests.FailureCount++
				}
				mu.Unlock()
			}
		}(ls, bar)
	}

	for _, req := range ls.Reqs {
		in <- req
	}

	close(in)
	wg.Wait()

	bar.Set(len(ls.Reqs))

	ls.AverageResponseTime = avg(ls.CallsResponseTime)
}

//Print show the load test result using tab writer.
func (ls *LoadStats) Print() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "AVG Response Time\tTotal Calls in error\t")
	fmt.Fprintf(w, fmt.Sprintf("%f\t%d", ls.AverageResponseTime, ls.FailureCount))
	fmt.Fprintln(w)
	w.Flush()
}

//doRequest does the http call using net/http.
func doRequest(req *http.Request, bar *uiprogress.Bar) (time.Duration, int) {
	timeStart := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)

	bar.Incr()

	if err != nil {
		log.Printf("Error fetching: %v", err)
	}
	defer resp.Body.Close()

	return time.Since(timeStart), resp.StatusCode
}

//avg calculate the average of a float slice.
func avg(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total / float64(len(values))

}
