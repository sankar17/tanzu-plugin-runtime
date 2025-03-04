// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Target is the namespace of the CLI to which plugin is applicable
type Target string

const (
	// TargetK8s is a kubernetes target of the CLI
	// This target applies if the plugin is interacting with a Kubernetes cluster
	TargetK8s Target = "kubernetes"
	targetK8s Target = "k8s"

	// TargetTMC is a Tanzu Mission Control target of the CLI
	// This target applies if the plugin is interacting with a Tanzu Mission Control endpoint
	TargetTMC Target = "mission-control"
	targetTMC Target = "tmc"

	// TargetGlobal is used for plugins that are not associated with any target
	TargetGlobal Target = "global"

	// TargetUnknown specifies that the target is not currently known
	TargetUnknown Target = ""
)

var (
	// SupportedTargets is a list of all supported Target
	SupportedTargets = []Target{TargetK8s, TargetTMC}
)

const (
	// AllUnstableVersions allows all plugin versions
	AllUnstableVersions VersionSelectorLevel = "all"
	// AlphaUnstableVersions allows all alpha tagged versions
	AlphaUnstableVersions VersionSelectorLevel = "alpha"
	// ExperimentalUnstableVersions includes all pre-releases, minus +build tags
	ExperimentalUnstableVersions VersionSelectorLevel = "experimental"
	// NoUnstableVersions allows no unstable plugin versions, format major.minor.patch only
	NoUnstableVersions VersionSelectorLevel = "none"
)

const (
	// FeatureCli allows a feature to be set at the CLI level (globally) rather than for a single plugin
	FeatureCli string = "cli"
	// EditionStandard refers to the standard edition
	// Edition value (in config) affects branding and cluster creation
	EditionStandard = "tkg"
	// EditionCommunity refers to the community edition
	EditionCommunity = "tce"
)

// EditionSelector allows selecting edition versions based on config file
type EditionSelector string

// VersionSelectorLevel allows selecting plugin versions based on semver properties
type VersionSelectorLevel string

// IsGlobal tells if the server is global.
//
// Deprecated: This API is deprecated. Use Context.Target instead.
func (s *Server) IsGlobal() bool {
	return s.Type == GlobalServerType
}

// IsManagementCluster tells if the server is a management cluster.
//
// Deprecated: This API is deprecated. Use context.IsManagementCluster instead.
func (s *Server) IsManagementCluster() bool {
	return s.Type == ManagementClusterServerType
}

// GetCurrentServer returns the current server.
//
// Deprecated: This API is deprecated. Use GetCurrentContext() instead.
func (c *ClientConfig) GetCurrentServer() (*Server, error) {
	for _, server := range c.KnownServers {
		if server.Name == c.CurrentServer {
			return server, nil
		}
	}
	return nil, fmt.Errorf("current server %q not found", c.CurrentServer)
}

// HasServer tells whether the Server by the given name exists.
func (c *ClientConfig) HasServer(name string) bool {
	for _, s := range c.KnownServers {
		if s.Name == name {
			return true
		}
	}
	return false
}

// GetContext by name.
func (c *ClientConfig) GetContext(name string) (*Context, error) {
	for _, ctx := range c.KnownContexts {
		if ctx.Name == name {
			return ctx, nil
		}
	}
	return nil, fmt.Errorf("could not find context %q", name)
}

// HasContext tells whether the Context by the given name exists.
func (c *ClientConfig) HasContext(name string) bool {
	_, err := c.GetContext(name)
	return err == nil
}

// GetCurrentContext returns the current context for the given type.
func (c *ClientConfig) GetCurrentContext(target Target) (*Context, error) {
	ctxName := c.CurrentContext[target]
	if ctxName == "" {
		return nil, fmt.Errorf("no current context set for target %q", target)
	}
	ctx, err := c.GetContext(ctxName)
	if err != nil {
		return nil, fmt.Errorf("unable to get current context: %s", err.Error())
	}
	return ctx, nil
}

// GetAllCurrentContextsMap returns all current context per Target
func (c *ClientConfig) GetAllCurrentContextsMap() (map[Target]*Context, error) {
	currentContexts := make(map[Target]*Context)
	for _, target := range SupportedTargets {
		context, err := c.GetCurrentContext(target)
		if err == nil && context != nil {
			currentContexts[target] = context
		}
	}
	return currentContexts, nil
}

// GetAllCurrentContextsList returns all current context names as list
func (c *ClientConfig) GetAllCurrentContextsList() ([]string, error) {
	var serverNames []string
	currentContextsMap, err := c.GetAllCurrentContextsMap()
	if err != nil {
		return nil, err
	}

	for _, context := range currentContextsMap {
		serverNames = append(serverNames, context.Name)
	}
	return serverNames, nil
}

