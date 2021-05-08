package nssrv

import (
	"context"
	"strings"
	"time"

	log "github.com/buglloc/simplelog"
	"github.com/karlseguin/ccache"
	"github.com/miekg/dns"
	"golang.org/x/sync/errgroup"

	"github.com/buglloc/rip/v2/pkg/cfg"
)

type NSSrv struct {
	tcpServer *dns.Server
	udpServer *dns.Server
	cache     *ccache.Cache
}

func NewSrv() (*NSSrv, error) {
	srv := &NSSrv{
		tcpServer: &dns.Server{
			Addr:         cfg.Addr,
			Net:          "tcp",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		udpServer: &dns.Server{
			Addr:         cfg.Addr,
			Net:          "udp",
			UDPSize:      65535,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		cache: ccache.New(ccache.Configure().MaxSize(cfg.CacheSize)),
	}

	srv.udpServer.Handler = srv.newDNSRouter()
	srv.tcpServer.Handler = srv.newDNSRouter()
	return srv, nil
}

func (s *NSSrv) ListenAndServe() error {
	var g errgroup.Group
	g.Go(func() error {
		log.Info("starting TCP-server", "addr", s.tcpServer.Addr)
		err := s.tcpServer.ListenAndServe()
		if err != nil {
			log.Error("can't start TCP-server", "err", err)
		}
		return err
	})

	g.Go(func() error {
		log.Info("starting UDP-server", "addr", s.tcpServer.Addr)
		err := s.udpServer.ListenAndServe()
		if err != nil {
			log.Error("can't start UDP-server", "err", err)
		}
		return err
	})

	return g.Wait()
}

func (s *NSSrv) Shutdown(ctx context.Context) error {
	var g errgroup.Group
	g.Go(func() error {
		return s.tcpServer.ShutdownContext(ctx)
	})

	g.Go(func() error {
		return s.udpServer.ShutdownContext(ctx)
	})

	done := make(chan error)
	go func() {
		defer close(done)
		done <- g.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (s *NSSrv) newDNSRouter() *dns.ServeMux {
	out := dns.NewServeMux()
	for _, zone := range cfg.Zones {
		if !strings.HasSuffix(zone, ".") {
			zone += "."
		}

		out.HandleFunc(zone, func(zone string) func(w dns.ResponseWriter, req *dns.Msg) {
			return func(w dns.ResponseWriter, req *dns.Msg) {
				defer func() { _ = w.Close() }()

				l := log.Child("client", w.RemoteAddr().String())
				msg := s.handleRequest(zone, req, &l)
				_ = w.WriteMsg(msg)
			}
		}(zone))
	}
	return out
}
