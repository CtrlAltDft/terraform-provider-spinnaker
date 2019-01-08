package provider

import (
	"strings"

	"github.com/jgramoll/terraform-provider-spinnaker/client"
)

type deployStageClusterCapacity struct {
	Desired int `mapstructure:"desired"`
	Max     int `mapstructure:"max"`
	Min     int `mapstructure:"min"`
}

type deployStageCluster struct {
	Account                             string                       `mapstructure:"account"`
	Application                         string                       `mapstructure:"application"`
	AvailabilityZones                   []map[string][]string        `mapstructure:"availability_zones"`
	Capacity                            []deployStageClusterCapacity `mapstructure:"capacity"`
	CloudProvider                       string                       `mapstructure:"cloud_provider"`
	Cooldown                            int                          `mapstructure:"cooldown"`
	CopySourceCustomBlockDeviceMappings bool                         `mapstructure:"copy_source_custom_block_device_mappings"`
	EBSOptimized                        bool                         `mapstructure:"ebs_optimized"`
	EnabledMetrics                      []string                     `mapstructure:"enabled_metrics"`
	FreeFormDetails                     string                       `mapstructure:"free_form_details"`
	HealthCheckGracePeriod              string                       `mapstructure:"health_check_grace_period"`
	HealthCheckType                     string                       `mapstructure:"health_check_type"`
	IAMRole                             string                       `mapstructure:"iam_role"`
	InstanceMonitoring                  bool                         `mapstructure:"instance_monitoring"`
	InstanceType                        string                       `mapstructure:"instance_type"`
	KeyPair                             string                       `mapstructure:"key_pair"`
	LoadBalancers                       []string                     `mapstructure:"load_balancers"`
	Moniker                             []moniker                    `mapstructure:"moniker"`
	Provider                            string                       `mapstructure:"provider"`
	SecurityGroups                      []string                     `mapstructure:"security_groups"`
	SpelLoadBalancers                   []string                     `mapstructure:"spel_load_balancers"`
	SpelTargetGroups                    []string                     `mapstructure:"spel_target_groups"`
	SpotPrice                           string                       `mapstructure:"spot_price"`
	Stack                               string                       `mapstructure:"stack"`
	Strategy                            string                       `mapstructure:"strategy"`
	SubnetType                          string                       `mapstructure:"subnet_type"`
	SuspendedProcesses                  []string                     `mapstructure:"suspended_processes"`
	Tags                                map[string]string            `mapstructure:"tags"`
	TargetGroups                        []string                     `mapstructure:"target_groups"`
	TargetHealthyDeployPercentage       int                          `mapstructure:"target_healthy_deploy_percentage"`
	TerminationPolicies                 []string                     `mapstructure:"termination_policies"`
	UseAmiBlockDeviceMappings           bool                         `mapstructure:"use_ami_block_device_mappings"`
	UseSourceCapacity                   bool                         `mapstructure:"use_source_capacity"`
}

func (c *deployStageCluster) clientCapacity() client.DeployStageClusterCapacity {
	if len(c.Capacity) > 0 {
		return client.DeployStageClusterCapacity(c.Capacity[0])
	}
	return client.DeployStageClusterCapacity{}
}

func (c *deployStageCluster) clientMoniker() client.Moniker {
	if len(c.Moniker) > 0 {
		return client.Moniker(c.Moniker[0])
	}
	return client.Moniker{}
}

func (c *deployStageCluster) clientAvailabilityZones() *map[string][]string {
	newAZ := map[string][]string{}
	for _, regions := range c.AvailabilityZones {
		for region, zones := range regions {
			if len(zones) == 0 {
				continue
			}
			// TODO unit test
			newAZ[strings.Replace(region, "_", "-", -1)] = zones
		}
	}
	return &newAZ
}

func (c *deployStageCluster) importAvailabilityZones(clientCluster *client.DeployStageCluster) {
	for region, zones := range clientCluster.AvailabilityZones {
		newZone := map[string][]string{
			strings.Replace(region, "-", "_", -1): zones,
		}
		// TODO unit test
		c.AvailabilityZones = append(c.AvailabilityZones, newZone)
	}
}

