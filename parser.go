package main

import (
	gerr "github.com/corebreaker/goerrors"
	"github.com/gorhill/cronexpr"
)

func ParseSpec(spec string) (*cronexpr.Expression, error) {
	expr, err := cronexpr.Parse(spec)
	if err != nil {
		return nil, gerr.DecorateError(err)
	}

	return expr, nil
}