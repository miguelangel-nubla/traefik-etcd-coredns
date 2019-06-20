package spec

type Record struct {
	DNSName  string `json:"-"`
	Host     string `json:"host,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Mail     bool   `json:"mail,omitempty"` // Be an MX record. Priority becomes Preference.
	TTL      uint32 `json:"ttl,omitempty"`
}

const (
	RecordTypeA     = "A"
	RecordTypeCNAME = "CNAME"
)
