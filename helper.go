package dyn

// NewTXTRecord creates a new txt record struct
func NewTXTRecord(txtdata string, ttl int) *Record {
	return &Record{
		TTL: ttl,
		RData: DataBlock{
			TxtData: txtdata,
		},
	}
}
