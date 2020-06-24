package consent

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
)


var ErrInvalidPubRestrictionType = errors.New("consent: invalid pub restriction type")

type ConsentV2 struct {
	Created           time.Time
	LastUpdated       time.Time
	CmpID             int
	CmpVersion        int
	ConsentScreen     byte
	ConsentLanguage   [2]byte
	VendorListVersion int

	TcfPolicyVersion  byte
	IsServiceSpecific bool
	UseNonStandardStacks bool
	SpecialFeatureOptIns [12]bool
	PurposesConsent [24]bool // used to be called PurposesAllowed
	PurposesLITransparency [24]bool

	// specific jurisdiction disclosures
	PurposeOneTreatment bool
	PublisherCC [2]byte // ISO 3166-1 alpha-2 code, encoded with 6 bits per char

	// vendor consent section
	VendorConsent VendorSet

	// Vendor legitimate interest section
	VendorLegitimateInterest VendorSet

	// publisher restrictions section
	PubRestrictions []PubRestriction

	// optional:
	DisclosedVendors VendorSet
	AllowedVendors VendorSet
	PublisherTC *PublisherTC
}

type PubRestriction struct {
	PurposeID int
	RestrictionType byte
	Vendors VendorSet
}

func (c *ConsentV2) Version() byte {
	return 2
}

type VendorSet struct {
	maxVendorID int
	Set map[int]bool
}

func (c *ConsentV2) String() string {
	var b bitWriter
	c.writeCoreString(&b)
	coreBytes := b.Bytes()
	totalEncodedByteLen := base64.RawURLEncoding.EncodedLen(len(coreBytes))

	var disclosedVendorBytes, allowedVendorBytes, publisherTCBytes []byte
	if len(c.DisclosedVendors.Set) != 0 {
		var b bitWriter
		b.AppendByte(disclosedVendorsType, 3)
		c.DisclosedVendors.AppendRangeOrBitField(&b)
		disclosedVendorBytes = b.Bytes()
		totalEncodedByteLen += base64.RawURLEncoding.EncodedLen(len(disclosedVendorBytes)) + 1
	}
	if len(c.AllowedVendors.Set) != 0 {
		var b bitWriter
		b.AppendByte(allowedVendorsType, 3)
		c.AllowedVendors.AppendRangeOrBitField(&b)
		allowedVendorBytes = b.Bytes()
		totalEncodedByteLen += base64.RawURLEncoding.EncodedLen(len(allowedVendorBytes)) + 1
	}

	if c.PublisherTC != nil {
		var b bitWriter
		b.AppendByte(publisherTCType, 3)
		c.PublisherTC.Write(&b)
		publisherTCBytes = b.Bytes()
		totalEncodedByteLen += base64.RawURLEncoding.EncodedLen(len(publisherTCBytes)) + 1
	}

	dst := make([]byte, totalEncodedByteLen)
	base64.RawURLEncoding.Encode(dst, coreBytes)
	offset := base64.RawURLEncoding.EncodedLen((len(coreBytes)))
	if len(disclosedVendorBytes) > 0 {
		dst[offset] = '.'
		offset++
		base64.RawURLEncoding.Encode(dst[offset:], disclosedVendorBytes)
		offset += base64.RawURLEncoding.EncodedLen(len(disclosedVendorBytes))
	}
	if len(allowedVendorBytes) > 0 {
		dst[offset] = '.'
		offset += 1
		base64.RawURLEncoding.Encode(dst[offset:], allowedVendorBytes)
		offset += base64.RawURLEncoding.EncodedLen(len(allowedVendorBytes))
	}
	if len(publisherTCBytes) > 0 {
		dst[offset] = '.'
		offset += 1
		base64.RawURLEncoding.Encode(dst[offset:], publisherTCBytes)
		offset += base64.RawURLEncoding.EncodedLen(len(publisherTCBytes))
	}

	return string(dst)
}

