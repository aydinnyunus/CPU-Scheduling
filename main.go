package main

import (
	"fmt"
	"math/rand"
	"time"
	_ "time"
)

const (
	NUMBER_OF_PROCESS = 100
)

type process struct {
	actualBurstTime, estimatedBurstTime, arrivalTime float64
}

func newProcess(arrivalTime float64) *process {

	p := process{arrivalTime: arrivalTime}
	return &p
}

var processList []*process
var actualBurstTime []float64
var estimatedBurstTime []float64
var arrivalTime []int64
var alpha =0.4

func calculateBurst(i int64) float64{
	process := alpha * processList[i-1].actualBurstTime + (1-alpha) * processList[i-1].estimatedBurstTime
	estimatedBurstTime = append(estimatedBurstTime, process)
	return process
}

func main() {
	/*
	Burst time(nth process) = alpha * actual burst time (n-1)th process + (1-alpha) * estimated burst time(n-1)th process

	Where alpha is a constant ; 0 <= alpha =< 1
	 */
	for i := 1 ; i <= NUMBER_OF_PROCESS; i++ {
		s1 := rand.NewSource(time.Now().UnixNano())

		r1 := rand.New(s1)
		burst := r1.Float64()*100
		fmt.Printf("\nActual Burst Time : %f\n", burst)

		now := time.Now()
		sec := now.Unix()      // number of seconds since January 1, 1970 UTC
		fmt.Printf("Arrival Time : %d\n", sec)
		process_1 := newProcess(float64(sec))
		if len(processList) > 1 {
			process_1.estimatedBurstTime = calculateBurst(int64(len(processList) - 1))
		}
		process_1.actualBurstTime = burst
		processList = append(processList, process_1)
		fmt.Printf("Process %d Estimated Burst Time : %f\n", i,process_1.estimatedBurstTime)
	}
	fmt.Println(processList)



}