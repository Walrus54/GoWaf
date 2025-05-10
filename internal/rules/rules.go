package rules

import (
	"regexp"
)

type Rule struct {
	ID         string
	Name       string
	Pattern    *regexp.Regexp `yaml:"-"`
	RawPattern string         `yaml:"pattern"`
}

func (r *Rule) Compile() error {
	pattern, err := regexp.Compile(r.RawPattern)
	if err != nil {
		return err
	}
	r.Pattern = pattern
	return nil
}