func (c *deployStageCluster) toClientCluster() *client.DeployStageCluster {
	// TODO better way?
	return &client.DeployStageCluster{
		Account:           c.Account,
		Application:       c.Application,
		AvailabilityZones: *c.clientAvailabilityZones(),
		Capacity:          c.clientCapacity(),
		CloudProvider:     c.CloudProvider,
		Cooldown:          c.Cooldown,

		CopySourceCustomBlockDeviceMappings: c.CopySourceCustomBlockDeviceMappings,

		EBSOptimized:           c.EBSOptimized,
		EnabledMetrics:         c.EnabledMetrics,
		FreeFormDetails:        c.FreeFormDetails,
		HealthCheckGracePeriod: c.HealthCheckGracePeriod,
		HealthCheckType:        c.HealthCheckType,
		IAMRole:                c.IAMRole,
		InstanceMonitoring:     c.InstanceMonitoring,
		InstanceType:           c.InstanceType,
		KeyPair:                c.KeyPair,
		LoadBalancers:          c.LoadBalancers,
		Moniker:                c.clientMoniker(),
		Provider:               c.Provider,
		SecurityGroups:         c.SecurityGroups,
		SpelLoadBalancers:      c.SpelLoadBalancers,
		SpelTargetGroups:       c.SpelTargetGroups,
		SpotPrice:              c.SpotPrice,
		Stack:                  c.Stack,
		Strategy:               c.Strategy,
		SubnetType:             c.SubnetType,
		SuspendedProcesses:     c.SuspendedProcesses,
		Tags:                   c.Tags,
		TargetGroups:           c.TargetGroups,

		TargetHealthyDeployPercentage: c.TargetHealthyDeployPercentage,
		TerminationPolicies:           c.TerminationPolicies,
		UseAmiBlockDeviceMappings:     c.UseAmiBlockDeviceMappings,
		UseSourceCapacity:             c.UseSourceCapacity,
	}
}

func newClusterFromClientCluster(c *client.DeployStageCluster) *deployStageCluster {
	newCluster := deployStageCluster{
		Account:       c.Account,
		Application:   c.Application,
		CloudProvider: c.CloudProvider,
		Cooldown:      c.Cooldown,

		CopySourceCustomBlockDeviceMappings: c.CopySourceCustomBlockDeviceMappings,

		EBSOptimized:           c.EBSOptimized,
		EnabledMetrics:         c.EnabledMetrics,
		FreeFormDetails:        c.FreeFormDetails,
		HealthCheckGracePeriod: c.HealthCheckGracePeriod,
		HealthCheckType:        c.HealthCheckType,
		IAMRole:                c.IAMRole,
		InstanceMonitoring:     c.InstanceMonitoring,
		InstanceType:           c.InstanceType,
		KeyPair:                c.KeyPair,
		LoadBalancers:          c.LoadBalancers,
		Provider:               c.Provider,
		SecurityGroups:         c.SecurityGroups,
		SpelLoadBalancers:      c.SpelLoadBalancers,
		SpelTargetGroups:       c.SpelTargetGroups,
		SpotPrice:              c.SpotPrice,
		Stack:                  c.Stack,
		Strategy:               c.Strategy,
		SubnetType:             c.SubnetType,
		SuspendedProcesses:     c.SuspendedProcesses,
		Tags:                   c.Tags,
		TargetGroups:           c.TargetGroups,

		TargetHealthyDeployPercentage: c.TargetHealthyDeployPercentage,
		TerminationPolicies:           c.TerminationPolicies,
		UseAmiBlockDeviceMappings:     c.UseAmiBlockDeviceMappings,
		UseSourceCapacity:             c.UseSourceCapacity,
	}
	newCluster.importAvailabilityZones(c)
	newCluster.Capacity = append(newCluster.Capacity, deployStageClusterCapacity(c.Capacity))
	newCluster.Moniker = append(newCluster.Moniker, moniker(c.Moniker))
	return &newCluster
}