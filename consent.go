package consent

func ParseConsentVersion(s string) (byte, error) {
	if len(s) == 0 {
		return 0, ErrUnexpectedEnd
	}
	// string is base64-encoded and version is encoded in the first 6 bits (i.e. first char in b64 string)
	switch s[0] {
	case 'B':
		return 1, nil
	case 'C':
		return 2, nil
	}
	return 0, ErrUnsupported
}

type Consent interface {
	Version() byte
	String() string
}

func Parse(s string) (Consent, error) {
	v, err := ParseConsentVersion(s)
	if err != nil {
		return nil, err
	}
	switch v {
	case 1:
		return ParseV1(s)
	case 2:
		return ParseV2(s)
	}
	return nil, nil // unreachable
}

func Validate(s string) error {
	v, err := ParseConsentVersion(s)
	if err != nil {
		return err
	}
	switch v {
	case 1:
		return ValidateV1(s)
	case 2:
		return ValidateV2(s)
	}
	return nil // unreachable
}
