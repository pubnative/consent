package consent

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

type encodingType byte

const (
	bitFieldType encodingType = 0
	rangeType    encodingType = 1
)

var (
	// ErrUnexpectedEnd is returned when a consent string is too short
	ErrUnexpectedEnd = errors.New("consent: unexpected end")

	// ErrUnsupported is returned when a string version is not 1
	ErrUnsupported = errors.New("consent: version is not supported")
)

// Consent is a golang representation of an IAB consent string
//
// Implementation is done as per IAB Consent String format v1.1:
// https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/68f5e0012a7bdb00867ce9fee57fb67cfe9153e3/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md
type Consent struct {
	Version           byte
	Created           time.Time
	LastUpdated       time.Time
	CmpID             int
	CmpVersion        int
	ConsentScreen     byte
	ConsentLanguage   string
	VendorListVersion int
	PurposesAllowed   [24]bool

	MaxVendorID  int
	encodingType encodingType

	Vendors map[int]bool
}

// New creates a new instance of a Consent struct
func New(cmpID, cmpVersion int, consentScreen byte, lang string, vendorListVerson int, purposesAllowed [24]bool, allowedVendors map[int]bool) Consent {

	vendors := make(map[int]bool)
	for k, v := range allowedVendors {
		vendors[k] = v
	}

	return Consent{
		Version:           1,
		Created:           time.Now(),
		LastUpdated:       time.Now(),
		CmpID:             cmpID,
		CmpVersion:        cmpVersion,
		ConsentScreen:     consentScreen,
		ConsentLanguage:   lang,
		VendorListVersion: vendorListVerson,
		PurposesAllowed:   purposesAllowed,
		Vendors:           vendors,
	}
}

// Bytes returns raw, i.e. not base64 encoded, consent string bytes
func (c *Consent) Bytes() []byte {
	var b bitWriter
	b.AppendByte(c.Version, 6)
	b.AppendInt(c.Created.UnixNano()/int64(time.Second/10), 36)
	b.AppendInt(c.LastUpdated.UnixNano()/int64(time.Second/10), 36)
	b.AppendInt(int64(c.CmpID), 12)
	b.AppendInt(int64(c.CmpVersion), 12)
	b.AppendByte(c.ConsentScreen, 6)
	lang := []byte(strings.ToUpper(c.ConsentLanguage))
	if len(lang) == 2 {
		b.AppendByte(lang[0]-byte('A'), 6)
		b.AppendByte(lang[1]-byte('A'), 6)
	} else {
		b.AppendByte(byte('X')-byte('A'), 6)
		b.AppendByte(byte('X')-byte('A'), 6)
	}
	b.AppendInt(int64(c.VendorListVersion), 12)
	b.AppendBools(c.PurposesAllowed[:])
	b.AppendInt(int64(c.MaxVendorID), 16)

	switch ecType, defConsent, rngCount := c.findSmallest(); ecType {
	case bitFieldType:
		b.AppendByte(byte(bitFieldType), 1) // encoding type
		for i := 1; i <= c.MaxVendorID; i++ {
			var v byte
			if c.Vendors[i] {
				v = 1
			}
			b.AppendByte(v, 1)
		}
	case rangeType:
		b.AppendByte(byte(rangeType), 1) // encoding type
		if defConsent {
			b.AppendByte(1, 1)
		} else {
			b.AppendByte(0, 1)
		}

		b.AppendInt(rngCount, 12)
		for i := 1; i <= c.MaxVendorID; i++ {
			start := i
			for ; start <= c.MaxVendorID && c.Vendors[start] == defConsent; start++ {
			}
			end := start
			for ; end <= c.MaxVendorID && c.Vendors[end] != defConsent; end++ {
			}

			if end == start {
				break
			} else if end-start == 1 {
				b.AppendByte(0, 1)
				b.AppendInt(int64(start), 16)
			} else {
				b.AppendByte(1, 1)
				b.AppendInt(int64(start), 16)
				b.AppendInt(int64(end-1), 16)
			}
			i = end
		}
	}

	return b.Bytes()
}

