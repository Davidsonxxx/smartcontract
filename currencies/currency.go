package currencies

type Currency int8

const (
	// don't use iota to make it more explicit
	Bitcoin Currency = 0
	Ether Currency = 1
)
