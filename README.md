# Go package with improved reflect.DeepEqual

[![License MIT](https://img.shields.io/badge/license-MIT-royalblue.svg)](LICENSE)
[![License BSD3](https://img.shields.io/badge/license-BSD3-royalblue.svg)](LICENSE-go)
[![Go version](https://img.shields.io/github/go-mod/go-version/powerman/deepequal?color=blue)](https://go.dev/)
[![Test](https://img.shields.io/github/actions/workflow/status/powerman/deepequal/test.yml?label=test)](https://github.com/powerman/deepequal/actions/workflows/test.yml)
[![Coverage Status](https://raw.githubusercontent.com/powerman/deepequal/gh-badges/coverage.svg)](https://github.com/powerman/deepequal/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/powerman/deepequal?color=blue)](https://github.com/powerman/deepequal/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/deepequal.svg)](https://pkg.go.dev/github.com/powerman/deepequal)

![Linux | amd64 arm64 armv7 ppc64le s390x riscv64](https://img.shields.io/badge/Linux-amd64%20arm64%20armv7%20ppc64le%20s390x%20riscv64-royalblue)
![macOS | amd64 arm64](https://img.shields.io/badge/macOS-amd64%20arm64-royalblue)
![Windows | amd64 arm64](https://img.shields.io/badge/Windows-amd64%20arm64-royalblue)

Most of the code is copied from Go reflect package with slight
modifications.

Differences from reflect.DeepEqual:

- If compared value implements `.Equal(valueOfSameType) bool` method then
  it will be called instead of comparing values as is.
- If called `Equal` method will panics then whole DeepEqual will panics too.

This means you can use this DeepEqual method to correctly compare types
like time.Time or decimal.Decimal, without taking in account unimportant
differences (like time zone or exponent).
