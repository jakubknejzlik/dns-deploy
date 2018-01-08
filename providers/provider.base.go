package providers

import (
	"github.com/jakubknejzlik/dns-deploy/model"
	"golang.org/x/oauth2"
)

// DNSProvider ...
type DNSProvider interface {
	ListDomains() ([]model.Domain, error)
	GetDomain(name string) (model.Domain, error)
	CreateDomain(name string) (model.Domain, error)
	DeleteDomain(name string) error

	// GetDomainZone(domain string) (model.DomainZone, error)
	// UpdateDomainZone(domain string, diff model.DomainZoneDiff) error
	ListDomainRecords(domain string) ([]model.DomainRecord, error)
	// GetDomainRecord(domain, recordID string) (model.DomainRecord, error)
	CreateDomainRecord(domain string, record model.DomainRecord) error
	UpdateDomainRecord(domain string, record model.DomainRecord) error
	DeleteDomainRecord(domain, recordID string) error
}

func GetProvider(code, token string) (DNSProvider, error) {
	return NewDigitalOceanClient(token)
}

// TokenSource ...
type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