func (c *ConsentV2) writeCoreString(b *bitWriter) {
	b.AppendByte(c.Version(), 6)
	b.AppendInt(c.Created.UnixNano()/int64(time.Second/10), 36)
	b.AppendInt(c.LastUpdated.UnixNano()/int64(time.Second/10), 36)
	b.AppendInt(int64(c.CmpID), 12)
	b.AppendInt(int64(c.CmpVersion), 12)
	b.AppendByte(c.ConsentScreen, 6)
	b.AppendByte(c.ConsentLanguage[0]-byte('A'), 6)
	b.AppendByte(c.ConsentLanguage[1]-byte('A'), 6)
	b.AppendInt(int64(c.VendorListVersion), 12)
	b.AppendByte(c.TcfPolicyVersion, 6)
	b.AppendBit(c.IsServiceSpecific)
	b.AppendBit(c.UseNonStandardStacks)

	for _, optIn := range c.SpecialFeatureOptIns {
		b.AppendBit(optIn)
	}

	for _, consent := range c.PurposesConsent {
		b.AppendBit(consent)
	}

	for _, tr := range c.PurposesLITransparency {
		b.AppendBit(tr)
	}

	b.AppendBit(c.PurposeOneTreatment)

	b.AppendByte(c.PublisherCC[0]-byte('A'), 6)
	b.AppendByte(c.PublisherCC[1]-byte('A'), 6)

	c.VendorConsent.AppendRangeOrBitField(b)
	c.VendorLegitimateInterest.AppendRangeOrBitField(b)

	b.AppendInt(int64(len(c.PubRestrictions)), 12)
	for _, pubRestriction := range c.PubRestrictions {
		b.AppendInt(int64(pubRestriction.PurposeID), 6)
		b.AppendByte(pubRestriction.RestrictionType, 2)
		numEntries, _ := pubRestriction.Vendors.getRangeSizes()
		AppendRange(b, pubRestriction.Vendors.Set, pubRestriction.Vendors.maxVendorID, numEntries)
	}
}