func (c *Consent) findSmallest() (encodingType, bool, int64) {
	var bfScore, rtScore, rfScore int
	// bitfield, range with default consent == true, and range with false
	var rtRecords, rfRecords int64

	bfScore = c.MaxVendorID
	rtScore = 1 + 12
	rfScore = 1 + 12

	for i := 1; i <= c.MaxVendorID; i++ {
		cur := c.Vendors[i]
		start := i
		for ; i <= c.MaxVendorID && c.Vendors[i] == cur; i++ {
		}
		end := i
		i--

		size := 1
		if end-start == 1 {
			size += 16
		} else {
			size += 16 + 16
		}

		if cur {
			rfScore += size
			rfRecords++
		} else {
			rtScore += size
			rtRecords++
		}
	}

	min := bfScore
	if min > rtScore {
		min = rtScore
	}
	if min > rfScore {
		min = rfScore
	}

	if min == rfScore {
		return rangeType, false, rfRecords
	} else if min == rtScore {
		return rangeType, true, rtRecords
	} else {
		return bitFieldType, false, 0
	}
}

func (c *Consent) outputRange(defConsent bool) (int64, bitWriter) {
	var b bitWriter
	var count int64
	return count, b
}

// String return a base64-encoded consent string
func (c *Consent) String() string {
	return base64.RawURLEncoding.EncodeToString(c.Bytes())
}

// ParseRaw converts a raw, i.e. non base64-encoded, consent string into the
// Consent struct
func (c *Consent) ParseRaw(binary []byte) error {
	if len(binary) < 21 {
		return ErrUnexpectedEnd
	}

	b := newBitReader(binary)

	c.Version, _ = b.ReadByte(6)
	if c.Version != 1 {
		return ErrUnsupported
	}

	dt, _ := b.ReadInt(36)
	c.Created = time.Unix(dt/10, dt%10*100*1000*1000)

	dt, _ = b.ReadInt(36)
	c.LastUpdated = time.Unix(dt/10, dt%10*100*1000*1000)

	dt, _ = b.ReadInt(12)
	c.CmpID = int(dt)

	dt, _ = b.ReadInt(12)
	c.CmpVersion = int(dt)

	c.ConsentScreen, _ = b.ReadByte(6)

	l1, _ := b.ReadByte(6)
	l2, _ := b.ReadByte(6)
	c.ConsentLanguage = string([]byte{l1 + byte('A'), l2 + byte('A')})

	dt, _ = b.ReadInt(12)
	c.VendorListVersion = int(dt)

	for i := range c.PurposesAllowed {
		by, _ := b.ReadByte(1)
		v := false
		if by == 1 {
			v = true
		}
		c.PurposesAllowed[i] = v
	}

	dt, _ = b.ReadInt(16)
	c.MaxVendorID = int(dt)

	by, _ := b.ReadByte(1)
	c.encodingType = encodingType(by)
	c.Vendors = make(map[int]bool)
	switch c.encodingType {
	case bitFieldType:
		for i := 1; i <= c.MaxVendorID; i++ {
			if by, ok := b.ReadByte(1); !ok {
				return ErrUnexpectedEnd
			} else if by == 1 {
				c.Vendors[i] = true
			}
		}
	case rangeType:
		by, ok := b.ReadByte(1)
		if !ok {
			return ErrUnexpectedEnd
		}
		defCons := false
		if by == 1 {
			defCons = true

			for i := 1; i <= int(c.MaxVendorID); i++ {
				c.Vendors[i] = true
			}
		}

		numEntries, ok := b.ReadInt(12)
		if !ok {
			return ErrUnexpectedEnd
		}
		for i := 0; i < int(numEntries); i++ {
			singleOrRange, ok := b.ReadByte(1)
			if !ok {
				return ErrUnexpectedEnd
			}
			if singleOrRange == 0 { // Single
				dt, ok = b.ReadInt(16)
				if !ok {
					return ErrUnexpectedEnd
				}
				c.Vendors[int(dt)] = !defCons
			} else { // Range
				start, ok := b.ReadInt(16)
				if !ok {
					return ErrUnexpectedEnd
				}
				end, ok := b.ReadInt(16)
				if !ok {
					return ErrUnexpectedEnd
				}
				for j := start; j <= end; j++ {
					if defCons {
						delete(c.Vendors, int(j))
					} else {
						c.Vendors[int(j)] = true
					}
				}
			}
		}
	}

	return nil
}

// Parse parses base64 encoded consent
func (c *Consent) Parse(data string) error {
	bytes := []byte(data)
	bin := make([]byte, base64.RawStdEncoding.DecodedLen(len(bytes)))
	_, err := base64.RawURLEncoding.Decode(bin, bytes)
	if err != nil {
		return err
	}
	return c.ParseRaw(bin)
}

// Parse parses base64 encoded consent
func Parse(data string) (*Consent, error) {
	c := new(Consent)
	return c, c.Parse(data)
}
