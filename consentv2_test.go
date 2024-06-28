package consent

import (
	"testing"
)

func TestConsentV2Parse(t *testing.T) {
	s := "COtybn4PA_zT4KjACBENAPCIAEBAAECAAIAAAAAAAAAA"
	c := ConsentV2{}

	err := c.Parse(s)
	noError(t, err)
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
		37:  true,
		47:  true,
		48:  true,
		53:  true,
		65:  true,
		98:  true,
		129: true,
	})

	equal(t, c.String(), s)
}

func TestConsentV2ParseDisclosedAndAllowedVendors(t *testing.T) {
	s := "COrEAV4OrXx94ACABBENAHCIAD-AAAAAAACAAxAAAAgAIAwgAgAAAAEAgQAAAAAEAYQAQAAAACAAAABAAA.IBAgAAAgAIAwgAgAAAAEAAAACA.QAagAQAgAIAwgA"
	c := ConsentV2{}
	err := c.Parse(s)
	noError(t, err)
	equal(t, c.AllowedVendors.Set, map[int]bool{
		12: true,
		23: true,
		37: true,
		47: true,
		48: true,
		53: true,
	})

	equal(t, c.DisclosedVendors.Set, map[int]bool{
		23:  true,
		37:  true,
		47:  true,
		48:  true,
		53:  true,
		65:  true,
		98:  true,
		129: true,
	})

	equal(t, c.String(), s)
}

