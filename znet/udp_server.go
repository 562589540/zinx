package znet

import (
	"fmt"
	"github.com/562589540/zinx/zconf"
	"github.com/562589540/zinx/zlog"
	"net"
)

func (s *Server) ListenUDPConn() {
	// udp拓展 udp是无链接的
	uDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.IP, s.UdpPort))

	if err != nil {
		zlog.Ins().ErrorF("[START] resolve UDP addr err: %v\n", err)
		return
	}

	uDPConn, err := net.ListenUDP("udp", uDPAddr)

	if err != nil {
		zlog.Ins().ErrorF("[START] listen UDP err: %v\n", err)
		return
	}

	zlog.Ins().InfoF("[START] UDP server listening at IP: %s, Port %d", s.IP, s.UdpPort)

	go func() {
		for {
			buffer := make([]byte, zconf.GlobalObject.IOReadBuffSize)
			n, addr, err := uDPConn.ReadFromUDP(buffer)
			if err != nil {
				zlog.Ins().ErrorF("Read UDP err: %v", err)
				continue
			}

			// 处理业务逻辑
			go s.handleUDPData(uDPConn, buffer[:n], addr)
		}
	}()
	select {
	case <-s.exitChan:
		err := uDPConn.Close()
		if err != nil {
			zlog.Ins().ErrorF("listener close err: %v", err)
		}
		return
	}
}

// 为udp创建单独的管理器 处理udp业务
func (s *Server) handleUDPData(uDPConn *net.UDPConn, data []byte, addr *net.UDPAddr) {
	fmt.Println("---------->UDP读到数据", data)
	_, err := uDPConn.WriteToUDP(data, addr)
	if err != nil {
		fmt.Println("---------->WriteToUDP error", err)
		return
	}
}
