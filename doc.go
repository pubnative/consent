/*
Package consent contains IAB consent string encoding and decoding
implementation.

To create a new consent use the consent.New function:

 cmpID := 1
 cmpVersion := 1
 consentScreen := byte(1)
 lang := "EN"
 vendorListVersion := 42
 purposesAllowed := [24]bool{true, true, true, true}
 allowedVendors := map[int]bool{10: true, 64: true}

 c1 := consent.New(cmpID, cmpVersion, consentScreen, lang,
 	vendorListVersion, purposesAllowed, allowedVendors)

To decode an existing consent use consent.Parse:

 c2, err := consent.Parse("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")

In both cases you can modify the consent struct as you want:

 c2.Vendors[999] = true
 c2.LastUpdated = time.Now()

Then you can encode it as IAB base64 encoded string:

 c2.String()

*/
package consent
