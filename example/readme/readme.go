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

	c1 := consent.NewConsentV1(cmpID, cmpVersion, consentScreen, lang,
		vendorListVersion, purposesAllowed, allowedVendors)

	fmt.Println(c1.String())

	// Decode an existing consent
	c2, err := consent.ParseV1("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Last modified: %s, vendors allowed: %v\n",
		c1.LastUpdated, c1.Vendors)

	// Patch an existing consent
	c2.Vendors[999] = true
	fmt.Println(c2.String())

	// Decode a consent v2 string
	cv2, err := consent.ParseV2("COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Last modified: %s, vendors allowed: %v\n",
		cv2.LastUpdated, cv2.VendorConsent)

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
}
