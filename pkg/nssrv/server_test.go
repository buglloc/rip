package nssrv_test

import (
	"context"
	"fmt"
	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/nssrv"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() { _ = l.Close() }()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func newSRV(t *testing.T) *nssrv.NSSrv {
	port, err := getFreePort()
	require.NoError(t, err)

	cfg.Zones = []string{"tst"}
	cfg.Addr = fmt.Sprintf("localhost:%d", port)
	srv, err := nssrv.NewSrv()
	require.NoError(t, err)

	go func() {
		err = srv.ListenAndServe()
	}()

	// TODO((buglloc): too ugly
	time.Sleep(1 * time.Second)
	if err != nil {
		_ = srv.Shutdown(context.Background())
		require.NoError(t, err)
	}

	return srv
}

func resolve(t *testing.T, client *dns.Client, msg *dns.Msg) net.IP {
	res, _, err := client.Exchange(msg, cfg.Addr)
	require.NoError(t, err)
	require.NotEmpty(t, res.Answer)

	var ip net.IP
	switch v := res.Answer[0].(type) {
	case *dns.A:
		ip = v.A.To4()
	case *dns.AAAA:
		ip = v.AAAA.To16()
	}

	return ip
}

func TestServer_simple(t *testing.T) {
	cases := []struct {
		in      string
		reqType uint16
		ip      net.IP
	}{
		{
			in:      "1-1-1-1.4.tst",
			reqType: dns.TypeA,
			ip:      net.ParseIP("1.1.1.1").To4(),
		},
		{
			in:      "1-1-1-1.v4.tst",
			reqType: dns.TypeA,
			ip:      net.ParseIP("1.1.1.1").To4(),
		},
		{
			in:      "1-1-1-1.v4.tst",
			reqType: dns.TypeA,
			ip:      net.ParseIP("1.1.1.1").To4(),
		},
		{
			in:      "fe80--fa94-c2ff-fee5-3cf6.6.tst",
			reqType: dns.TypeAAAA,
			ip:      net.ParseIP("fe80::fa94:c2ff:fee5:3cf6").To16(),
		},
		{
			in:      "fe80000000000000fa94c2fffee53cf6.v6.tst",
			reqType: dns.TypeAAAA,
			ip:      net.ParseIP("fe80::fa94:c2ff:fee5:3cf6").To16(),
		},
		{
			in:      "2-2-2-2.3-3-3-3.4.l.tst",
			reqType: dns.TypeA,
			ip:      net.ParseIP("3.3.3.3").To4(),
		},
		{
			in:      "2-2-2-2.3-3-3-3.4.s.tst",
			reqType: dns.TypeA,
			ip:      net.ParseIP("3.3.3.3").To4(),
		},
	}

	srv := newSRV(t)
	defer func() { _ = srv.Shutdown(context.Background()) }()

	client := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			msg := &dns.Msg{}
			msg.SetQuestion(dns.Fqdn(tc.in), tc.reqType)
			ip := resolve(t, client, msg)
			require.Equal(t, tc.ip, ip)
		})
	}
}

func TestServer_loop(t *testing.T) {
	srv := newSRV(t)
	defer func() { _ = srv.Shutdown(context.Background()) }()

	client := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn("1-1-1-1.v4.2-2-2-2.v4.loop.tst"), dns.TypeA)
	ip := resolve(t, client, msg)
	require.Equal(t, net.ParseIP("2.2.2.2").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("1.1.1.1").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("2.2.2.2").To4(), ip)
}

func TestServer_multiLoop(t *testing.T) {
	srv := newSRV(t)
	defer func() { _ = srv.Shutdown(context.Background()) }()

	client := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn("1-1-1-1.v4.2-2-2-2.v4.loop-cnt-2.3-3-3-3.v4.loop.tst"), dns.TypeA)
	ip := resolve(t, client, msg)
	require.Equal(t, net.ParseIP("3.3.3.3").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("2.2.2.2").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("1.1.1.1").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("3.3.3.3").To4(), ip)
}

func TestServer_multiLoopWithTTL(t *testing.T) {
	srv := newSRV(t)
	defer func() { _ = srv.Shutdown(context.Background()) }()

	client := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn("1-1-1-1.v4.2-2-2-2.v4.loop-ttl-20s.3-3-3-3.v4.loop.tst"), dns.TypeA)
	ip := resolve(t, client, msg)
	require.Equal(t, net.ParseIP("3.3.3.3").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("2.2.2.2").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("1.1.1.1").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("2.2.2.2").To4(), ip)
	ip = resolve(t, client, msg)
	require.Equal(t, net.ParseIP("1.1.1.1").To4(), ip)
}