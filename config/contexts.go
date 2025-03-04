// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/tanzu-plugin-runtime/config/nodeutils"

	configtypes "github.com/vmware-tanzu/tanzu-plugin-runtime/config/types"
)

// GetContext retrieves the context by name
func GetContext(name string) (*configtypes.Context, error) {
	// Retrieve client config node
	node, err := getClientConfigNode()
	if err != nil {
		return nil, err
	}
	return getContext(node, name)
}

// AddContext add or update context and currentContext
func AddContext(c *configtypes.Context, setCurrent bool) error {
	return SetContext(c, setCurrent)
}

// SetContext add or update context and currentContext
//
//nolint:gocyclo
func SetContext(c *configtypes.Context, setCurrent bool) error {
	// Retrieve client config node
	AcquireTanzuConfigLock()
	defer ReleaseTanzuConfigLock()
	node, err := getClientConfigNodeNoLock()
	if err != nil {
		return err
	}
	// Add or update the context
	persist, err := setContext(node, c)
	if err != nil {
		return err
	}
	if persist {
		err = persistConfig(node)
		if err != nil {
			return err
		}
	}
	// Set current context
	if setCurrent {
		persist, err = setCurrentContext(node, c)
		if err != nil {
			return err
		}
		if persist {
			err = persistConfig(node)
			if err != nil {
				return err
			}
		}
	}

	// Back-fill servers based on contexts
	s := convertContextToServer(c)

	// Add or update server
	persist, err = setServer(node, s)
	if err != nil {
		return err
	}
	if persist {
		err = persistConfig(node)
		if err != nil {
			return err
		}
	}

	// Set current server
	if setCurrent && s.Type == configtypes.ManagementClusterServerType { // nolint:staticcheck
		persist, err = setCurrentServer(node, s.Name)
		if err != nil {
			return err
		}
		if persist {
			err = persistConfig(node)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// DeleteContext delete a context by name
func DeleteContext(name string) error {
	return RemoveContext(name)
}

// RemoveContext delete a context by name
func RemoveContext(name string) error {
	// Retrieve client config node
	AcquireTanzuConfigLock()
	defer ReleaseTanzuConfigLock()
	node, err := getClientConfigNodeNoLock()
	if err != nil {
		return err
	}
	ctx, err := getContext(node, name)
	if err != nil {
		return err
	}
	err = removeCurrentContext(node, ctx)
	if err != nil {
		return err
	}
	err = removeContext(node, name)
	if err != nil {
		return err
	}
	err = removeServer(node, name)
	if err != nil {
		return err
	}
	err = removeCurrentServer(node, name)
	if err != nil {
		return err
	}
	return persistConfig(node)
}

// ContextExists checks if context by name already exists
func ContextExists(name string) (bool, error) {
	exists, _ := GetContext(name)
	return exists != nil, nil
}

// GetCurrentContext retrieves the current context for the specified target
func GetCurrentContext(target configtypes.Target) (c *configtypes.Context, err error) {
	// Retrieve client config node
	node, err := getClientConfigNode()
	if err != nil {
		return nil, err
	}
	return getCurrentContext(node, target)
}

// GetAllCurrentContextsMap returns all current context per Target
func GetAllCurrentContextsMap() (map[configtypes.Target]*configtypes.Context, error) {
	node, err := getClientConfigNodeNoLock()
	if err != nil {
		return nil, err
	}
	return getAllCurrentContextsMap(node)
}

// GetAllCurrentContextsList returns all current context names as list
func GetAllCurrentContextsList() ([]string, error) {
	currentContextsMap, err := GetAllCurrentContextsMap()
	if err != nil {
		return nil, err
	}
	var serverNames []string
	for _, context := range currentContextsMap {
		serverNames = append(serverNames, context.Name)
	}
	return serverNames, nil
}

// SetCurrentContext sets the current context to the specified name if context is present
func SetCurrentContext(name string) error {
	// Retrieve client config node
	AcquireTanzuConfigLock()
	defer ReleaseTanzuConfigLock()
	node, err := getClientConfigNodeNoLock()
	if err != nil {
		return err
	}

	ctx, err := getContext(node, name)
	if err != nil {
		return err
	}
	persist, err := setCurrentContext(node, ctx)
	if err != nil {
		return err
	}
	if persist {
		err = persistConfig(node)
		if err != nil {
			return err
		}
	}
	if ctx.Target == configtypes.TargetK8s {
		persist, err = setCurrentServer(node, name)
		if err != nil {
			return err
		}
		if persist {
			err = persistConfig(node)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// RemoveCurrentContext removed the current context of specified context type
func RemoveCurrentContext(target configtypes.Target) error {
	// Retrieve client config node
	AcquireTanzuConfigLock()
	defer ReleaseTanzuConfigLock()
	node, err := getClientConfigNodeNoLock()
	if err != nil {
		return err
	}
	c, err := getCurrentContext(node, target)
	if err != nil {
		return err
	}
	err = removeCurrentContext(node, &configtypes.Context{Target: target})
	if err != nil {
		return err
	}
	err = removeCurrentServer(node, c.Name)
	if err != nil {
		return err
	}
	return persistConfig(node)
}

// EndpointFromContext retrieved the endpoint from the specified context
func EndpointFromContext(s *configtypes.Context) (endpoint string, err error) {
	switch s.Target {
	case configtypes.TargetK8s:
		return s.ClusterOpts.Endpoint, nil
	case configtypes.TargetTMC:
		return s.GlobalOpts.Endpoint, nil
	default:
		return endpoint, fmt.Errorf("unknown server type %q", s.Target)
	}
}

func getContext(node *yaml.Node, name string) (*configtypes.Context, error) {
	cfg, err := convertNodeToClientConfig(node)
	if err != nil {
		return nil, err
	}
	for _, ctx := range cfg.KnownContexts {
		if ctx.Name == name {
			return ctx, nil
		}
	}
	return nil, fmt.Errorf("context %v not found", name)
}

func getCurrentContext(node *yaml.Node, target configtypes.Target) (*configtypes.Context, error) {
	cfg, err := convertNodeToClientConfig(node)
	if err != nil {
		return nil, err
	}
	return cfg.GetCurrentContext(target)
}

func getAllCurrentContextsMap(node *yaml.Node) (map[configtypes.Target]*configtypes.Context, error) {
	cfg, err := convertNodeToClientConfig(node)
	if err != nil {
		return nil, err
	}
	return cfg.GetAllCurrentContextsMap()
}

func setContexts(node *yaml.Node, contexts []*configtypes.Context) (err error) {
	for _, c := range contexts {
		_, err = setContext(node, c)
		if err != nil {
			return err
		}
	}
	return err
}

func setContext(node *yaml.Node, ctx *configtypes.Context) (persist bool, err error) {
	// Get Patch Strategies from config metadata
	patchStrategies, err := GetConfigMetadataPatchStrategy()
	if err != nil {
		patchStrategies = make(map[string]string)
	}

	var persistDiscoverySources bool

	// Convert context to node
	newContextNode, err := convertContextToNode(ctx)
	if err != nil {
		return persist, err
	}

	// Find the contexts node from the root node
	keys := []nodeutils.Key{
		{Name: KeyContexts, Type: yaml.SequenceNode},
	}
	contextsNode := nodeutils.FindNode(node.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys(keys))
	if contextsNode == nil {
		return persist, err
	}

	exists := false
	var result []*yaml.Node
	// Skip duplicate for context and server similar logic
	//nolint:dupl
	for _, contextNode := range contextsNode.Content {
		if index := nodeutils.GetNodeIndex(contextNode.Content, "name"); index != -1 &&
			contextNode.Content[index].Value == ctx.Name {
			exists = true
			// replace the nodes as per patch strategy
			_, err = nodeutils.DeleteNodes(newContextNode.Content[0], contextNode, nodeutils.WithPatchStrategyKey(KeyContexts), nodeutils.WithPatchStrategies(patchStrategies))
			if err != nil {
				return false, err
			}
			persist, err = nodeutils.MergeNodes(newContextNode.Content[0], contextNode)
			if err != nil {
				return false, err
			}
			persistDiscoverySources, err = setDiscoverySources(contextNode, ctx.DiscoverySources, nodeutils.WithPatchStrategyKey(fmt.Sprintf("%v.%v", KeyContexts, KeyDiscoverySources)), nodeutils.WithPatchStrategies(patchStrategies))
			if err != nil {
				return false, err
			}
			// merge the discovery sources to context
			if persistDiscoverySources {
				_, err = nodeutils.MergeNodes(newContextNode.Content[0], contextNode)
				if err != nil {
					return false, err
				}
			}
			result = append(result, contextNode)
			continue
		}
		result = append(result, contextNode)
	}
	if !exists {
		result = append(result, newContextNode.Content[0])
		persist = true
	}
	contextsNode.Content = result
	return persistDiscoverySources || persist, err
}

func setCurrentContext(node *yaml.Node, ctx *configtypes.Context) (persist bool, err error) {
	// Find current context node in the yaml node
	keys := []nodeutils.Key{
		{Name: KeyCurrentContext, Type: yaml.MappingNode},
	}
	currentContextNode := nodeutils.FindNode(node.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys(keys))
	if currentContextNode == nil {
		return persist, nodeutils.ErrNodeNotFound
	}
	if index := nodeutils.GetNodeIndex(currentContextNode.Content, string(ctx.Target)); index != -1 {
		if currentContextNode.Content[index].Value != ctx.Name {
			currentContextNode.Content[index].Value = ctx.Name
			currentContextNode.Content[index].Style = 0
			persist = true
		}
	} else {
		currentContextNode.Content = append(currentContextNode.Content, nodeutils.CreateScalarNode(string(ctx.Target), ctx.Name)...)
		persist = true
	}
	return persist, err
}

func removeCurrentContext(node *yaml.Node, ctx *configtypes.Context) error {
	// Find current context node in the yaml node
	keys := []nodeutils.Key{
		{Name: KeyCurrentContext},
	}

	currentContextNode := nodeutils.FindNode(node.Content[0], nodeutils.WithKeys(keys))
	if currentContextNode == nil {
		return nil
	}
	targetNodeIndex := nodeutils.GetNodeIndex(currentContextNode.Content, string(ctx.Target))
	if targetNodeIndex == -1 {
		return nil
	}
	if currentContextNode.Content[targetNodeIndex].Value == ctx.Name || ctx.Name == "" {
		targetNodeIndex--
		currentContextNode.Content = append(currentContextNode.Content[:targetNodeIndex], currentContextNode.Content[targetNodeIndex+2:]...)
	}
	return nil
}

//nolint:dupl
func removeContext(node *yaml.Node, name string) error {
	// Find the contexts node in the yaml node
	keys := []nodeutils.Key{
		{Name: KeyContexts},
	}
	contextsNode := nodeutils.FindNode(node.Content[0], nodeutils.WithKeys(keys))
	if contextsNode == nil {
		return nil
	}
	var contexts []*yaml.Node
	for _, contextNode := range contextsNode.Content {
		if index := nodeutils.GetNodeIndex(contextNode.Content, "name"); index != -1 && contextNode.Content[index].Value == name {
			continue
		}
		contexts = append(contexts, contextNode)
	}
	contextsNode.Content = contexts
	return nil
}
