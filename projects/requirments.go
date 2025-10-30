package main

/*
import (

	 "time"
	 "testing"
	 "sync"
	 "strings"
     "net/http"
	 "log"
	 "io"
	 "context"
	 "database"

	 "framework - gin"
	 "grpc api"
	 "kafka"
	 "redis"
	 "mysql"

)
*/

import (
	"fmt"
	"time"
	"testing"
)

func times() {

	now := time.Now()
	fmt.Println("currrent time", now)
	fmt.Println("day", now.Day())
	fmt.Println("date only", now.Format("2006-01-02 15:04:05"))

	timers:=time.NewTimer(3* time.Second)
   fmt.Println("timers")
   fmt.Println(<-timers.C)
}

func carboys( value int) int{

	 return value - 500000
}

func testcars(t *testing.T){

	rest:=carboys(1000000)
	exected:=500000

	if rest!=exected{
		t.Errorf("test case failed %d %d",rest,exected)
	}
}


func main() {
	times()
}
