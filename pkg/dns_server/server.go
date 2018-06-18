package dns_server

import (
	"strings"
	"time"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
)

func RunBackground() error {
	tcpHandler := dns.NewServeMux()
	udpHandler := dns.NewServeMux()

	for _, zone := range cfg.Zones {
		if !strings.HasSuffix(zone, ".") {
			zone += "."
		}
		tcpHandler.HandleFunc(zone, NewHandler(zone))
		udpHandler.HandleFunc(zone, NewHandler(zone))
	}

	tcpServer := &dns.Server{
		Addr:         cfg.Addr,
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	udpServer := &dns.Server{
		Addr:         cfg.Addr,
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("starting TCP-server", "addr", tcpServer.Addr)
		if err := tcpServer.ListenAndServe(); err != nil {
			log.Error("TCP-server start failed", "err", err.Error())
		}
	}()
	go func() {
		log.Info("starting UDP-server", "addr", udpServer.Addr)
		if err := udpServer.ListenAndServe(); err != nil {
			log.Error("UDP-server start failed", "err", err.Error())
		}
	}()

	return nil
}
