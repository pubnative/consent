package consent

import (
	"testing"
)

const (
	csv1 = "BOEFEAyOEFEAyAHABDENAI4AAAB9vABAASA"
	csv2 = "COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA"
)

func TestParseVersion(t *testing.T) {
	v, err := ParseConsentVersion(csv1)
	noError(t, err)
	equal(t, byte(1), v)

	v, err = ParseConsentVersion(csv2)
	noError(t, err)
	equal(t, byte(2), v)

	v, err = ParseConsentVersion("X32g")
	equal(t, err, ErrUnsupported)
}

func TestParseConsentString(t *testing.T) {
	v, err := Parse(csv1)
	equal(t, nil, err)
	equal(t, byte(1), v.Version())
	v1, ok := v.(*ConsentV1)
	equal(t, ok, true)
	if v1 == nil {
		t.Fail()
	}
	equal(t, "EN", v1.ConsentLanguage)

	v, err = Parse(csv2)
	equal(t, nil, err)
	equal(t, byte(2), v.Version())
	v2, ok := v.(*ConsentV2)
	equal(t, ok, true)
	if v2 == nil {
		t.Fail()
	}
	equal(t, 2, v2.CmpVersion)
}
