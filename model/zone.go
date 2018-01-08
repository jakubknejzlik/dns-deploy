package model

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ghodss/yaml"
)

type DomainZone struct {
	Domain  Domain         `yaml:"domain,omitempty"`
	Records []DomainRecord `yaml:"records"`
}

type DomainZoneDiff struct {
	KeepRecords   []DomainRecord `yaml:"keep"`
	UpdateRecords []DomainRecord `yaml:"update"`
	AddRecords    []DomainRecord `yaml:"add"`
	DeleteRecords []DomainRecord `yaml:"delete"`
}

func NewDomainZoneFromFile(filepath string) (DomainZone, error) {
	var zone DomainZone

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return zone, err
	}

	fileExt := path.Ext(filepath)
	zone.Domain = Domain{Name: path.Base(strings.TrimSuffix(filepath, fileExt))}

	switch fileExt {
	case ".json":
		if err := json.Unmarshal(data, &zone); err != nil {
			return zone, err
		}
		break
	case ".yaml":
	case ".yml":
		if err := yaml.Unmarshal(data, &zone); err != nil {
			return zone, err
		}
		break
	}
	return zone, nil
}

func CreateDomainZoneDiff(from, to DomainZone) DomainZoneDiff {
	diff := DomainZoneDiff{}

	// find missing records
	for _, toRecord := range to.Records {
		exists := false
		for _, fromRecord := range from.Records {
			if toRecord.isTypeNameEqual(fromRecord) {
				exists = true
				break
			}
		}
		if !exists {
			diff.AddRecords = append(diff.AddRecords, toRecord)
		}
	}

	for _, fromRecord := range from.Records {
		found := false
		for _, toRecord := range to.Records {
			if fromRecord.isTypeNameEqual(toRecord) {
				found = true

				if fromRecord.isEqual(toRecord) {
					diff.KeepRecords = append(diff.KeepRecords, toRecord)
				} else {
					toRecord.ID = fromRecord.ID
					diff.UpdateRecords = append(diff.UpdateRecords, toRecord)
				}
			}
		}
		if !found && fromRecord.Type != "NS" {
			diff.DeleteRecords = append(diff.DeleteRecords, fromRecord)
		}
	}
	return diff
}
