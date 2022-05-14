package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/agoussia/godes"
)

//Input Parameters
const (
	SHUTDOWN_TIME = 8 * 60.
	ALPHA         = float64(0.4)
)

// the arrival and service are two random number generators for the exponential  distribution
var arrival = godes.NewExpDistr(true)
var service = godes.NewExpDistr(true)

// true when any counter is available
var counterSwt = godes.NewBooleanControl()

// FIFO Queue for the arrived customers
var processArrivalQueue = godes.NewFIFOQueue("0")
var allProcessQueue = godes.NewFIFOQueue("0")

var tellers *Queues
var measures [][]float64
var titles = []string{
	"Elapsed Time",
	"Queue Length",
	"Queueing Time",
	"Service Time",
}

var availableQueues int = 0

// the Queues is a Passive Object represebting resource
type Queues struct {
	max int
}

func (queues *Queues) Catch(customer *Process) {
	for {
		counterSwt.Wait(true)
		if processArrivalQueue.GetHead().(*Process).id == customer.id {
			break
		} else {
			godes.Yield()
		}
	}
	availableQueues++
	if availableQueues == queues.max {
		counterSwt.Set(false)
	}
}

func (queues *Queues) Release() {
	availableQueues--
	counterSwt.Set(true)
}

// Process the Process is a Runner
type Process struct {
	*godes.Runner
	id                                                                                                                                                     int
	actualBurstTime, estimatedBurstTime, arrivalTime, remainingTime, serviceTime, waitTime, turnAroundTime, avgArrivalTime, avgWaitTime, avgTurnAroundTime float64
	priority                                                                                                                                               int64
}

func (process *Process) Run() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	no := service.Get(1. / r1.Float64())
	process.serviceTime = float64(time.Now().Unix()) + no
	a0 := godes.GetSystemTime()
	tellers.Catch(process)
	a1 := godes.GetSystemTime()
	processArrivalQueue.Get()
	qlength := float64(processArrivalQueue.Len())
	godes.Advance(no)
	a2 := godes.GetSystemTime()
	tellers.Release()
	collectionArray := []float64{a2 - a0, qlength, a1 - a0, a2 - a1}
	measures = append(measures, collectionArray)
	fmt.Printf("Estimated Burst Time %f", process.estimatedBurstTime)
	fmt.Println()
	fmt.Printf("Arrival Time %f", process.arrivalTime)
	fmt.Println()

}

func getProcessByID() {

}

func calculateBurst(process *Process) float64 {
	calculate := ALPHA*process.actualBurstTime + (1-ALPHA)*process.estimatedBurstTime
	return calculate
}

func calculateAvgArr(process *Process) {
	//process.avgArrivalTime = process.arrivalTime / NUMBER_OF_PROCESS

}

func calculateAvgWait(process *Process) {
	process.waitTime = process.serviceTime - process.arrivalTime

	//process.avgWaitTime = process.waitTime / NUMBER_OF_PROCESS

	//https://www.gatevidyalay.com/round-robin-round-robin-scheduling-examples/

	//Waiting time = Turn Around time – Burst time

}

func calculateAvgTurn(process *Process) {

	process.turnAroundTime = process.waitTime + process.actualBurstTime

	//process.avgTurnAroundTime = process.turnAroundTime / NUMBER_OF_PROCESS

}

func calculateAvgQueue() {

}

func main() {
	measures = [][]float64{}
	tellers = &Queues{3}
	godes.Run()
	counterSwt.Set(true)
	count := 0
	for {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		no := arrival.Get(1. / r1.Float64())
		customer := &Process{&godes.Runner{}, count, 0, 0, float64(time.Now().Unix()) + no, 0, 0, 0, 0, 0}
		processArrivalQueue.Place(customer)
		allProcessQueue.Place(customer)
		if count > 1 {
			customer.estimatedBurstTime = calculateBurst(customer)
		}
		godes.AddRunner(customer)
		godes.Advance(no)
		if godes.GetSystemTime() > SHUTDOWN_TIME {
			break
		}
		count++
	}
	godes.WaitUntilDone() // waits for all the runners to finish the Run()
	collector := godes.NewStatCollector(titles, measures)
	collector.PrintStat()
	fmt.Printf("Finished \n")
}

/* OUTPUT
Variable		#	Average	Std Dev	L-Bound	U-Bound	Minimum	Maximum
Elapsed Time	944	 2.591	 1.959	 2.466	 2.716	 0.005	11.189
Queue Length	944	 2.411	 3.069	 2.215	 2.607	 0.000	13.000
Queueing Time	944	 1.293	 1.533	 1.195	 1.391	 0.000	 6.994
Service Time	944	 1.298	 1.247	 1.219	 1.378	 0.003	 7.824
*/
