package providers

import (
	"context"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/jakubknejzlik/dns-deploy/model"
	"golang.org/x/oauth2"
)

// DigitalOceanProvider ...
type DigitalOceanProvider struct {
	client *godo.Client
}

func NewDigitalOceanClient(token string) DNSProvider {
	tokenSource := &TokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	return DigitalOceanProvider{client}
}

func (p DigitalOceanProvider) ListDomains() ([]model.Domain, error) {

	domains := []model.Domain{}

	_domains, _, err := p.client.Domains.List(context.Background(), nil)
	if err != nil {
		return domains, err
	}

	for _, d := range _domains {
		domains = append(domains, model.Domain{Name: d.Name})
	}

	return domains, nil
}

func (p DigitalOceanProvider) GetDomain(name string) (model.Domain, error) {
	var domain model.Domain

	_domain, _, err := p.client.Domains.Get(context.Background(), name)
	if err != nil {
		return domain, err
	}

	domain = model.Domain{Name: _domain.Name}

	return domain, nil
}

func (p DigitalOceanProvider) CreateDomain(name string) (model.Domain, error) {
	var domain model.Domain

	req := &godo.DomainCreateRequest{Name: name}
	_domain, _, err := p.client.Domains.Create(context.Background(), req)
	if err != nil {
		return domain, err
	}

	domain = model.Domain{Name: _domain.Name}

	return domain, nil
}
func (p DigitalOceanProvider) DeleteDomain(name string) error {
	_, err := p.client.Domains.Delete(context.Background(), name)
	return err
}

func (p DigitalOceanProvider) ListDomainRecords(domain string) ([]model.DomainRecord, error) {
	records := []model.DomainRecord{}

	opts := &godo.ListOptions{PerPage: 999}
	_records, _, err := p.client.Domains.Records(context.Background(), domain, opts)
	if err != nil {
		return records, err
	}

	for _, record := range _records {
		rec := model.DomainRecord{
			ID:       strconv.Itoa(record.ID),
			Type:     record.Type,
			Name:     record.Name,
			Data:     record.Data,
			Priority: record.Priority,
			Port:     record.Port,
			TTL:      record.TTL,
			Weight:   record.Weight,
			Flags:    record.Flags,
			Tag:      record.Tag,
		}
		records = append(records, rec)
	}

	return records, nil
}

func createEditRequest(record model.DomainRecord) *godo.DomainRecordEditRequest {
	data := record.Data

	if record.Type == "CNAME" {
		data += "."
	}

	return &godo.DomainRecordEditRequest{
		Type:     record.Type,
		Name:     record.Name,
		Data:     data,
		Priority: record.Priority,
		Port:     record.Port,
		TTL:      record.TTL,
		Weight:   record.Weight,
		Flags:    record.Flags,
		Tag:      record.Tag,
	}

}
func (p DigitalOceanProvider) CreateDomainRecord(domain string, record model.DomainRecord) error {
	req := createEditRequest(record)
	_, _, err := p.client.Domains.CreateRecord(context.Background(), domain, req)
	return err
}

func (p DigitalOceanProvider) UpdateDomainRecord(domain string, record model.DomainRecord) error {
	recordID, err := strconv.ParseInt(record.ID, 10, 32)
	if err != nil {
		return err
	}
	req := createEditRequest(record)
	_, _, err = p.client.Domains.EditRecord(context.Background(), domain, int(recordID), req)
	return err
}

func (p DigitalOceanProvider) DeleteDomainRecord(domain string, recordID string) error {
	id, err := strconv.ParseInt(recordID, 10, 12)
	if err != nil {
		return err
	}
	_, err = p.client.Domains.DeleteRecord(context.Background(), domain, int(id))
	return err
}
