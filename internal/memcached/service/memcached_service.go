package memcachedservice

import (
	lg "discogo/internal/logger"
	lgm "discogo/internal/logger/messages"
	"net"
	// "discogo/internal/memcached/memcachedclient"
)

type MemcachedService struct {
	ServerIP string
	// MemcachedClient *memcachedclient.MemcachedClient
}

func NewMemcachedService(serverIP string) *MemcachedService {
	return &MemcachedService{
		ServerIP: serverIP,
	}
}

func (ms *MemcachedService) DiscoverService() error {
	if ms.ServerIP == "" {
		lg.Error(lgm.MessageR(lgm.ErrorMemcachedFailedToRetrieveServerAddress, ms.ServerIP))
		return nil
	}

	conn, err := net.Dial("tcp", ms.ServerIP)
	if err != nil {
		lg.Error(lgm.MessageR(lgm.ErrorMemcachedConnectionFailed, err))
		return err
	}
	defer conn.Close()

	lg.Info(lgm.MessageR(lgm.InfoMemcachedKeyNotFound, ms.ServerIP))
	return nil
}

func (ms *MemcachedService) GetServerIP() string {
	return ms.ServerIP
}
