package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/jakubknejzlik/dns-deploy/model"
	"github.com/jakubknejzlik/dns-deploy/providers"
	"github.com/urfave/cli"
)

// DeployCommand ...
func RunCommand() cli.Command {
	return cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "provider",
				EnvVar: "DNS_PROVIDER",
			},
			cli.StringFlag{
				Name:   "token",
				EnvVar: "DNS_PROVIDER_TOKEN",
			},
		},
		Action: func(c *cli.Context) error {

			providerCode := c.String("provider")
			token := c.String("token")

			if providerCode == "" {
				return cli.NewExitError("--provider option or DNS_PROVIDER envvar must be provided", 1)
			}
			if token == "" {
				return cli.NewExitError("--token option or DNS_TOKEN envvar must be provided", 1)
			}

			provider, err := providers.GetProvider(providerCode, token)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := run(provider); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func run(provider providers.DNSProvider) error {
	currentDomains, err := provider.ListDomains()
	if err != nil {
		return err
	}
	currentDomainsMap := map[string]model.Domain{}
	for _, d := range currentDomains {
		currentDomainsMap[d.Name] = d
	}

	fmt.Println("current domains", currentDomains)

	dir, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, info := range dir {
		if !info.IsDir() {
			zone, err := model.NewDomainZoneFromFile(info.Name())
			if err != nil {
				return err
			}
			if currentDomainsMap[zone.Domain.Name].Name == "" {
				provider.CreateDomain(zone.Domain.Name)
			}
			if err := updateDomainZone(provider, zone); err != nil {
				return err
			}
		}
	}
	return nil
}

func getDomainZone(provider providers.DNSProvider, domain string) (model.DomainZone, error) {
	zone := model.DomainZone{}

	records, err := provider.ListDomainRecords(domain)
	if err != nil {
		return zone, err
	}

	zone.Records = records

	return zone, nil

	// records := provider
	// zone := model.DomainZone{}

	// opts := &godo.ListOptions{PerPage: 999}
	// records, _, err := p.client.Domains.Records(context.Background(), domain, opts)
	// if err != nil {
	// 	return zone, err
	// }

	// for _, record := range records {
	// 	rec := model.DomainRecord{
	// 		ID:       string(record.ID),
	// 		Type:     record.Type,
	// 		Name:     record.Name,
	// 		Data:     record.Data,
	// 		Priority: record.Priority,
	// 		Port:     record.Port,
	// 		TTL:      record.TTL,
	// 		Weight:   record.Weight,
	// 		Flags:    record.Flags,
	// 		Tag:      record.Tag,
	// 	}
	// 	zone.Records = append(zone.Records, rec)
	// }

	// return zone, nil
}

func updateDomainZone(provider providers.DNSProvider, zone model.DomainZone) error {
	remoteZone, err := getDomainZone(provider, zone.Domain.Name)
	if err != nil {
		return err
	}
	fmt.Println("Updating zone", zone.Domain.Name)
	diff := model.CreateDomainZoneDiff(remoteZone, zone)
	yml, _ := yaml.Marshal(diff)
	fmt.Println("diff", string(yml))

	return applyDomainZoneDiff(provider, zone.Domain.Name, diff)
}

func applyDomainZoneDiff(p providers.DNSProvider, domain string, diff model.DomainZoneDiff) error {

	for _, record := range diff.AddRecords {
		if err := p.CreateDomainRecord(domain, record); err != nil {
			return err
		}
	}
	for _, record := range diff.UpdateRecords {
		if err := p.UpdateDomainRecord(domain, record); err != nil {
			return err
		}
	}
	for _, record := range diff.DeleteRecords {
		if err := p.DeleteDomainRecord(domain, record.ID); err != nil {
			return err
		}
	}

	return nil
}
