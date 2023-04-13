package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Queue struct {
	lock    sync.Mutex
	hashmap map[string][]int
}

// var queue Queue

var queue = Queue{}

func init() {
	queue.hashmap = map[string][]int{}
}

func (qu *Queue) Push(name string, list []int) {
	qu.lock.Lock()
	defer qu.lock.Unlock()

	tmp, _ := qu.hashmap[name]

	tmp = append(tmp, list...)
	qu.hashmap[name] = tmp
}

func (qu *Queue) Pop(name string, w http.ResponseWriter) {
	qu.lock.Lock()
	defer qu.lock.Unlock()

	tmp, ok := qu.hashmap[name]
	if !ok || len(tmp) == 0 {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "null")
		return
	}
	data := tmp[len(tmp)-1]
	tmp = tmp[:len(tmp)-1]
	qu.hashmap[name] = tmp
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, data)
}

func (qu *Queue) BPop(name string, duration int, w http.ResponseWriter) {
	// qu.lock.Lock()
	// defer qu.lock.Unlock()

	tmp, ok := qu.hashmap[name]
	if !ok || len(tmp) == 0 {
		//Wait for given duration for push operation
		for i := 0; i < duration; i++ {
			time.Sleep(time.Second)
			qu.lock.Lock()
			tmp1, _ := qu.hashmap[name]
			if len(tmp1) > 0 {
				defer qu.lock.Unlock()
				data := tmp1[len(tmp1)-1]
				tmp1 = tmp1[:len(tmp1)-1]
				qu.hashmap[name] = tmp1
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, data)
				return
			}
			qu.lock.Unlock()
		}
		// if still no push operation then return null
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "null")
		return
	}
	data := tmp[len(tmp)-1]
	tmp = tmp[:len(tmp)-1]
	qu.hashmap[name] = tmp
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, data)
}

func QPush(cmd []string, w http.ResponseWriter) {

	newQueue := []int{}

	newList := cmd[2:]
	for _, v := range newList {
		v, _ := strconv.ParseInt(v, 0, 64)
		newQueue = append(newQueue, int(v))
	}

	queue.Push(cmd[1], newQueue)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "ok")

}

func QPop(cmd []string, w http.ResponseWriter) {
	queue.Pop(cmd[1], w)
}

func BQPop(cmd []string, w http.ResponseWriter) {
	duration, _ := strconv.ParseInt(cmd[2], 0, 64)
	queue.BPop(cmd[1], int(duration), w)
}
