package consent

import (
	"testing"
)

func TestConsentV2Parse(t *testing.T) {
	s := "COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA"
	c := ConsentV2{}

	err := c.Parse(s)
	equal(t, err, nil)
	equal(t, int(c.Version()), 2)
	equal(t, c.Created.String(), "2020-01-26 18:01:00 +0100 CET")
	equal(t, c.LastUpdated.String(), "2021-02-02 18:01:00 +0100 CET")
	equal(t, c.CmpID, 675)
	equal(t, c.CmpVersion, 2)
	equal(t, int(c.ConsentScreen), 1)
	equal(t, string(c.ConsentLanguage[:]), "EN")
	equal(t, c.VendorListVersion, 15)
	equal(t, int(c.TcfPolicyVersion), 2)
	equal(t, c.IsServiceSpecific, false)
	equal(t, c.UseNonStandardStacks, false)
	equal(t, string(c.PublisherCC[:]), "AA")
	equal(t, c.PurposeOneTreatment, true)
	equal(t, c.SpecialFeatureOptIns, [12]bool{true, false})
	equal(t, c.PurposesConsent, [24]bool{false, true, false, false, false, false, false, false, false, true})
	equal(t, c.PurposesLITransparency, [24]bool{false, true, false, false, false, false, false, false, true})

	equal(t, c.String(), s)
}


func TestConsentV2ParseVendorLegitimateInterest(t *testing.T) {
	s := "COrEAV4OrXx94ACABBENAHCIAD-AAAAAAACAAxAAAAgAIAwgAgAAAAEAgQAAAAAEAYQAQAAAACAAAABAAA"
	c := ConsentV2{}

	err := c.Parse(s)
	equal(t, err, nil)
	equal(t, c.VendorLegitimateInterest.Set, map[int]bool{
		37: true,
		47: true,
		48: true,
		53: true,
		65: true,
		98: true,
		129: true,
	})

	equal(t, c.String(), s)
}

func TestConsentV2ParseDisclosedAndAllowedVendors(t *testing.T) {
	s := "COrEAV4OrXx94ACABBENAHCIAD-AAAAAAACAAxAAAAgAIAwgAgAAAAEAgQAAAAAEAYQAQAAAACAAAABAAA.IBAgAAAgAIAwgAgAAAAEAAAACA.QAagAQAgAIAwgA"
	c := ConsentV2{}
	err := c.Parse(s)
	equal(t, err, nil)
	equal(t, c.AllowedVendors.Set, map[int]bool{
		12: true,
		23: true,
		37: true,
		47: true,
		48: true,
		53: true,
	})

	equal(t, c.DisclosedVendors.Set, map[int]bool{
		23: true,
		37: true,
		47: true,
		48: true,
		53: true,
		65: true,
		98: true,
		129: true,
	})

	equal(t, c.String(), s)
}

func TestConsentV2ParsePubRestrictions(t *testing.T) {
	s :=  "COuQACgOuQACgM-AAAENAPCAAAAAAAAAAAAAAAAAAABgoAAQAAHAAA"
	c := ConsentV2{}
	err := c.Parse(s)
	equal(t, err, nil)
	expectedRestrictions := []PubRestriction{
		{PurposeID: 1, RestrictionType: RestrictionTypeRequireConsent, Vendors: VendorSet{}},
		{PurposeID: 2, RestrictionType: RestrictionTypeNotAllowed, Vendors: VendorSet{}},
		{PurposeID: 3, RestrictionType: RestrictionTypeRequireLegitimateInterest, Vendors: VendorSet{}},
	}

	for i := range expectedRestrictions {
		equal(t, c.PubRestrictions[i].PurposeID, expectedRestrictions[i].PurposeID)
		equal(t, c.PubRestrictions[i].RestrictionType, expectedRestrictions[i].RestrictionType)
	}
	equal(t, c.String(), s)

	vs := VendorSet{Set: map[int]bool{1: true, 3: true, 4: true}, maxVendorID: 4}
	c.PubRestrictions[0].Vendors = vs
	s2 := c.String()
	c2 := ConsentV2{}
	err = c2.Parse(s2)
	equal(t, err, nil)
	expectedRestrictions[0].Vendors = vs
	for i := range expectedRestrictions {
		equal(t, c2.PubRestrictions[i].PurposeID, expectedRestrictions[i].PurposeID)
		equal(t, c2.PubRestrictions[i].RestrictionType, expectedRestrictions[i].RestrictionType)
	}
	equal(t, c2.String(), s2)
}


func TestConsentV2ParsePublisherTC(t *testing.T) {
	s := "COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA.cAAAAAAAITg"
	c := ConsentV2{}
	err := c.Parse(s)
	equal(t, err, nil)

	equal(t, c.PublisherTC.PubPurposesConsent, [24]bool{true})
	expectedPPLIT := [24]bool{}
	expectedPPLIT[23] = true
	equal(t, c.PublisherTC.PubPurposesLITransparency, expectedPPLIT)
	equal(t, c.PublisherTC.CustomPurposesConsent, []bool{false, true})
	equal(t, c.PublisherTC.CustomPurposesLITransparency, []bool{true, true})

	equal(t, c.String(), s)
}
