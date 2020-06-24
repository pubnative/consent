/*
Package consent contains IAB consent string encode and decode implementations.

Version 1

Create a new consent (v1) with the consent.NewConsentV1 function:

 cmpID := 1
 cmpVersion := 1
 consentScreen := byte(1)
 lang := "EN"
 vendorListVersion := 42
 purposesAllowed := [24]bool{true, true, true, true}
 allowedVendors := map[int]bool{10: true, 64: true}

 c1 := consent.NewConsentV1(cmpID, cmpVersion, consentScreen, lang,
 	vendorListVersion, purposesAllowed, allowedVendors)

Decode an existing string with consent.ParseV1:

 c2, err := consent.ParseV1("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")

In both cases you can read and modify a struct as you want:

 if c2.Vendors[55] {
	 // Vendor #55 is permitted
 }
 c2.Vendors[999] = true
 c2.LastUpdated = time.Now()

At anytime export it as an IAB base64-encoded consent string:

 c2.String()

Parsing a consent string if you don't know its version

 c, err = consent.Parse("COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA")
 println(c.Version()) // prints 2
 // you can cast the parsed value to the matching Consent struct:
 cv2, ok := v.(*ConsentV2) 

You can also only parse the version information from a consent string
without parsing the whole string:

 version, err = ParseConsentVersion("BOEFEAyOEFEAyAHABDENAI4AAAB9vABAASA")
 println(version) // prints 1

*/
package consent
