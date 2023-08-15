package tools

import (
	_ "embed"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func DetectService() {
	hostname := flag.String("hostname", "", "hostname to test")
	startPort := flag.Int("start-port", 6379, "the port on which the scanning starts")
	endPort := flag.Int("end-port", 6379, "the port from which the scanning ends")
	timeout := flag.Duration("timeout", time.Millisecond*200, "timeout")
	flag.Parse()

	ports := []int{}

	wg := &sync.WaitGroup{}
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			opened := isAlive(*hostname, p, *timeout)
			if opened {
				ports = append(ports, p)
			}
			wg.Done()
		}(port)
	}

	wg.Wait()
	fmt.Printf("opened ports: %v\n", ports)
}

func isAlive(host string, port int, timeout time.Duration) bool {
	time.Sleep(time.Millisecond * 1)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err == nil {
		_ = conn.Close()
		return true
	}

	return false
}
