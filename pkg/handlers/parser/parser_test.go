package parser_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/cname"
	"github.com/buglloc/rip/v2/pkg/handlers/defaultip"
	"github.com/buglloc/rip/v2/pkg/handlers/ipv4"
	"github.com/buglloc/rip/v2/pkg/handlers/ipv6"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
	"github.com/buglloc/rip/v2/pkg/handlers/loop"
	"github.com/buglloc/rip/v2/pkg/handlers/parser"
	"github.com/buglloc/rip/v2/pkg/handlers/proxy"
)

func init() {
	// TODO(buglloc): so ugly
	cfg.AllowProxy = true
}

func TestParser(t *testing.T) {
	cases := []struct {
		in       string
		handlers []handlers.Handler
	}{
		{
			in: "",
		},
		{
			in: "lalala",
		},
		{
			in: "1-1-1-1",
			handlers: []handlers.Handler{
				&ipv4.Handler{IP: net.ParseIP("1.1.1.1").To4()},
			},
		},
		{
			in: "1010101",
			handlers: []handlers.Handler{
				&ipv4.Handler{IP: net.ParseIP("1.1.1.1").To4()},
			},
		},
		{
			in: "1-1-1-1.4",
			handlers: []handlers.Handler{
				&ipv4.Handler{IP: net.ParseIP("1.1.1.1").To4()},
			},
		},
		{
			in: "1010101.4",
			handlers: []handlers.Handler{
				&ipv4.Handler{IP: net.ParseIP("1.1.1.1").To4()},
			},
		},
		{
			in: "fe80--fa94-c2ff-fee5-3cf6",
			handlers: []handlers.Handler{
				&ipv6.Handler{IP: net.ParseIP("fe80::fa94:c2ff:fee5:3cf6").To16()},
			},
		},
		{
			in: "fe80000000000000fa94c2fffee53cf6",
			handlers: []handlers.Handler{
				&ipv6.Handler{IP: net.ParseIP("fe80::fa94:c2ff:fee5:3cf6").To16()},
			},
		},
		{
			in: "fe80--fa94-c2ff-fee5-3cf6.6",
			handlers: []handlers.Handler{
				&ipv6.Handler{IP: net.ParseIP("fe80::fa94:c2ff:fee5:3cf6").To16()},
			},
		},
		{
			in: "fe80000000000000fa94c2fffee53cf6.6",
			handlers: []handlers.Handler{
				&ipv6.Handler{IP: net.ParseIP("fe80::fa94:c2ff:fee5:3cf6").To16()},
			},
		},
		{
			in: "lalala.d.d.d",
			handlers: []handlers.Handler{
				&defaultip.Handler{},
				&defaultip.Handler{},
				&defaultip.Handler{},
			},
		},
		{
			in: "lalala.example-com.c",
			handlers: []handlers.Handler{
				&cname.Handler{TargetFQDN: "example.com."},
			},
		},
		{
			in: "lalala.example-com.p",
			handlers: []handlers.Handler{
				&proxy.Handler{TargetFQDN: "example.com."},
			},
		},
		{
			in: "lalala.d.lala-com.p.d.example-com.c.d",
			handlers: []handlers.Handler{
				&defaultip.Handler{},
				&cname.Handler{TargetFQDN: "example.com."},
				&defaultip.Handler{},
				&proxy.Handler{TargetFQDN: "lala.com."},
				&defaultip.Handler{},
			},
		},
		{
			in: "1-1-1-1.4.2-2-2-2.4.3-3-3-3.4.l",
			handlers: []handlers.Handler{
				&loop.Handler{
					Len: 3,
					Nested: []handlers.Handler{
						&ipv4.Handler{IP: net.ParseIP("3.3.3.3").To4()},
						&ipv4.Handler{IP: net.ParseIP("2.2.2.2").To4()},
						&ipv4.Handler{IP: net.ParseIP("1.1.1.1").To4()},
					},
				},
			},
		},
		{
			in: "1-1-1-1.4-ttl-10s.2-2-2-2.4-cnt-1.3-3-3-3.4-cnt-2.l",
			handlers: []handlers.Handler{
				&loop.Handler{
					Len: 3,
					Nested: []handlers.Handler{
						&ipv4.Handler{
							IP: net.ParseIP("3.3.3.3").To4(),
							BaseHandler: handlers.BaseHandler{
								Limiters: limiter.Limiters{
									&limiter.Count{
										Max: 2,
									},
								},
							},
						},
						&ipv4.Handler{
							IP: net.ParseIP("2.2.2.2").To4(),
							BaseHandler: handlers.BaseHandler{
								Limiters: limiter.Limiters{
									&limiter.Count{
										Max: 1,
									},
								},
							},
						},
						&ipv4.Handler{
							IP: net.ParseIP("1.1.1.1").To4(),
							BaseHandler: handlers.BaseHandler{
								Limiters: limiter.Limiters{
									&limiter.TTL{
										TTL: 10 * time.Second,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			in: "1-1-1-1.4-ttl-10s.2-2-2-2.3-3-3-3.s",
			handlers: []handlers.Handler{
				&loop.Handler{
					Len: 3,
					Nested: []handlers.Handler{
						&ipv4.Handler{
							IP: net.ParseIP("3.3.3.3").To4(),
						},
						&ipv4.Handler{
							IP: net.ParseIP("2.2.2.2").To4(),
							BaseHandler: handlers.BaseHandler{
								Limiters: limiter.Limiters{
									&limiter.TTL{
										TTL: cfg.StickyTTL,
									},
								},
							},
						},
						&ipv4.Handler{
							IP: net.ParseIP("1.1.1.1").To4(),
							BaseHandler: handlers.BaseHandler{
								Limiters: limiter.Limiters{
									&limiter.TTL{
										TTL: 10 * time.Second,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			hndlrs, err := parser.NewParser(tc.in).All()
			require.NoError(t, err)
			require.EqualValues(t, tc.handlers, hndlrs)
		})
	}
}
