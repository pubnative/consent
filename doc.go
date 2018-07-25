/*
Package consent contains IAB consent string encode and decode implementations.

Create a new consent with the consent.New function:

 cmpID := 1
 cmpVersion := 1
 consentScreen := byte(1)
 lang := "EN"
 vendorListVersion := 42
 purposesAllowed := [24]bool{true, true, true, true}
 allowedVendors := map[int]bool{10: true, 64: true}

 c1 := consent.New(cmpID, cmpVersion, consentScreen, lang,
 	vendorListVersion, purposesAllowed, allowedVendors)

Decode an existing string with consent.Parse:

 c2, err := consent.Parse("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA")

In both cases you can read and modify a struct as you want:

 if c2.Vendors[55] {
	 // Vendor #55 is permitted
 }
 c2.Vendors[999] = true
 c2.LastUpdated = time.Now()

At anytime export it as an IAB base64-encoded consent string:

 c2.String()

*/
package consent
