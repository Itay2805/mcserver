package main

import (
	"flag"
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Profiling related code
	//
	var f *os.File
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}

		log.Println("CPU Profiler has started")
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			if *memprofile != "" {
				f, err := os.Create(*memprofile)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				defer f.Close()
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
			}

			pprof.StopCPUProfile()

			f.Close()
			os.Exit(0)
		}
	}()
	//
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	log.Println("Everything is ready to start")

	// start the game loop
	go game.StartGameLoop()

	// start the server
	server.StartServer()
}

