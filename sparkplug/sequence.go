package sparkplug

// NextSequenceNumber returns the next sequence number based on current value.
// All messages published have a sequence number starting at 0 and ending at 255 after which it is reset to 0.
func NextSequenceNumber(current uint64) uint64 {
	if current > 255 {
		panic("sequence number must be between 0 and 255")
	}
	if current == 255 {
		return 0
	}
	return current + 1
}
