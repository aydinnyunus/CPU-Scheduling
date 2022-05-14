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
	ShutdownTime = 8 * 60.
	ALPHA        = float64(0.4)
	HIGH = 1
	NORMAL = 2
	LOW = 3
)

var processList []*Process
var queueList []*Queues
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

var totalTimeCounted = float64(0)
var waitTime = float64(0)
var turnaroundTime = float64(0)
var totalTime = float64(0)

// Queues the Queues is a Passive Object representing resource
type Queues struct {
	id int
	max int
	priority int
	availableQueues int
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
	queues.availableQueues++
	if queues.availableQueues == queues.max {
		counterSwt.Set(false)
	}
}

func (queues *Queues) Release() {
	queues.availableQueues--
	counterSwt.Set(true)
}

// Process the Process is a Runner
type Process struct {
	*godes.Runner
	id                                                                                                                                                     int
	exitTime, actualBurstTime, estimatedBurstTime, arrivalTime, remainingTime, serviceTime, waitTime, turnAroundTime, avgArrivalTime, avgWaitTime, avgTurnAroundTime float64
	isCalculated																																		   bool
}

func (process *Process) Run() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	no := service.Get(1. / r1.Float64())
	process.serviceTime = no
	a0 := godes.GetSystemTime()
	min := queueList[0]
	max := queueList[0]
	middle := queueList[0]
	setPriorities()
	for _, value := range queueList {
		if 	value.priority< min.priority {
			min = value
		}
		if value.priority > max.priority {
			max = value
		}
	}

	for i:=0;i<3;i++{
		if middle.id == max.id{
			continue
		}else if middle.id == min.id{
			middle = queueList[i+1]
			continue
		}
	}
	max.Catch(process)
	middle.Catch(process)
	min.Catch(process)


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
	for i := range processList{
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

func calculateAvgArr() float64 {
	arr := float64(0)
	for i := range processList{
		arr += processList[i].arrivalTime
	}
	return arr / float64(len(processList))

}

func calculateAvgWait() float64{

	return waitTime / (float64)(len(processList))

}

func calculateAvgTurn() float64 {

	return turnaroundTime / (float64)(len(processList))

}

func calculateAvgQueue() {

}

func findMinAndMaxAvailable() (min *Queues, max *Queues) {
	min = queueList[0]
	max = queueList[0]
	for _, value := range queueList {
		if 	value.availableQueues< min.availableQueues {
			min.availableQueues = value.availableQueues
		}
		if value.availableQueues > max.availableQueues {
			max.availableQueues = value.availableQueues
		}
	}
	return min, max
}

func setPriorities (){
	min, max := findMinAndMaxAvailable()
	for i := range queueList{
		if queueList[i] == max{
			max.priority = HIGH
		} else if queueList[i] == min{
			min.priority = NORMAL
		} else {
			queueList[i].priority = LOW
		}
	}
}
func roundRobin(){
	timeQuantum := float64(5)
	totalTime = calculateTotalTime()
	fmt.Println("Total Time ", totalTime)
	for math.Floor(totalTime*10000) / 10000 != 0{
		for i := range processList {
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
	for i := range processList{
		totalTime += processList[i].actualBurstTime
	}
	return totalTime
}

func main() {
	measures = [][]float64{}
	for i:=0;i<3;i++{
		tellers = &Queues{i,3,0,0}
		queueList = append(queueList, tellers)
	}
	godes.Run()
	counterSwt.Set(true)
	count := 0
	for {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		no := arrival.Get(1. / r1.Float64())
		burst := r1.Float64()*100
		customer := &Process{&godes.Runner{}, count, burst, 0, no, burst, 0, 0, 0, 0,0,0,0,false}
		processArrivalQueue.Place(customer)
		processList = append(processList, customer)
		if count > 1 {
			customer.estimatedBurstTime = calculateBurst(customer)
		}
		godes.AddRunner(customer)
		godes.Advance(no)
		if godes.GetSystemTime() > ShutdownTime {
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
