package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	URL string
	sem chan struct{}
}

func NewServer(url string) *Server {
	return &Server{
		URL: url,
		sem: make(chan struct{}, 3),
	}
}
func (s *Server) HandleReq(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sem <- struct{}{}
	defer func() {
		<-s.sem
	}()

	fmt.Printf("\nServer %s  is handling URL %s\n", s.URL, url)
	time.Sleep(2 * time.Second)
	openBrowser(url)

	fmt.Printf("server %s finished handling URL %s \n ", s.URL, url)

}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	if err := exec.Command(cmd, args...).Start(); err != nil {
		fmt.Printf("Failed to open URL or Browser %v\n", err)
	}
}

type Loadbalancer struct {
	servers []*Server
	index   int
	mu      sync.Mutex
}

func NewLoadBaoancer(servers []*Server) *Loadbalancer {
	return &Loadbalancer{servers: servers}
}
func (lb *Loadbalancer) GetNextServer() *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}

func main() {
	servers := []*Server{
		NewServer("server 1"),
		NewServer("server 2"),
		NewServer("server 3"),
	}

	lb := NewLoadBaoancer(servers)

	var wg sync.WaitGroup

	for {
		var url string
		fmt.Printf("Enter the URL (or Exit to quit requesting):")
		fmt.Scanln(&url)
		if url == "exit" {
			break
		}
		server := lb.GetNextServer()

		wg.Add(1)
		go server.HandleReq(url, &wg)
	}
	wg.Wait()
}
