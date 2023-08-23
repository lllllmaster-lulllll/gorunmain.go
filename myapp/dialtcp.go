package myapp

import (
	_ "embed"
	"fmt"
	"net"
	"time"
)

func DetectService(hostname string, port int, timeout time.Duration) (b bool, err error) {
	// hostname := flag.String("hostname", "", "hostname to test")
	// startPort := flag.Int("start-port", 6379, "the port on which the scanning starts")
	// endPort := flag.Int("end-port", 6379, "the port from which the scanning ends")
	// timeout := flag.Duration("timeout", time.Millisecond*200, "timeout")
	// flag.Parse()

	// ports := []int{}
	var opened bool
	// wg := &sync.WaitGroup{}
	// for port := *startPort; port <= *endPort; port++ {
	// wg.Add(1)
	// go func(p int) {
	opened, err = isAlive(hostname, port, timeout)
	if err != nil {
		return false, err
	}
	// if opened {
	// 	ports = append(ports, p)
	// }
	// wg.Done()
	// }(port)

	// wg.Wait()
	// fmt.Printf("opened ports: %v\n", ports)
	if !opened {
		return false, nil
	}
	return true, nil
}

func isAlive(host string, port int, timeout time.Duration) (b bool, err error) {
	time.Sleep(time.Millisecond * 1)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err == nil {
		_ = conn.Close()
		return true, nil
	}

	return false, err
}
