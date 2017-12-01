package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var numOfRows int

type Vertex struct {
	ID                 int
	leader             *Vertex
	followers          []*Vertex
	clusterMaxDistance int
}

func (v Vertex) String() string {

	var followers []int

	for _, x := range v.followers {
		followers = append(followers, x.ID)
	}

	return fmt.Sprintf("\nID:\t%v\nLeader:\t%v\nFollower:\t%v\nclusterMaxDistance:\t%v\n\n\n",
		v.ID, v.leader.ID, followers, v.clusterMaxDistance)
}

var LeaderMap = make(map[int]*Vertex, numOfRows)

// ConsumeFollowers is a union-like interface for vertices
func (v *Vertex) ConsumeFollowers(w *Vertex) {

	w.leader = v

	if v.clusterMaxDistance < w.clusterMaxDistance {
		v.clusterMaxDistance = w.clusterMaxDistance
	}

	for _, x := range w.followers {
		x.leader = v
		x.clusterMaxDistance = v.clusterMaxDistance
		v.followers = append(v.followers, x)
	}

	w.followers = nil

	delete(LeaderMap, w.ID)
}

type Edge struct {
	Vertex1 *Vertex
	Vertex2 *Vertex
	Cost    int
	Index   int
}

func (e Edge) String() string {
	return fmt.Sprintf("\nVertex1:\t%v\nVertex2:\t%v\nCost:\t%v\n\n\n",
		e.Vertex1.ID, e.Vertex2.ID, e.Cost)
}

// A EdgeHeap returns the Edge with the lowest Cost
type EdgeHeap []*Edge

func (eh EdgeHeap) Len() int { return len(eh) }

func (eh EdgeHeap) Less(i, j int) bool {
	return eh[i].Cost < eh[j].Cost
}

func (eh EdgeHeap) Swap(i, j int) {
	eh[i], eh[j] = eh[j], eh[i]
	eh[i].Index = i
	eh[j].Index = j
}

// Push adds Edges to EdgeHeaps
func (eh *EdgeHeap) Push(x interface{}) {
	n := len(*eh)
	e := x.(*Edge)
	e.Index = n
	*eh = append(*eh, e)
}

// Pop returns the Edge with the lowest value of Cost
func (eh *EdgeHeap) Pop() interface{} {
	old := *eh
	n := len(old)
	v := old[n-1]
	v.Index = -1 // for safety, identify it's no longer in heap
	*eh = old[0 : n-1]
	return v
}

var eh EdgeHeap

func main() {

	readFile(os.Args[1])

	kCluster(3)

	maxDistance := 0

	for _, v := range LeaderMap {
		if v.clusterMaxDistance > maxDistance {
			maxDistance = v.clusterMaxDistance
		}
	}

	fmt.Println(maxDistance)

}

func readFile(filename string) {
	i := 0

	file, err := os.Open(filename) //should read in file named in CLI

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Scan first line
	if scanner.Scan() {

		firstLine := strings.Fields(scanner.Text())

		numOfRows, err = strconv.Atoi(firstLine[0])

		if err != nil {
			log.Fatalf("couldn't convert number: %v\n", err)
		}

	}

	for scanner.Scan() {

		//remove spaces
		thisLine := strings.Fields(scanner.Text())

		firstVertex, err := strconv.Atoi(thisLine[0])
		secondVertex, err := strconv.Atoi(thisLine[1])
		edgeCost, err := strconv.Atoi(thisLine[2])

		if err != nil {
			log.Fatal(err)
		}

		w, ok := LeaderMap[firstVertex]

		if !ok {
			w = &Vertex{firstVertex, nil, []*Vertex{}, 0}
			w.leader = w
			w.followers = append(w.followers, w)
			LeaderMap[firstVertex] = w
		}

		u, ok := LeaderMap[secondVertex]

		if !ok {
			u = &Vertex{secondVertex, nil, []*Vertex{}, 0}
			u.leader = u
			u.followers = append(u.followers, u)
			LeaderMap[secondVertex] = u
		}

		e := &Edge{w, u, edgeCost, i}

		eh = append(eh, e)

		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	heap.Init(&eh)

}

func kCluster(s int) {

	for len(LeaderMap) > s {

		e := heap.Pop(&eh).(*Edge)

		x := e.Vertex1
		y := e.Vertex2

		if e.Cost > x.leader.clusterMaxDistance {
			x.leader.clusterMaxDistance = e.Cost
		}

		if e.Cost > y.clusterMaxDistance {
			y.leader.clusterMaxDistance = e.Cost
		}

		if x.leader != y.leader {
			if len(x.leader.followers) > len(y.leader.followers) {
				x.leader.ConsumeFollowers(y.leader)
			} else {
				y.leader.ConsumeFollowers(x.leader)
			}
		}
	}
}
