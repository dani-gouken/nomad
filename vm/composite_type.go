package vm

import "strings"

type CompositeType struct {
	Alias string
	Cases []string
}

func NewCompositeType(alias string, types ...string) CompositeType {
	return CompositeType{
		Cases: types,
		Alias: alias,
	}
}

func (c *CompositeType) GetSignature() string {
	return strings.Join(c.Cases, "|")
}
