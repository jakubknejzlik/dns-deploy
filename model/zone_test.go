package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDomainZoneFromFile(t *testing.T) {

	exptectedZone := DomainZone{
		Domain: Domain{Name: "example.com"},
		Records: []DomainRecord{
			DomainRecord{
				Type: "A",
				Data: "8.8.8.8",
			},
		},
	}

	zone, err := NewDomainZoneFromFile("../test/example.com.yml")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	assert.Equal(t, exptectedZone, zone, "new zone should be equal to expected zone")
}

func TestCreateDomainZoneDiff(t *testing.T) {

	expectedDiff := DomainZoneDiff{
		KeepRecords:   []DomainRecord{DomainRecord{Name: "test2", Type: "CNAME", Data: "test2.diff2.com."}},
		UpdateRecords: []DomainRecord{DomainRecord{Name: "test", Type: "CNAME", Data: "test.diff2.com."}},
		AddRecords:    []DomainRecord{DomainRecord{Name: "@", Type: "AAAA", Data: "8.8.8.8"}},
		DeleteRecords: []DomainRecord{DomainRecord{Name: "@", Type: "A", Data: "8.8.8.8"}},
	}

	zone1, _ := NewDomainZoneFromFile("../test/diff.com.yml")
	zone2, _ := NewDomainZoneFromFile("../test/diff2.com.yml")

	diff := CreateDomainZoneDiff(zone1, zone2)

	assert.Equal(t, expectedDiff.KeepRecords, diff.KeepRecords, "keep records should be same")
	assert.Equal(t, expectedDiff.UpdateRecords, diff.UpdateRecords, "update records should be same")
	assert.Equal(t, expectedDiff.AddRecords, diff.AddRecords, "add records should be same")
	assert.Equal(t, expectedDiff.DeleteRecords, diff.DeleteRecords, "delete records should be same")

}