func (c *ConsentV2) ParseCore(binary []byte) error {
	if len(binary) * 8 < 213 {
		return ErrUnexpectedEnd
	}
	b := newBitReader(binary)
	version, _ := b.ReadByte(6)
	if version != 2 {
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
	c.ConsentLanguage = [2]byte{l1 + byte('A'), l2 + byte('A')}

	dt, _ = b.ReadInt(12)
	c.VendorListVersion = int(dt)

	c.TcfPolicyVersion, _ = b.ReadByte(6)
	c.IsServiceSpecific, _ = b.ReadBit()
	c.UseNonStandardStacks, _ = b.ReadBit()
	for i:= 0; i < 12; i++ {
		c.SpecialFeatureOptIns[i], _ = b.ReadBit()
	}
	for i:= 0; i < 24; i++ {
		c.PurposesConsent[i], _ = b.ReadBit()
	}
	for i:= 0; i < 24; i++ {
		c.PurposesLITransparency[i], _ = b.ReadBit()
	}
	c.PurposeOneTreatment, _ = b.ReadBit()
	l1, _ = b.ReadByte(6)
	l2, _ = b.ReadByte(6)
	c.PublisherCC = [2]byte{l1 + byte('A'), l2 + byte('A')}

	var err error
	c.VendorConsent, err = parseVendorSet(&b)
	if err != nil {
		return err
	}

	c.VendorLegitimateInterest, err = parseVendorSet(&b)
	if err != nil {
		return err
	}

	c.PubRestrictions, err = parsePubRestrictions(&b)
	if err != nil {
		return err
	}

	return nil
}

const (
	disclosedVendorsType byte = 1
	allowedVendorsType byte = 2
	publisherTCType byte = 3
)

type PublisherTC struct {
	PubPurposesConsent [24]bool // true: consent, false: no consent
	PubPurposesLITransparency [24]bool // true: legitimate interest established
	CustomPurposesConsent []bool // consent/no consent
	CustomPurposesLITransparency []bool // legitimate interest established/not established
}

func (p *PublisherTC) Write(w *bitWriter) {
	for _, b := range(p.PubPurposesConsent) {
		w.AppendBit(b)
	}
	for _, b := range(p.PubPurposesLITransparency) {
		w.AppendBit(b)
	}
	w.AppendInt(int64(len(p.CustomPurposesConsent)), 6)
	for _, b := range(p.CustomPurposesConsent) {
		w.AppendBit(b)
	}

	for _, b := range(p.CustomPurposesLITransparency) {
		w.AppendBit(b)
	}
}

func(c *ConsentV2) ParseNonCoreSegment(binary []byte) error {
	b := newBitReader(binary)
	segmentType, ok := b.ReadByte(3)
	if !ok {
		return ErrUnexpectedEnd
	}
	switch segmentType {
	case disclosedVendorsType:
		return c.parseDisclosedVendors(&b)
	case allowedVendorsType:
		return c.parseAllowedVendors(&b)
	case publisherTCType:
		return c.parsePublisherTC(&b)
	default:
		return nil
	}
}

// assumes segment type bits have already been consumed
func (c *ConsentV2) parseDisclosedVendors(b *bitReader) error {
	m, err := parseVendorSet(b)
	if err != nil {
		return err
	}
	c.DisclosedVendors = m
	return nil
}

// assumes segment type bits have already been consumed
func (c *ConsentV2) parseAllowedVendors(b *bitReader) error {
	m, err := parseVendorSet(b)
	if err != nil {
		return err
	}
	c.AllowedVendors = m
	return nil
}

// assumes segment type bits have already been consumed
func (c *ConsentV2) parsePublisherTC(b *bitReader) error {
	publisherTC := &PublisherTC{}
	var ok bool
	for i := range publisherTC.PubPurposesConsent {
		publisherTC.PubPurposesConsent[i], ok = b.ReadBit()
		if !ok {
			return ErrUnexpectedEnd
		}
	}

	for i := range publisherTC.PubPurposesLITransparency {
		publisherTC.PubPurposesLITransparency[i], ok = b.ReadBit()
		if !ok {
			return ErrUnexpectedEnd
		}
	}

	numCustomPurposes, ok := b.ReadInt(6)
	if !ok {
		return ErrUnexpectedEnd
	}

	customPurposesConsent := make([]bool, numCustomPurposes)
	for i := 0; i < int(numCustomPurposes); i++ {
		customPurposesConsent[i], ok = b.ReadBit()
		if !ok {
			return ErrUnexpectedEnd
		}
	}

	customPurposesLITransparency := make([]bool, numCustomPurposes)
	for i := 0; i < int(numCustomPurposes); i++ {
		customPurposesLITransparency[i], ok = b.ReadBit()
		if !ok {
			return ErrUnexpectedEnd
		}
	}
	publisherTC.CustomPurposesConsent = customPurposesConsent
	publisherTC.CustomPurposesLITransparency = customPurposesLITransparency
	c.PublisherTC = publisherTC
	return nil
}

func ParseV2(s string) (*ConsentV2, error) {
	c := new(ConsentV2)
	return c, c.Parse(s)
}

func (c *ConsentV2) Parse(data string) error {
	segments := strings.Split(data, ".")
	if len(segments) == 0 {
		return ErrUnexpectedEnd
	}
	coreSegment := segments[0]
	bytes := []byte(coreSegment)
	bin := make([]byte, base64.RawStdEncoding.DecodedLen(len(bytes)))
	_, err := base64.RawURLEncoding.Decode(bin, bytes)
	if err != nil {
		return err
	}

	err = c.ParseCore(bin)
	if err != nil {
		return err
	}

	for i := 1; i < len(segments); i++ {
		segment := segments[i]
		bytes := []byte(segment)
		bin := make([]byte, base64.RawStdEncoding.DecodedLen(len(bytes)))
		_, err := base64.RawURLEncoding.Decode(bin, bytes)
		if err != nil {
			return nil
		}
		err = c.ParseNonCoreSegment(bin)
		if err != nil {
			return nil
		}
	}
	return nil
}

const (
	RestrictionTypeNotAllowed byte = 0
	RestrictionTypeRequireConsent byte = 1
	RestrictionTypeRequireLegitimateInterest byte = 2
)


func parsePubRestrictions(b *bitReader) ([]PubRestriction, error) {
	numPubRestrictions, ok := b.ReadInt(12)
	if !ok {
		return nil, ErrUnexpectedEnd
	}
	restrictions := make([]PubRestriction, 0, int(numPubRestrictions))
	for i := 0; i < int(numPubRestrictions); i++ {
		pr := PubRestriction{}
		purposeID, ok := b.ReadInt(6)
		if !ok {
			return nil, ErrUnexpectedEnd
		}
		pr.PurposeID = int(purposeID)
		pr.RestrictionType, ok = b.ReadByte(2)
		if !ok {
			return nil, ErrUnexpectedEnd
		}
		if pr.RestrictionType > 2 {
			return []PubRestriction{}, ErrInvalidPubRestrictionType
		}
		vendors, maxVendorID, err := parseRange(b) // NOTE: this one can't be a bit field for some reason
		if err != nil {
			return nil, err
		}
		pr.Vendors = VendorSet{Set: vendors, maxVendorID: maxVendorID}
		restrictions = append(restrictions, pr)
	}
	return restrictions, nil
}

func parseVendorSet(b *bitReader) (VendorSet, error) {
	maxVendorID, ok := b.ReadInt(16)
	if !ok {
		return VendorSet{}, ErrUnexpectedEnd
	}
	isRangeEncoding, ok := b.ReadBit()
	if !ok {
		return VendorSet{}, ErrUnexpectedEnd
	}
	var err error
	var set map[int]bool
	if isRangeEncoding {
		set, _, err = parseRange(b)
	} else {
		set, err = parseBitField(b, int(maxVendorID))
	}

	if err != nil {
		return VendorSet{}, err
	}

	return VendorSet{Set: set, maxVendorID: int(maxVendorID)}, nil
}

func parseRange(b *bitReader) (map[int]bool, int, error) { // second return value is max vendor id in the set
	numEntries, ok := b.ReadInt(12)
	var maxVendorID int
	if !ok {
		return nil, 0, ErrUnexpectedEnd
	}
	vendors := make(map[int]bool)
	for i := 0; i < int(numEntries); i++ {
		isARange, _ := b.ReadBit()
		startOrOnlyVendorID, ok := b.ReadInt(16)
		if !ok {
			return nil, 0, ErrUnexpectedEnd
		}
		endVendorID := startOrOnlyVendorID
		if isARange {
			endVendorID, ok = b.ReadInt(16)
			if !ok {
				return nil, 0, ErrUnexpectedEnd
			}
		}
		for id := startOrOnlyVendorID; id <= endVendorID; id++ {
			vendors[int(id)] = true
			maxVendorID = int(id)
		}
	}

	return vendors, maxVendorID, nil
}

func parseBitField(b *bitReader, maxVendorID int) (map[int]bool, error) {
	vendors := map[int]bool{}
	for vendorID := 1; vendorID <= maxVendorID; vendorID++ {
		vendorBit, ok := b.ReadBit()
		if !ok {
			return nil, ErrUnexpectedEnd
		}
		if vendorBit {
			vendors[vendorID] = true
		}
	}
	return vendors, nil
}

func (s *VendorSet) AppendRangeOrBitField(b *bitWriter) {
	numRangeEntries, rangeSizeInBits := s.getRangeSizes()
	isRange := false
	if rangeSizeInBits < s.maxVendorID {
		isRange = true
	}

	b.AppendInt(int64(s.maxVendorID), 16)
	if isRange {
		b.AppendByte(1, 1) // encoding type
		AppendRange(b, s.Set, s.maxVendorID, numRangeEntries)
	} else {
		b.AppendByte(0, 1) // encoding type
		AppendBitField(b, s.Set, s.maxVendorID)
	}
}

func AppendBitField(b *bitWriter, vendorSet map[int]bool, maxVendorID int) {
	for vendorID := 1; vendorID <= maxVendorID; vendorID++ {
		var bit byte = 0
		if vendorSet[vendorID] {
			bit = 1
		}
		b.AppendByte(bit, 1)
	}
}

func (v *VendorSet) getRangeSizes() (int, int) { // numEntries, number of bits
	var numEntries, bitCount int
	for vendorID := 1; vendorID <= v.maxVendorID; vendorID++ {
		if v.Set[vendorID] {
			rangeStart := vendorID
			rangeEnd := vendorID
			for i := rangeStart + 1; i <= v.maxVendorID && v.Set[i]; i++ {
				rangeEnd = i
			}
			if rangeStart == rangeEnd { // single entry
				bitCount += 17
			} else { // range
				bitCount += 33
			}

			numEntries += 1
			vendorID = rangeEnd + 1
		}
	}
	return numEntries, bitCount + 12
}

func AppendRange(b *bitWriter, vendorSet map[int]bool, maxVendorID int, numEntries int) {
	b.AppendInt(int64(numEntries), 12)
	for vendorID := 1; vendorID <= maxVendorID; vendorID++ {
		if vendorSet[vendorID] {
			rangeStart := vendorID
			rangeEnd := vendorID
			for i := rangeStart + 1; i <= maxVendorID && vendorSet[i]; i++ {
				rangeEnd = i
			}
			if rangeStart == rangeEnd { // single entry
				b.AppendByte(0, 1)// isARange
				b.AppendInt(int64(rangeStart), 16) // only vendor id
			} else { // range
				b.AppendByte(1, 1)// isARange
				b.AppendInt(int64(rangeStart), 16) // first vendor id
				b.AppendInt(int64(rangeEnd), 16) // last vendor id
			}

			vendorID = rangeEnd + 1
		}
	}
}
