package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/agoussia/godes"
)

//Input Parameters
const (
	SHUTDOWN_TIME = 8 * 60.
	ALPHA         = float64(0.4)
	REALTIME = 0
	HIGH = 1
	NORMAL = 2
	LOW = 3
)

var processList []*Process
// the arrival and service are two random number generators for the exponential  distribution
var arrival = godes.NewExpDistr(true)
var service = godes.NewExpDistr(true)
var exit	= godes.NewExpDistr(true)

// true when any counter is available
var counterSwt = godes.NewBooleanControl()

// FIFO Queue for the arrived customers
var processArrivalQueue = godes.NewFIFOQueue("0")

var tellers *Queues
var measures [][]float64
var titles = []string{
	"Elapsed Time",
	"Queue Length",
	"Queueing Time",
	"Service Time",
}

var availableQueues = 0
var totalTimeCounted = float64(0)
var waitTime = float64(0)
var turnaroundTime = float64(0)
var totalTime = float64(0)

// Queues the Queues is a Passive Object represebting resource
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
	exitTime, actualBurstTime, estimatedBurstTime, arrivalTime, remainingTime, serviceTime, waitTime, turnAroundTime, avgArrivalTime, avgWaitTime, avgTurnAroundTime float64
	priority                                                                                                                                               int64
	isCalculated																																		   bool
}

func (process *Process) Run() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	no := service.Get(1. / r1.Float64())
	process.serviceTime = no
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

func getProcessByID(id int) *Process{
	for i, _ := range processList{
		if processList[i].id == id{
			return processList[i]
		}
	}
	return nil
}

func calculateBurst(process *Process) float64 {
	id := getProcessByID(process.id - 1)
	calculate := ALPHA*id.actualBurstTime + (1-ALPHA)*id.estimatedBurstTime
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

func roundRobin(){
	timeQuantum := float64(5)
	totalTime = calculateTotalTime()
	fmt.Println("Total Time ", totalTime)
	for math.Floor(totalTime*10000) / 10000 != 0{
		for i, _ := range processList {
			if processList[i].remainingTime <= timeQuantum && processList[i].remainingTime > 0{
				totalTimeCounted += processList[i].remainingTime
				totalTime -= processList[i].remainingTime

				processList[i].remainingTime = 0
				//fmt.Println(totalTime)
				//fmt.Println("process id is finished ", processList[i].id)
			} else if processList[i].remainingTime > 0 {

				processList[i].remainingTime -= timeQuantum
				totalTime -= timeQuantum
				totalTimeCounted += timeQuantum
				//fmt.Println(totalTime)

			}
			if processList[i].remainingTime == 0 && !processList[i].isCalculated {
				processList[i].waitTime = totalTimeCounted - processList[i].arrivalTime - processList[i].actualBurstTime
				waitTime += processList[i].waitTime

				processList[i].turnAroundTime = totalTimeCounted - processList[i].arrivalTime
				turnaroundTime += processList[i].turnAroundTime

				processList[i].exitTime = processList[i].arrivalTime + processList[i].turnAroundTime
				processList[i].isCalculated = true
				//fmt.Println("process id is calculated ", processList[i].id)
			}
		}
	}
	//fmt.Println(wait_time)
}

func calculateTotalTime() float64 {
	for i,_ := range processList{
		totalTime += processList[i].actualBurstTime
	}
	return totalTime
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
		burst := r1.Float64()*100
		customer := &Process{&godes.Runner{}, count, burst, 0, no, burst, 0, 0, 0, 0,0,0,0,0,false}
		processArrivalQueue.Place(customer)
		processList = append(processList, customer)
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
	roundRobin()
	fmt.Println("Wait Time", waitTime/float64(len(processList)))
	fmt.Println("TurnAround Time", turnaroundTime/float64(len(processList)))

	fmt.Printf("Finished \n")
}

/* OUTPUT
Variable		#	Average	Std Dev	L-Bound	U-Bound	Minimum	Maximum
Elapsed Time	944	 2.591	 1.959	 2.466	 2.716	 0.005	11.189
Queue Length	944	 2.411	 3.069	 2.215	 2.607	 0.000	13.000
Queueing Time	944	 1.293	 1.533	 1.195	 1.391	 0.000	 6.994
Service Time	944	 1.298	 1.247	 1.219	 1.378	 0.003	 7.824
*/
