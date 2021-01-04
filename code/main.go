package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	fcRequestIDHeader = "x-fc-request-id"
	fcControlPath     = "x-fc-control-path"

	fcLogInvokeStartPrefix     = "FC Invoke Start RequestId: %s"
	fcLogInvokeEndPrefix       = "FC Invoke End RequestId: %s"
	fcLogInitializeStartPrefix = "FC Initialize Start RequestId: %s"
	fcLogInitializeEndPrefix   = "FC Initialize End RequestId: %s"
	fcLogPreFreezeStartPrefix  = "FC PreFreeze Start RequestId: %s"
	fcLogPreFreezeEndPrefix    = "FC PreFreeze End RequestId: %s"
	fcLogPreStopStartPrefix    = "FC PreStop Start RequestId: %s"
	fcLogPreStopEndPrefix      = "FC PreStop End RequestId: %s"
)

// init 与容器生命周期绑定
// 我们用 custom runtime 嘛，不用 custom container 嘛，还要打镜像
func main() {
	fmt.Println("FunctionCompute golang runtime inited.")

	http.HandleFunc("/", handle)
	port := os.Getenv("FC_SERVER_PORT")
	if port == "" {
		port = "9000"
	}
	http.ListenAndServe(":"+port, nil)
}

func handle(w http.ResponseWriter, req *http.Request) {
	controlPath := req.Header.Get(fcControlPath)
	fmt.Println("controlPath", controlPath)
	if controlPath == "/initialize" {
		initializeHandler(w, req)
	} else if controlPath == "/pre-freeze" {
		preFreezeHandler(w, req)
	} else if controlPath == "/pre-stop" {
		preStopHandler(w, req)
	} else {
		invokeHandler(w, req)
	}
}

var app *newrelic.Application

func initializeHandler(w http.ResponseWriter, req *http.Request) {
	requestID := req.Header.Get(fcRequestIDHeader)
	fmt.Println(fmt.Sprintf(fcLogInitializeStartPrefix, requestID))
	defer func() {
		fmt.Println(fmt.Sprintf(fcLogInitializeEndPrefix, requestID))
	}()

	if application, err := newrelic.NewApplication(
		newrelic.ConfigAppName("newrelic-demo"),
		newrelic.ConfigLicense("97dc68c3ea1926da0d6988fbce83e9652193NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true)); err != nil {
		fmt.Println("An error occured when new application")
		panic(err)
		app = application
	}

	w.Write([]byte(""))
}

func preFreezeHandler(w http.ResponseWriter, req *http.Request) {
	requestID := req.Header.Get(fcRequestIDHeader)
	fmt.Println(fmt.Sprintf(fcLogPreFreezeStartPrefix, requestID))
	defer func() {
		fmt.Println(fmt.Sprintf(fcLogPreFreezeEndPrefix, requestID))
	}()
	time.Sleep(2 * time.Second)
	w.Write([]byte(""))
}

func preStopHandler(w http.ResponseWriter, req *http.Request) {
	requestID := req.Header.Get(fcRequestIDHeader)
	fmt.Println(fmt.Sprintf(fcLogPreStopStartPrefix, requestID))
	defer func() {
		fmt.Println(fmt.Sprintf(fcLogPreStopEndPrefix, requestID))
	}()
	w.Write([]byte(""))
}

func invokeHandler(w http.ResponseWriter, req *http.Request) {
	txn := app.StartTransaction("invoke")
	defer txn.End()
	requestID := req.Header.Get(fcRequestIDHeader)
	fmt.Println(fmt.Sprintf(fcLogInvokeStartPrefix, requestID))
	defer func() {
		fmt.Println(fmt.Sprintf(fcLogInvokeEndPrefix, requestID))
	}()

	headerComponent := txn.StartSegment("header")
	n := rand.Intn(20) // n will be between 0 and 10
	fmt.Printf("Sleeping %d ms...\n", 20+n)
	time.Sleep(time.Duration(20+n) * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	headerComponent.End()

	bodyComponent := txn.StartSegment("header")
	n = rand.Intn(30) // n will be between 0 and 10
	fmt.Printf("Sleeping %d ms...\n", 30+n)
	time.Sleep(time.Duration(30+n) * time.Millisecond)
	bodyComponent.End()
	w.Write([]byte(fmt.Sprintf("Hello, golang http invoke!")))
}
