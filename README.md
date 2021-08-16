# Go package with improved reflect.DeepEqual

[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/deepequal.svg)](https://pkg.go.dev/github.com/powerman/deepequal)
[![CI/CD](https://github.com/powerman/deepequal/workflows/CI/CD/badge.svg?event=push)](https://github.com/powerman/deepequal/actions?query=workflow%3ACI%2FCD)
[![Coverage Status](https://coveralls.io/repos/github/powerman/deepequal/badge.svg?branch=master)](https://coveralls.io/github/powerman/deepequal?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/powerman/deepequal)](https://goreportcard.com/report/github.com/powerman/deepequal)
[![Release](https://img.shields.io/github/v/release/powerman/deepequal)](https://github.com/powerman/deepequal/releases/latest)

Most of the code is copied from Go reflect package with slight
modifications.

Differences from reflect.DeepEqual:

- If compared value implements `.Equal(valueOfSameType) bool` method then
  it will be called instead of comparing values as is.
- If called `Equal` method will panics then whole DeepEqual will panics too.

This means you can use this DeepEqual method to correctly compare types
like time.Time or decimal.Decimal, without taking in account unimportant
differences (like time zone or exponent).
