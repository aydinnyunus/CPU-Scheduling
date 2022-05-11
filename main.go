package main

import (
	"fmt"
	"math/rand"
	"time"
	_ "time"
)

const (
	NUMBER_OF_PROCESS = 100
	REALTIME = 0
	HIGH = 1
	NORMAL = 2
	LOW = 3
)

var processList []*process
var _ []float64
var estimatedBurstTime []float64
var _ []int64
var alpha =0.4
var variance = float64(0)

type process struct {
	actualBurstTime, estimatedBurstTime, arrivalTime float64
	priority int64
}

func newProcess(arrivalTime float64) *process {

	p := process{arrivalTime: arrivalTime}
	return &p
}


func calculateBurst(i int64) float64{
	process := alpha * processList[i-1].actualBurstTime + (1-alpha) * processList[i-1].estimatedBurstTime
	estimatedBurstTime = append(estimatedBurstTime, process)
	return process
}

func priorityQueue(){

}

func calculateAvgArr(){

}

func calculateAvgWait(){

}

func calculateAvgTurn(){

}

func calculateAvgQueue(){
	
}

func main() {
	/*
	Burst time(nth process) = alpha * actual burst time (n-1)th process + (1-alpha) * estimated burst time(n-1)th process

	Where alpha is a constant ; 0 <= alpha =< 1

	3 priority queue
	Average Queue Length
	Average Waiting Time
	Average Arrival Time
	Average Turnaround Time
	FIFO Round Robin
	Poisson enter
	Exponential leave
	 */
	now := time.Now()
	sec := float64(now.Unix()) // number of seconds since January 1, 1970 UTC
	s1 := rand.NewSource(time.Now().UnixNano())

	r1 := rand.New(s1)
	last := float64(sec) + r1.Float64() * float64(1000)
	diff := last - float64(sec)

	/*
	//Or however many you might need + buffer.
	c := make(chan int, 300)

	//Push
	c <- value

	//Pop
	x <- c
	 */
	for i := 1 ; i <= NUMBER_OF_PROCESS; i++ {
		s1 := rand.NewSource(time.Now().UnixNano())

		r1 := rand.New(s1)
		burst := r1.Float64()*100
		fmt.Printf("\nActual Burst Time : %f\n", burst)


		if len(processList) > 1{
			mean := NUMBER_OF_PROCESS/diff
			fmt.Println(mean)
			poisson := Poisson{Lambda: float64(mean)}
			variance = poisson.Variance()
			fmt.Printf("%f\n", poisson.StdDev())
			fmt.Printf("%f\n", variance)
		}

		fmt.Printf("Arrival Time : %d\n", float64(sec) + variance)

		sec = float64(sec) + variance
		process1 := newProcess(float64(sec) + variance)
		if i == 1{
			process1.arrivalTime = float64(sec)
		}
		if len(processList) > 1 {
			process1.estimatedBurstTime = calculateBurst(int64(len(processList) - 1))
		}
		r1 = rand.New(s1)
		priority := r1.Int63n(3)
		process1.priority = priority
		process1.actualBurstTime = burst
		processList = append(processList, process1)

		fmt.Printf("Process %d Estimated Burst Time : %f\n", i, process1.estimatedBurstTime)
	}
	fmt.Println(processList)

}