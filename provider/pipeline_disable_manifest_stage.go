package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jgramoll/terraform-provider-spinnaker/client"
)

type disableManifestStage struct {
	baseStage `mapstructure:",squash"`

	Account       string `mapstructure:"account"`
	App           string `mapstructure:"app"`
	CloudProvider string `mapstructure:"cloud_provider"`
	Cluster       string `mapstructure:"cluster"`
	Criteria      string `mapstructure:"criteria"`
	Kind          string `mapstructure:"kind"`
	// kinds          string `mapstructure:"kinds"`
	// labelSelectors string `mapstructure:"labelSelectors"`
	Location     string `mapstructure:"location"`
	ManifestName string `mapstructure:"manifest_name"`
	Mode         string `mapstructure:"mode"`
}

func newDisableManifestStage() *disableManifestStage {
	return &disableManifestStage{
		baseStage: *newBaseStage(),
	}
}

func (s *disableManifestStage) toClientStage(config *client.Config, refID string) (client.Stage, error) {
	cs := client.NewDisableManifestStage()
	err := s.baseToClientStage(&cs.BaseStage, refID, newDefaultNotificationInterface)
	if err != nil {
		return nil, err
	}

	cs.Account = s.Account
	cs.App = s.App
	cs.CloudProvider = s.CloudProvider
	cs.Cluster = s.Cluster
	cs.Criteria = s.Criteria
	cs.Kind = s.Kind
	cs.Location = s.Location
	cs.ManifestName = s.ManifestName
	cs.Mode = s.Mode

	return cs, nil
}

func (*disableManifestStage) fromClientStage(cs client.Stage) (stage, error) {
	clientStage := cs.(*client.DisableManifestStage)
	newStage := newDisableManifestStage()
	err := newStage.baseFromClientStage(&clientStage.BaseStage, newDefaultNotificationInterface)
	if err != nil {
		return nil, err
	}

	newStage.Account = clientStage.Account
	newStage.App = clientStage.App
	newStage.CloudProvider = clientStage.CloudProvider
	newStage.Cluster = clientStage.Cluster
	newStage.Criteria = clientStage.Criteria
	newStage.Kind = clientStage.Kind
	newStage.Location = clientStage.Location
	newStage.ManifestName = clientStage.ManifestName
	newStage.Mode = clientStage.Mode

	return newStage, nil
}

func (s *disableManifestStage) SetResourceData(d *schema.ResourceData) error {
	err := s.baseSetResourceData(d)
	if err != nil {
		return err
	}

	err = d.Set("account", s.Account)
	if err != nil {
		return err
	}
	err = d.Set("app", s.App)
	if err != nil {
		return err
	}
	err = d.Set("cloud_provider", s.CloudProvider)
	if err != nil {
		return err
	}
	err = d.Set("cluster", s.Cluster)
	if err != nil {
		return err
	}
	err = d.Set("criteria", s.Criteria)
	if err != nil {
		return err
	}
	err = d.Set("kind", s.Kind)
	if err != nil {
		return err
	}
	err = d.Set("location", s.Location)
	if err != nil {
		return err
	}
	err = d.Set("manifest_name", s.ManifestName)
	if err != nil {
		return err
	}
	err = d.Set("mode", s.Mode)
	if err != nil {
		return err
	}

	return nil
}
