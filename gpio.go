package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stianeikeland/go-rpio"
)

type response struct {
	Status interface{} `json:"status,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type pinStatus struct {
	Pin    uint8 `json:"pin"`
	Status bool  `json:"status"`
}

var (
	usablePins = map[uint8]uint8{ //  BCM: wPi(GPIO)
		4:  7,
		5:  21,
		6:  22,
		10: 99,
		13: 23,
		12: 26,
		16: 27,
		17: 0,
		18: 1,
		19: 24,
		22: 3,
		20: 28,
		21: 29,
		23: 4,
		24: 5,
		26: 25,
		27: 2}
)

// // Define an ordered map
type orderedMap []pinStatus

func (s orderedMap) Len() int {
	return len(s)
}
func (s orderedMap) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s orderedMap) Less(i, j int) bool {
	return s[i].Pin < s[j].Pin
}

// router.Path("/api/gpio/{id}").Methods("GET").HandlerFunc(getStatus).Methods("GET")
func getListPins() {
	fmt.Println("get usable pin list")
	var pins []int
	for k := range usablePins {
		pins = append(pins, int(k))
	}
	sort.Ints(pins)
	i, err := json.Marshal(pins)
	if err != nil {
		fmt.Println(fmt.Sprintf(`{error: %s}`, err))
	}
	fmt.Println(string(i))
}

// router.Path("/api/gpio/{id}").Methods("GET").HandlerFunc(getStatus).Methods("GET")
func getPins(w http.ResponseWriter, q *http.Request) {
	fmt.Println("getPins")
	var pins []int
	for k := range usablePins {
		pins = append(pins, int(k))
	}
	sort.Ints(pins)
	i, err := json.Marshal(pins)
	if err != nil {
		fmt.Println(fmt.Sprintf(`{error: %s}`, err))
	}
	fmt.Fprintf(w, string(i))
}

// router.Path("/api/gpio/{id}").Methods("GET").HandlerFunc(getStatus).Methods("GET")
func getStatus(w http.ResponseWriter, q *http.Request) {
	params := mux.Vars(q)
	// w http.ResponseWriter, body []byte, params map[string]string) {
	idStr := params["id"]
	if idStr == "" {
		sendResponse(w, response{Error: "ID not specified"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendResponse(w, response{Error: err.Error()})
		return
	}
	fmt.Printf("get status of %s", idStr)

	if _, ok := usablePins[uint8(id)]; ok {
		pin := rpio.Pin(id)
		// status := pin.Read()
		if pin.Read() == rpio.Low {
			sendResponse(w, response{Status: true})
		} else {
			sendResponse(w, response{Status: false})
		}
	} else {
		sendResponse(w, response{Error: "not usable"})
	}
}

// router.Path("/api/gpio").Methods("GET").HandlerFunc(auth.SecureFunc(getGPIOStatus)).Methods("GET")}
func getGPIOStatus(w http.ResponseWriter, q *http.Request) {
	fmt.Println("getGPIOStatus")
	var pins orderedMap
	for k := range usablePins {
		pin := rpio.Pin(k)
		st := false
		if pin.Read() == rpio.Low {
			st = true
		}
		pins = append(pins, pinStatus{Pin: k, Status: st})
	}
	sort.Sort(pins)
	i, err := json.Marshal(pins)
	if err != nil {
		fmt.Println(fmt.Sprintf(`{error: %s}`, err))
	}
	fmt.Fprintf(w, string(i))
}

// router.Path("/api/gpio/{id}/high").Methods("POST").HandlerFunc(high)
func high(w http.ResponseWriter, q *http.Request) {
	params := mux.Vars(q)
	// w http.ResponseWriter, body []byte, params map[string]string) {
	idStr := params["id"]
	if idStr == "" {
		sendResponse(w, response{Error: "ID not specified"})
		return
	}
	fmt.Printf("set HIGH to %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendResponse(w, response{Error: err.Error()})
		return
	}

	if _, ok := usablePins[uint8(id)]; ok {
		pin := rpio.Pin(id)
		pin.Write(rpio.High)
	} else {
		sendResponse(w, response{Error: "not usable"})
	}
}

// router.Path("/api/gpio/{id}/low").Methods("POST").HandlerFunc(low)
func low(w http.ResponseWriter, q *http.Request) {
	params := mux.Vars(q)
	// w http.ResponseWriter, body []byte, params map[string]string) {
	idStr := params["id"]
	if idStr == "" {
		sendResponse(w, response{Error: "ID not specified"})
		return
	}
	fmt.Printf("set LOW to %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendResponse(w, response{Error: err.Error()})
		return
	}

	if _, ok := usablePins[uint8(id)]; ok {
		pin := rpio.Pin(id)
		pin.Write(rpio.Low)
	} else {
		sendResponse(w, response{Error: "not usable"})
	}
}

// router.Path("/api/gpio/{id}/toggle").Methods("POST").HandlerFunc(toggle)
func toggle(w http.ResponseWriter, q *http.Request) {
	params := mux.Vars(q)
	// w http.ResponseWriter, body []byte, params map[string]string) {
	idStr := params["id"]
	if idStr == "" {
		sendResponse(w, response{Error: "ID not specified"})
		return
	}
	fmt.Printf("set TOGGLE to %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendResponse(w, response{Error: err.Error()})
		return
	}

	if _, ok := usablePins[uint8(id)]; ok {
		pin := rpio.Pin(id)
		pin.Toggle()
	} else {
		sendResponse(w, response{Error: "not usable"})
	}
}

func sendResponse(w http.ResponseWriter, r response) {
	i, err := json.Marshal(r)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf(`{error: %s}`, err))
	}
	fmt.Fprintf(w, string(i))
}

// func getRepos(w http.ResponseWriter, q *http.Request) {
//         user := ""
//         // opt := &github.RepositoryListOptions{Type: "owner,collaborator,organization_member", Sort: "updated", Direction: "desc"}
//         opt := &github.RepositoryListOptions{
//                 Sort:        "full_name", // created, updated, pushed, full_name. Default: full_name
//                 Direction:   "asc",       //  asc or desc
//                 ListOptions: github.ListOptions{PerPage: 100},
//         }
//
//         repos, _, err := client.Repositories.List(user, opt)
//         if err != nil {
//                 fmt.Println(err)
//         }
//         sendResponse(w, response{Status: "ok", Data: repos})
// }
