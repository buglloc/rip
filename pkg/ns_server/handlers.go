package ns_server

import (
	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

func NewHandler(zone string) func(w dns.ResponseWriter, req *dns.Msg) {
	return func(w dns.ResponseWriter, req *dns.Msg) {
		defer w.Close()
		msg := handle(zone, req)
		if msg != nil {
			w.WriteMsg(msg)
		}
	}
}

func handle(zone string, req *dns.Msg) *dns.Msg {
	response := &dns.Msg{}
	response.SetReply(req)
	for _, question := range req.Question {
		switch question.Qtype {
		case dns.TypeA, dns.TypeAAAA:
			msg, ip, err := parseName(question, zone)
			if err != nil {
				log.Error("failed to parse request", "type", typeToString(question.Qtype), "name", question.Name, "err", err.Error())
				continue
			}

			if msg != nil {
				// parser craft own response
				response.Answer = append(response.Answer, msg)
				continue
			}

			if ip == nil {
				continue
			}

			head := dns.RR_Header{
				Name:   question.Name,
				Rrtype: question.Qtype,
				Class:  dns.ClassINET,
				Ttl:    0,
			}

			var line dns.RR
			if question.Qtype == dns.TypeA {
				line = &dns.A{
					Hdr: head,
					A:   ip,
				}
			} else {
				line = &dns.AAAA{
					Hdr:  head,
					AAAA: ip,
				}
			}
			response.Answer = append(response.Answer, line)
		default:
			log.Debug("skip unknown request", "type", typeToString(question.Qtype))
			// TODO(buglloc): should we return SERVFAIL?
			//msg := &dns.Msg{}
			//msg.SetRcode(req, dns.RcodeServerFailure)
		}
	}

	return response
}
