# IAB Consent
[![GoDoc](https://godoc.org/github.com/pubnative/consent?status.png)](https://godoc.org/github.com/pubnative/consent)
[![CircleCI](https://circleci.com/gh/pubnative/consent.svg?style=svg)](https://circleci.com/gh/pubnative/consent)

A Go library for decoding and encoding [IAB consent][iab] strings.

[iab]: https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/68f5e0012a7bdb00867ce9fee57fb67cfe9153e3/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md

### Usage example

```go
package main

import (
	"fmt"
	"os"

	"github.com/pubnative/consent"
)

func main() {
	// Create a new consent string
	cmpID := 1
	cmpVersion := 1
	consentScreen := byte(1)
	lang := "EN"
	vendorListVersion := 42
	purposesAllowed := [24]bool{true, true, true, true}
	allowedVendors := map[int]bool{10: true, 64: true}

	c1 := consent.New(cmpID, cmpVersion, consentScreen, lang,
		vendorListVersion, purposesAllowed, allowedVendors)

	fmt.Println(c1.String())

	// Decode existing consent
	c2, err := consent.Parse("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Last modified: %s, vendors allowed: %v\n",
		c1.LastUpdated, c1.Vendors)

	// Patch the existing one
	c2.Vendors[999] = true
	fmt.Println(c2.String())
}
```
