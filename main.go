package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var cachedPages = make(map[string]CachedPage)

type CachedPage struct {
	ResponseData Response
	Requested    time.Time
}

type Response struct {
	Action string  `json:"action"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Answer float64 `json:"answer"`
	Cached bool    `json:"cached"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple math api. Make requests using x and y as GET parameters to /add, /subtract, /multiply or /divide")
}

func returnCachedPage(r *http.Request) (CachedPage, bool) {
	response, valid := cachedPages[r.RequestURI]
	if valid == true && time.Since(response.Requested) < time.Duration(60)*time.Second {
		validCache := true
		return response, validCache
	}
	validCache := false
	return response, validCache
}

func addCachedPage(r *http.Request, action string, x, y, answer float64) {
	cachedPages[r.RequestURI] = CachedPage{ResponseData: Response{Action: action, X: x, Y: y, Answer: answer, Cached: true}, Requested: time.Now()}
}

func verifyParams(r *http.Request) (string, float64, float64) {
	xList, ok := r.URL.Query()["x"]
	if !ok || len(xList) < 1 {
		return "Param x is missing", 0, 0
	}
	x, err := strconv.ParseFloat(xList[0], 64)
	if err != nil {
		return "Param x is not a number", 0, 0
	}

	yList, ok := r.URL.Query()["y"]
	if !ok || len(yList) < 1 {
		return "Param y is missing", 0, 0
	}
	y, err := strconv.ParseFloat(yList[0], 64)
	if err != nil {
		return "Param y is not a number", 0, 0
	}
	return "ok", x, y
}

func add(w http.ResponseWriter, r *http.Request) {
	cachedPage, validCache := returnCachedPage(r)
	if !validCache {
		paramCheckMessage, x, y := verifyParams(r)
		if paramCheckMessage != "ok" {
			fmt.Fprintf(w, paramCheckMessage)
			return
		}
		answer := x + y
		addResponse := Response{Action: "add", X: x, Y: y, Answer: answer, Cached: false}
		addCachedPage(r, "add", x, y, answer)
		json.NewEncoder(w).Encode(addResponse)
		return
	}
	addResponse := cachedPage.ResponseData
	json.NewEncoder(w).Encode(addResponse)
	return
}

func subtract(w http.ResponseWriter, r *http.Request) {
	cachedPage, validCache := returnCachedPage(r)
	if !validCache {
		paramCheckMessage, x, y := verifyParams(r)
		if paramCheckMessage != "ok" {
			fmt.Fprintf(w, paramCheckMessage)
			return
		}
		answer := x - y
		addResponse := Response{Action: "subtract", X: x, Y: y, Answer: answer, Cached: false}
		addCachedPage(r, "subtract", x, y, answer)
		json.NewEncoder(w).Encode(addResponse)
		return
	}
	addResponse := cachedPage.ResponseData
	json.NewEncoder(w).Encode(addResponse)
	return
}

func multiply(w http.ResponseWriter, r *http.Request) {
	cachedPage, validCache := returnCachedPage(r)
	if !validCache {
		paramCheckMessage, x, y := verifyParams(r)
		if paramCheckMessage != "ok" {
			fmt.Fprintf(w, paramCheckMessage)
			return
		}
		answer := x * y
		addResponse := Response{Action: "multiply", X: x, Y: y, Answer: answer, Cached: false}
		addCachedPage(r, "multiply", x, y, answer)
		json.NewEncoder(w).Encode(addResponse)
		return
	}
	addResponse := cachedPage.ResponseData
	json.NewEncoder(w).Encode(addResponse)
	return
}

func divide(w http.ResponseWriter, r *http.Request) {
	cachedPage, validCache := returnCachedPage(r)
	if !validCache {
		paramCheckMessage, x, y := verifyParams(r)
		if paramCheckMessage != "ok" {
			fmt.Fprintf(w, paramCheckMessage)
			return
		}
		if y == 0 {
			fmt.Fprint(w, "Cannot divide by zero")
			return
		}
		answer := x / y
		addResponse := Response{Action: "divide", X: x, Y: y, Answer: answer, Cached: false}
		addCachedPage(r, "divide", x, y, answer)
		json.NewEncoder(w).Encode(addResponse)
		return
	}
	addResponse := cachedPage.ResponseData
	json.NewEncoder(w).Encode(addResponse)
	return
}

func handleRequests() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/add", add)
	http.HandleFunc("/subtract", subtract)
	http.HandleFunc("/multiply", multiply)
	http.HandleFunc("/divide", divide)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
