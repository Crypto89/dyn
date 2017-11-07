package dyn

// Zone struct
type Zone struct {
	c    *Client
	name string
}

// Zone helper
func (c *Client) Zone(name string) *Zone {
	return &Zone{c: c, name: name}
}

// ARecord helper
func (z *Zone) ARecord(fqdn string) *RecordFunc {
	return &RecordFunc{c: z.c, object: "ARecord", zone: z, fqdn: fqdn}
}

// AAAARecord helper
func (z *Zone) AAAARecord(fqdn string) *RecordFunc {
	return &RecordFunc{c: z.c, object: "AAAARecord", zone: z, fqdn: fqdn}
}

// TXTRecord helper
func (z *Zone) TXTRecord(fqdn string) *RecordFunc {
	return &RecordFunc{c: z.c, object: "TXTRecord", zone: z, fqdn: fqdn}
}
