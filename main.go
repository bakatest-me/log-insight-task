package main

import (
	"log"
	"log-insign-task/src/insight"
	"log-insign-task/util"
)

func main() {
	filename := "pprof/heap.pprof"
	prof := util.NewProfiling(filename)
	defer prof.Close()

	is := insight.NewLogService()
	resp, err := is.Run("logfile.txt")
	if err != nil {
		log.Println(err)
		return
	}
	is.Print(resp)
}
