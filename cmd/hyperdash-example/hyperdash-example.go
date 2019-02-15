package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/lgarithm/hyperdash-sdk-go/hyperdash"
)

var (
	name = flag.String("n", "example-experiment", "name")
)

func main() {
	flag.Parse()
	hdc, err := hyperdash.NewClient(*name)
	if err != nil {
		log.Print(err)
		return
	}
	defer hdc.Close()
	if err := hdc.Start(); err != nil {
		log.Print(err)
		return
	}
	hdc.Param(map[string]interface{}{
		"x":     1,
		"label": `string`,
	})
	t0 := time.Now()
	tk := time.NewTicker(2 * time.Second)
	defer tk.Stop()
	for t := range tk.C {
		msg := fmt.Sprintf("%s\n", t)
		log.Printf("%s", msg)
		hdc.Log(msg)
		hdc.Metric("rate", t.Sub(t0).Seconds())
	}
}
