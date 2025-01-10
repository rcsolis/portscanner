package main

import (
	"flag"
	"log"
	"os"
	"sync"

	scanner "github.com/rcsolis/portscanner/internal/scanner"
)

var address string
var initialPort int
var finalPort int

const TIMEOUT_DEFAULT = 3
const NUM_GOROUTINES = 10
const CHANNEL_BUFFER = 10

func connect(address string, ports []int, scanWG *sync.WaitGroup, scanCH chan int) {
	defer scanWG.Done()
	scanner.Scan(address, ports, TIMEOUT_DEFAULT, scanCH)
}

func main() {
	var portsToScan []int
	var openPorts []int
	var scanCH = make(chan int, CHANNEL_BUFFER)
	var scanWG sync.WaitGroup

	var collectWG sync.WaitGroup
	var collectMT sync.Mutex
	// Parse flags and arguments
	flag.IntVar(&initialPort, "from", 1, "Initial port to scan")
	flag.IntVar(&initialPort, "f", 1, "Initial port to scan (shorthand)")
	flag.IntVar(&finalPort, "to", 1024, "Final port to scan")
	flag.IntVar(&finalPort, "t", 1024, "Final port to scan (shorthand)")
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatal("Usage: portscanner <address>")
	}
	address = os.Args[len(os.Args)-1]
	if address == "" {
		address = "scanme.nmap.org"
	}
	if initialPort > finalPort {
		initialPort = 1
		finalPort = 1024
	}
	log.Println("Scanning from", initialPort, "to", finalPort, "on", address)
	// Generate the list of ports to scan
	portsToScan = scanner.GeneratePorts(initialPort, finalPort)
	log.Println("Scanning ", len(portsToScan), "ports")
	// Calculate the number of ports to scan per goroutine
	portsPerGoRoutine := (len(portsToScan) + NUM_GOROUTINES) / NUM_GOROUTINES
	log.Println("Scanning ", portsPerGoRoutine, "per goroutine")
	// Generates the goroutines to scan the ports
	for i := 0; i < NUM_GOROUTINES; i++ {
		startIndex := i * portsPerGoRoutine
		endIndex := (i + 1) * portsPerGoRoutine
		if endIndex > len(portsToScan) {
			endIndex = len(portsToScan)
		}
		scanWG.Add(1)
		go connect(address, portsToScan[startIndex:endIndex], &scanWG, scanCH)
	}
	// Generates the goroutine to collect the results
	collectWG.Add(1)
	go func() {
		defer collectWG.Done()
		for port := range scanCH {
			if port != 0 {
				collectMT.Lock()
				openPorts = append(openPorts, port)
				collectMT.Unlock()
			}
		}
	}()

	scanWG.Wait()
	close(scanCH)
	collectWG.Wait()

	log.Println("Open ports:", openPorts)
}
