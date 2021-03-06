package rdns

import (
	"fmt"

	"github.com/miekg/dns"
)

// StaticResolver is a resolver that always returns the same answer, to any question.
// Typically used in combination with a blocklist to define fixed block responses or
// with a router when building a walled garden.
type StaticResolver struct {
	answer []dns.RR
	ns     []dns.RR
	extra  []dns.RR
	rcode  int
}

var _ Resolver = &StaticResolver{}

type StaticResolverOptions struct {
	// Records in zone-file format
	Answer []string
	NS     []string
	Extra  []string
	RCode  int
}

// NewStaticResolver returns a new instance of a StaticResolver resolver.
func NewStaticResolver(opt StaticResolverOptions) (*StaticResolver, error) {
	r := new(StaticResolver)

	for _, record := range opt.Answer {
		rr, err := dns.NewRR(record)
		if err != nil {
			return nil, err
		}
		r.answer = append(r.answer, rr)
	}
	for _, record := range opt.NS {
		rr, err := dns.NewRR(record)
		if err != nil {
			return nil, err
		}
		r.ns = append(r.ns, rr)
	}
	for _, record := range opt.Extra {
		rr, err := dns.NewRR(record)
		if err != nil {
			return nil, err
		}
		r.extra = append(r.extra, rr)
	}
	r.rcode = opt.RCode

	return r, nil
}

// Resolve a DNS query by returning a fixed response.
func (r *StaticResolver) Resolve(q *dns.Msg, ci ClientInfo) (*dns.Msg, error) {
	answer := new(dns.Msg)
	answer.SetReply(q)

	// Update the name of every answer record to match that of the query
	answer.Answer = make([]dns.RR, 0, len(r.answer))
	for _, rr := range r.answer {
		r := dns.Copy(rr)
		r.Header().Name = qName(q)
		answer.Answer = append(answer.Answer, r)
	}
	answer.Ns = r.ns
	answer.Extra = r.extra
	answer.Rcode = r.rcode

	return answer, nil
}

func (r *StaticResolver) String() string {
	return fmt.Sprintf("StaticResolver")
}
