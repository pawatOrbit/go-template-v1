package common

type TransportError struct {
	Code        int
	Description string
}

func (c TransportError) Error() string {
	return c.Description
}