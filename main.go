package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := influxdb2.NewClient("http://localhost:8086", "influxdb:influxdb")
	writeApi := client.WriteApiBlocking("", "test/autogen")
	defer client.Close()

	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "mem"},
		map[string]interface{}{"num": 23, "max": uint(42), "avg_ns": 111111111},
		time.Now())
	// the above will produce:
	// stat,unit=mem avg_ns=111111111i,max=42u,num=23i 1593456269331238808\n
	// which fails to be written in the db with empty error and status 400
	// if I delete manually all i, u and \n then it is written successfully.
	err := writeApi.WritePoint(ctx, p)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// this is ok
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%d", 23.5, 45)
	err = writeApi.WriteRecord(context.Background(), line)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
