package ns_server

import (
	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

func NewHandler(zone string) func(w dns.ResponseWriter, req *dns.Msg) {
	return func(w dns.ResponseWriter, req *dns.Msg) {
		defer w.Close()
		msg := handle(zone, req, log.Child("client", w.RemoteAddr().String()))
		if msg != nil {
			w.WriteMsg(msg)
		}
	}
}

func handle(zone string, req *dns.Msg, l log.Logger) *dns.Msg {
	response := &dns.Msg{}
	response.SetReply(req)
	for _, question := range req.Question {
		switch question.Qtype {
		case dns.TypeA, dns.TypeAAAA:
			answers, err := parseName(question, zone, l)
			if err != nil {
				l.Error("failed to parse request", "type", typeToString(question.Qtype), "name", question.Name, "err", err.Error())
				continue
			}

			if len(answers) == 0 {
				continue
			}
			response.Answer = append(response.Answer, answers...)
		default:
			l.Debug("skip unknown request", "type", typeToString(question.Qtype))
			// TODO(buglloc): should we return SERVFAIL?
			//msg := &dns.Msg{}
			//msg.SetRcode(req, dns.RcodeServerFailure)
		}
	}

	return response
}