// SetCurrentContext sets the current context for the given target.
func (c *ClientConfig) SetCurrentContext(target Target, ctxName string) error {
	if c.CurrentContext == nil {
		c.CurrentContext = make(map[Target]string)
	}
	c.CurrentContext[target] = ctxName
	ctx, err := c.GetContext(ctxName)
	if err != nil {
		return err
	}
	if ctx.IsManagementCluster() || ctx.Target == TargetTMC {
		c.CurrentServer = ctxName
	}
	return nil
}

// IsManagementCluster tells if the context is for a management cluster.
func (c *Context) IsManagementCluster() bool {
	return c != nil && c.Target == TargetK8s && c.ClusterOpts != nil && c.ClusterOpts.IsManagementCluster
}

// SetUnstableVersionSelector will help determine the unstable versions supported
// In order of restrictiveness:
// "all" -> "alpha" -> "experimental" -> "none"
// none: return stable versions only. the default for both the config and the old flag.
// alpha: only versions tagged with -alpha
// experimental: all pre-release versions without +build semver data
// all: return all unstable versions.
//
// Deprecated: This API is deprecated.
func (c *ClientConfig) SetUnstableVersionSelector(f VersionSelectorLevel) {
	if c.ClientOptions == nil {
		c.ClientOptions = &ClientOptions{}
	}
	if c.ClientOptions.CLI == nil {
		c.ClientOptions.CLI = &CLIOptions{}
	}
	switch f {
	case AllUnstableVersions, AlphaUnstableVersions, ExperimentalUnstableVersions, NoUnstableVersions:
		c.ClientOptions.CLI.UnstableVersionSelector = f
		return
	}
	c.ClientOptions.CLI.UnstableVersionSelector = AllUnstableVersions
}

// IsConfigFeatureActivated return true if the feature is activated, false if not. An error if the featurePath is malformed
func (c *ClientConfig) IsConfigFeatureActivated(featurePath string) (bool, error) {
	plugin, flag, err := c.SplitFeaturePath(featurePath)
	if err != nil {
		return false, err
	}

	if c.ClientOptions == nil || c.ClientOptions.Features == nil ||
		c.ClientOptions.Features[plugin] == nil || c.ClientOptions.Features[plugin][flag] == "" {
		return false, nil
	}

	booleanValue, err := strconv.ParseBool(c.ClientOptions.Features[plugin][flag])
	if err != nil {
		errMsg := "error converting " + featurePath + " entry '" + c.ClientOptions.Features[plugin][flag] + "' to boolean value: " + err.Error()
		return false, errors.New(errMsg)
	}
	return booleanValue, nil
}

// GetEnvConfigurations returns a map of environment variables to values
// it returns nil if configuration is not yet defined
func (c *ClientConfig) GetEnvConfigurations() map[string]string {
	if c.ClientOptions == nil || c.ClientOptions.Env == nil {
		return nil
	}
	return c.ClientOptions.Env
}

// SplitFeaturePath splits a feature's path into the pluginName and the featureName
// For example "features.management-cluster.dual-stack" returns "management-cluster", "dual-stack"
// An error results from a malformed path, including any path that does not start with "features."
func (c *ClientConfig) SplitFeaturePath(featurePath string) (string, string, error) {
	// parse the param
	paramArray := strings.Split(featurePath, ".")
	if len(paramArray) != 3 {
		return "", "", errors.New("unable to parse feature name config parameter into three parts [" + featurePath + "]  (was expecting features.<plugin>.<feature>)")
	}

	featuresLiteral := paramArray[0]
	plugin := paramArray[1]
	flag := paramArray[2]

	if featuresLiteral != "features" {
		return "", "", errors.New("unsupported feature config path parameter [" + featuresLiteral + "] (was expecting 'features.<plugin>.<feature>')")
	}
	return plugin, flag, nil
}

// SetEditionSelector indicates the edition of tanzu to be run
// EditionStandard is the default, EditionCommunity is also available.
// These values affect branding and cluster creation
//
// Deprecated: This API is deprecated.
func (c *ClientConfig) SetEditionSelector(edition EditionSelector) {
	if c.ClientOptions == nil {
		c.ClientOptions = &ClientOptions{}
	}
	if c.ClientOptions.CLI == nil {
		c.ClientOptions.CLI = &CLIOptions{}
	}
	switch edition {
	case EditionCommunity, EditionStandard:
		c.ClientOptions.CLI.Edition = edition
		return
	}
	c.ClientOptions.CLI.UnstableVersionSelector = EditionStandard
}
