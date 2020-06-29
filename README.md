# IAB Consent
[![GoDoc](https://godoc.org/github.com/pubnative/consent?status.png)](https://godoc.org/github.com/pubnative/consent)
[![CircleCI](https://circleci.com/gh/pubnative/consent.svg?style=svg)](https://circleci.com/gh/pubnative/consent)

A minimalistic Go library to encode and decode [IAB consent][iab] strings.

[iab]: https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/68f5e0012a7bdb00867ce9fee57fb67cfe9153e3/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md

### Usage examples

#### Version 1

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

	// Decode an existing consent
	c2, err := consent.Parse("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Last modified: %s, vendors allowed: %v\n",
		c1.LastUpdated, c1.Vendors)

	// Patch an existing consent
	c2.Vendors[999] = true
	fmt.Println(c2.String())
}
```

### Version 2

```go
	// Decode a consent v2 string
	cv2, err := consent.ParseV2("COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Last modified: %s, vendors allowed: %v\n",
		cv2.LastUpdated, cv2.VendorConsent)
```

### Consent string of unknown version

```go
	// decode a consent string without knowing the version beforehand
	cvx, err := consent.Parse("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")
	if err != nil {
	fmt.Println(err)
		os.Exit(1)
	}
	switch c3 := cvx.(type) {
	case *consent.ConsentV1:
		fmt.Printf("V1. PurposesAllowed: %v\n", c3.PurposesAllowed)
	case *consent.ConsentV2:
		fmt.Printf("V2. PurposesConsent: %v\n", c3.PurposesConsent)
	}
```
