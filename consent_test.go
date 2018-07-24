package consent

import (
	"reflect"
	"testing"
	"time"
)

func noError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
}

func equal(t *testing.T, expected, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected '%s', got '%s'", expected, got)
	}
}

func TestConsentString_OfficialDocTest(t *testing.T) {
	t2017, _ := time.Parse(time.RFC3339Nano, "2017-11-07T19:15:55.4Z")
	c := Consent{
		Version:           1,
		Created:           t2017,
		LastUpdated:       t2017,
		CmpID:             7,
		CmpVersion:        1,
		ConsentScreen:     3,
		ConsentLanguage:   "EN",
		VendorListVersion: 8,
		PurposesAllowed:   [24]bool{true, true, true, false},

		MaxVendorID:  2011,
		encodingType: rangeType,

		Vendors: make(map[int]bool),
	}
	for i := 0; i < 2012; i++ {
		c.Vendors[i] = true
	}
	delete(c.Vendors, 9)

	equal(t, "BOEFEAyOEFEAyAHABDENAI4AAAB9vABAASA", c.String())
}

func TestConsentParse_OfficialDocTest(t *testing.T) {
	var c Consent
	t2017, _ := time.Parse(time.RFC3339Nano, "2017-11-07T19:15:55.4Z")
	noError(t, c.Parse("BOEFEAyOEFEAyAHABDENAI4AAAB9vABAASA"))

	if c.Version != 1 {
		t.Fail()
	}
	if c.Created.UTC() != t2017 {
		t.Fail()
	}
	if c.LastUpdated.UTC() != t2017 {
		t.Fail()
	}

	equal(t, "EN", c.ConsentLanguage)
	if [24]bool{true, true, true, false} != c.PurposesAllowed {
		t.Fail()
	}

	vendors := make(map[int]bool)
	for i := 1; i < 2012; i++ {
		vendors[i] = true
	}
	vendors[9] = false
	if !reflect.DeepEqual(vendors, c.Vendors) {
		t.Fail()
	}
}

func TestConsentParse_JSSDK(t *testing.T) {
	var c Consent
	t2018, err := time.Parse("2006-01-02 15:04:05 MST", "2018-07-15 07:00:00 PDT")
	noError(t, err)
	noError(t, c.Parse("BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA"))

	equal(t, 1, int(c.Version))
	equal(t, t2018.UTC(), c.Created.UTC())
	equal(t, t2018.UTC(), c.LastUpdated.UTC())

	equal(t, "EN", c.ConsentLanguage)
	equal(t, [24]bool{true, true, false, false}, c.PurposesAllowed)
	equal(t, map[int]bool{1: true, 2: true, 4: true}, c.Vendors)
}

func TestConsentParse_BackNForth(t *testing.T) {
	fixtures := map[string]string{
		// example from official docs
		"BOEFEAyOEFEAyAHABDENAI4AAAB9vABAASA": "BOEFEAyOEFEAyAHABDENAI4AAAB9vABAASA", // Official

		// From Java SDK
		// BitField
		"BONMj34ONMj34ABACDENALqAAAAAplY": "BONMj34ONMj34ABACDENALqAAAAAplY",
		// Ranges
		"BN5lERiOMYEdiAOAWeFRAAYAAaAAptQ":                  "BN5lERiOMYEdiAOAWeFRAAYAAaAAptQ",
		"BN5lERiOMYEdiAKAWXEND1HoSBE6CAEAApAMgBmgEBxOIE6A": "BN5lERiOMYEdiAKAWXEND1HoSBE6CAEAApAMgBmgEBxOIE6A",
		"BOOMzbgOOQww_AtABAFRAb-AAAsvOA3gACAAkABgArgBaAF0AMAA1gBuAH8AQQBSgCoAL8AYQBigDIAM0AaABpgDYAOYAdgA8AB6gD4AQoAiABFQCMAI6ASABIgCTAEqAJeATIBQQCiAKSAU4BVQCtAK-AWYBaQC2ALcAXMAvAC-gGAAYcAxQDGAGQAMsAZsA0ADTAGqANcAbMA4ADjAHKAOiAdQB1gDtgHgAeMA9AD2AHzAP4BAACBAEEAIbAREBEgCKQEXARhZeYA": "BOOMzbgOOQww_AtABAFRAb-AAAsvPA2AAKACwAF4ANgAgABTADAAGMAM8AagBrgDoAOoAdwA8gB7gEMAQ4AiQBFgCPAEkAJQASwAmABQwClAKaAVYBWQCwALIAWoAuIBdAF2AL8AYgAx4BkgGUAMyAZwBngDUAGsANiAbQBvgDkgHMAc4A6QB2QDuAO-AeQB5wD3APiAfQB-gEBAIHAQUBDICHAIgAROAioCLQEZsvI",
		"BONZt-1ONZt-1AHABBENAO-AAAAHCAEAASABmADYAOAAeA": "BONZt-1ONZt-1AHABBENAO-AAAAHCAEAASABmADYAOAAeA",

		// From JS SDK
		// Ranges
		"BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA": "BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA",
	}

	for orig, short := range fixtures {
		t.Logf(orig)
		var c Consent
		err := c.Parse(orig)
		noError(t, err)
		equal(t, short, c.String())

		var c2 Consent
		err = c2.Parse(c.String())
		noError(t, err)
		equal(t, c.String(), c2.String())
	}
}
