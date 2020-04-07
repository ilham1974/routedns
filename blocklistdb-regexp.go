package rdns

import (
	"net"
	"regexp"

	"github.com/miekg/dns"
)

// RegexpDB holds a list of regular expressions against which it evaluates DNS queries.
type RegexpDB struct {
	filters []*regexp.Regexp
}

var _ BlocklistDB = &RegexpDB{}

// NewRegexpDB returns a new instance of a matcher for a list of regular expressions.
func NewRegexpDB(rules ...string) (*RegexpDB, error) {
	var filters []*regexp.Regexp
	for _, r := range rules {
		re, err := regexp.Compile(r)
		if err != nil {
			return nil, err
		}
		filters = append(filters, re)
	}

	return &RegexpDB{filters}, nil
}

func (m *RegexpDB) New(rules []string) (BlocklistDB, error) {
	return NewRegexpDB(rules...)
}

func (m *RegexpDB) Match(q dns.Question) (net.IP, bool) {
	for _, filter := range m.filters {
		if filter.MatchString(q.Name) {
			return nil, true
		}
	}
	return nil, false
}

func (m *RegexpDB) String() string {
	return "Regexp"
}