func TestConsentV2ParsePubRestrictions(t *testing.T) {
	s := "COuQACgOuQACgM-AAAENAPCAAAAAAAAAAAAAAAAAAABgoAAQAAHAAA"
	c := ConsentV2{}
	err := c.Parse(s)
	noError(t, err)
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
	noError(t, err)
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

func TestConsentV2_2Parse(t *testing.T) {
	s := "COtwGEAPA9dwAKjACBENA6EIAEBAAECAAIAAAAAAAAAA"
	c := ConsentV2{}

	err := c.Parse(s)
	noError(t, err)
	equal(t, int(c.Version()), 2)
	equal(t, c.Created.String(), "2020-01-26 04:00:00 +0400 +04")
	equal(t, c.LastUpdated.String(), "2021-02-02 04:00:00 +0400 +04")
	equal(t, c.CmpID, 675)
	equal(t, c.CmpVersion, 2)
	equal(t, int(c.ConsentScreen), 1)
	equal(t, string(c.ConsentLanguage[:]), "EN")
	equal(t, c.VendorListVersion, 58)
	equal(t, int(c.TcfPolicyVersion), 4)
	equal(t, c.IsServiceSpecific, false)
	equal(t, c.UseNonStandardStacks, false)
	equal(t, string(c.PublisherCC[:]), "AA")
	equal(t, c.PurposeOneTreatment, true)
	equal(t, c.SpecialFeatureOptIns, [12]bool{true, false})
	equal(t, c.PurposesConsent, [24]bool{false, true, false, false, false, false, false, false, false, true})
	equal(t, c.PurposesLITransparency, [24]bool{false, true, false, false, false, false, false, false, true})

	equal(t, c.String(), s)
}

func TestConsentV2_2ParseVendorLegitimateInterest(t *testing.T) {
	s := "COrBaYAOrVMAAACABBENA6EIAD-AAAAAAACAAxQBQALgAlABeADUAMQBAwBQASgAvABqAGIAQIAA"
	c := ConsentV2{}

	err := c.Parse(s)
	equal(t, err, nil)
	equal(t, c.VendorLegitimateInterest.Set, map[int]bool{
		37:  true,
		47:  true,
		53:  true,
		98:  true,
		129: true,
	})

	equal(t, c.String(), s)
}

func TestConsentV2_2ParseDisclosedAndAllowedVendors(t *testing.T) {
	s := "COrBaYAOrVMAAACABBENA6EAAAAAAAAAAACAAAAAAAAA.IBAwBgALgAlABeADUAMQAgQ.QAagAQAgAIAggA"
	c := ConsentV2{}
	err := c.Parse(s)
	noError(t, err)
	equal(t, c.AllowedVendors.Set, map[int]bool{
		12: true,
		23: true,
		37: true,
		47: true,
		53: true,
	})

	equal(t, c.DisclosedVendors.Set, map[int]bool{
		23:  true,
		37:  true,
		47:  true,
		53:  true,
		98:  true,
		129: true,
	})

	equal(t, c.String(), s)
}

func TestConsentV2_2ParsePublisherTC(t *testing.T) {
	s := "CP7GiMAP7GiMAEsACBFRAqEoAP_gAEPgAAwIIzJD_D7NbSNCwHp3aLsEMAhHRtCAQoQgAASBAmABSAKQIBQCgkAQFAygBCACAAAAICZBIQAECAAACUAAQAAAAAAEAEAAAAAIIAAAgAEAAAAIAAACAAgAEAAIAAAUEAAAmAgEAIIASAAAhAAAAACAAAAAAAAAAgAAAAAAAAEAAAAAAAAAAQAAgAAAAAAAABBGZAP4XZraRoWQ8K5BZghAEKKNoQCFAEAACQIEgACQBSBACAUggCAAkUAAQAAAAABASAJAABAAEAAAgAKAAAAAAAgAgAAAABBAAAAAAgAAAAAAAAQAEAAAABAAAAggAAESEIgBBACQAAAAAABAAAAAAAAAAAAEAAAAAAAAAgAAAAAAAAAAAAEAAAAAAAAAIAA.cAAAD_gAAAA"
	c := ConsentV2{}
	err := c.Parse(s)
	equal(t, err, nil)

	equal(t, c.PublisherTC.PubPurposesConsent, [24]bool{true})
	equal(t, c.PublisherTC.PubPurposesLITransparency, [24]bool{false, true, true, true, true, true, true, true, true, true})
	equal(t, c.PublisherTC.CustomPurposesConsent, []bool{})
	equal(t, c.PublisherTC.CustomPurposesLITransparency, []bool{})

	equal(t, c.String(), s)
}

func TestConsentV2_2ParsePubRestrictions(t *testing.T) {
	s := "CP5ti4AP5ti4AAcABBENAmEsAP_gAEPgACiQg1QYwACAAKgAYABoAFYALgAyABwAEEAJwAoABaADIAGgAOgAegBCgCIAIoASQAmABQAClAFoAXIAvAC_AGEAYgAzABogDaAN4AcwBAACGAEYAI4ASsApAClgFaAVwAygBqgDiAHPAO4A7wB4gD9AIMAQiAiYCKAEWAI6ASUAlQBLgCbgE7AKaAVkAuoBfADgAHtAP3Af8CAAEJgIWgRSBFQCNQEiAJPAToAqoBVwC3gFwALzAXsAwABjIDIwGWANAAbGA2UBxYDjQHzAP-AgOBDcCHYEbwJMQS9BMACYUExATGAmQBMqCZoJowTTBNaCbAJvQTgBOSCcwJ1QTsBPMINQg1QKQAFABbADMAH4ARwApQBlADuAIoAR0AkoB7QF5gMEAZGAywB5ID_gI3gSYAl9BMAEwQJhgTFAmOBMmCZgJnATSAmoBNiCbYJuQTeBN8CcME5QTmAnSBOuCdoJ3ATwAnmEFAINQB4KAIAKUAlwBOwEVBAAYA5wGRgPJCQFwAKgAZAA4ACAAGQANAAiABMAChAFsAXAA3gCGAEcAKUAVoA1QB3gD9AJ2AVkBIgC4AGxhoDQADIATgBQAD0AIQARAAmQBaAFwAMwAcwBAACGgFIAUsArQCuAHOAO8AkoBLgCsgHAAP2AioBIgC4AGMgMjAbGIgMwAnACgAHoAQgAiABMAChAFoAXAAzABzAEAAIaAUgBSwCtAK4AaoA5wCSgEuAJ2AVkA4AB-wEiALgAYyA2MVATACEAEwALQAXAAvABmADaAIAARwApABWgDKAHeASUAlwB-wFXALgAXmMgJABCACYAFwALwAZgA2gCAAEcAKQAVoAygBqgDvAJKAfsBVwC4AF5jgAwAFwA0ADnAHcAQgB5I6BGABUADIAHAAQAAyABoAEQAJIATAAoQBbAFwAL4AYgAzABvAEMAI6AUgBSwCtAK4AZQA7wB-gEWAJUAS4AnYBWQC6gIqASIAqoBcADLAGxgOLIAAwA0ADnAeSQgLgANAAyAEwALkAXgBfADEAGYANoAbwBAACOAFIAK0AZQA7wCSgEuAKyAgABCYCRAFXALgJAAgALgDuJQGwAMgAcACIAEwAKAAXAAvgBiADMAG0AQwAjoBSAFKAK0AZQA1QB3gElAJcAU0ArIBdQEJgIqASIAywBxpQAMABcARwA5wB3AF1AP-UgQAAVAAyABwAEAAMgAaABEACYAFAAKQAWgAuABfADEAGYANoAhoBSAFKAK0AZQA1QB3gD9AIsATsArICEwEiAKuAXAA2MBxo.f_wACHwAAAA"
	c := ConsentV2{}
	err := c.Parse(s)
	noError(t, err)
	expectedRestrictions := []PubRestriction{
		{
			PurposeID:       1,
			RestrictionType: RestrictionTypeRequireConsent,
			Vendors: VendorSet{
				554,
				map[int]bool{165: true, 302: true, 315: true, 554: true},
			},
		},
		{
			PurposeID:       2,
			RestrictionType: RestrictionTypeNotAllowed,
			Vendors: VendorSet{
				969,
				map[int]bool{231: true, 803: true, 969: true},
			},
		},
		{
			PurposeID:       2,
			RestrictionType: RestrictionTypeRequireConsent,
			Vendors: VendorSet{
				867,
				map[int]bool{21: true, 25: true, 28: true, 32: true, 50: true, 52: true, 68: true, 76: true, 80: true, 91: true, 92: true, 111: true, 134: true, 142: true, 165: true, 173: true, 213: true, 239: true, 253: true, 315: true, 345: true, 580: true, 736: true, 867: true},
			},
		},
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
	noError(t, err)
	expectedRestrictions[0].Vendors = vs
	for i := range expectedRestrictions {
		equal(t, c2.PubRestrictions[i].PurposeID, expectedRestrictions[i].PurposeID)
		equal(t, c2.PubRestrictions[i].RestrictionType, expectedRestrictions[i].RestrictionType)
	}
	equal(t, c2.String(), s2)
}
