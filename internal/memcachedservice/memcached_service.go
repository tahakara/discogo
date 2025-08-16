package memcachedservice

import (
	"log"
	"net"

	"discogo/internal/memcached"
)

type MemcachedService struct {
	ServerIP        string
	MemcachedClient *memcached.MemcachedClient
}

func NewMemcachedService(serverIP string, memcachedClient *memcached.MemcachedClient) *MemcachedService {
	return &MemcachedService{
		ServerIP:        serverIP,
		MemcachedClient: memcachedClient,
	}
}

func (ms *MemcachedService) DiscoverService() error {
	if ms.ServerIP == "" {
		log.Println("Server IP address is not set.")
		return nil
	}

	conn, err := net.Dial("tcp", ms.ServerIP)
	if err != nil {
		log.Printf("Failed to connect to service at %s: %v", ms.ServerIP, err)
		return err
	}
	defer conn.Close()

	log.Printf("Successfully connected to service at %s", ms.ServerIP)
	return nil
}

func (ms *MemcachedService) GetServerIP() string {
	return ms.ServerIP
}
