// Copyright 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package context contains all the cross version api compatibility tests for context apis
package context

import (
	"fmt"

	"github.com/vmware-tanzu/tanzu-plugin-runtime/test/compatibility/core"
	"github.com/vmware-tanzu/tanzu-plugin-runtime/test/compatibility/framework"
	"github.com/vmware-tanzu/tanzu-plugin-runtime/test/compatibility/framework/compatibilitytests/common"
)

// DefaultSetContextInputOptions helper method to construct SetContext API input options
func DefaultSetContextInputOptions(version core.RuntimeVersion, contextName string) *framework.SetContextInputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.SetContextInputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			ContextOpts: &framework.ContextOpts{
				Name:   contextName,
				Target: framework.TargetK8s,
				GlobalOpts: &framework.GlobalServerOpts{
					Endpoint: "default-compatibility-test-endpoint",
				},
			},
		}
	case core.Version0254:
		return &framework.SetContextInputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.Version0254,
			},
			ContextOpts: &framework.ContextOpts{
				Name: contextName,
				Type: framework.CtxTypeK8s,
				GlobalOpts: &framework.GlobalServerOpts{
					Endpoint: "default-compatibility-test-endpoint",
				},
			},
		}
	}
	return nil
}

// DefaultGetContextInputOptions helper method to construct GetContext API input options
func DefaultGetContextInputOptions(version core.RuntimeVersion, contextName string) *framework.GetContextInputOptions {
	return &framework.GetContextInputOptions{
		RuntimeAPIVersion: &core.RuntimeAPIVersion{
			RuntimeVersion: version,
		},
		ContextName: contextName,
	}
}

// DefaultGetContextOutputOptions helper method to construct GetContext API output options
func DefaultGetContextOutputOptions(version core.RuntimeVersion, contextName string) *framework.GetContextOutputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.GetContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			ContextOpts: &framework.ContextOpts{
				Name:   contextName,
				Target: framework.TargetK8s,
				GlobalOpts: &framework.GlobalServerOpts{
					Endpoint: common.DefaultEndpoint,
				},
			},
			ValidationStrategy: core.ValidationStrategyStrict,
		}
	case core.Version0254:
		return &framework.GetContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.Version0254,
			},
			ContextOpts: &framework.ContextOpts{
				Name: contextName,
				Type: framework.CtxTypeK8s,
				GlobalOpts: &framework.GlobalServerOpts{
					Endpoint: common.DefaultEndpoint,
				},
			},
		}
	}
	return nil
}

// DefaultGetContextOutputOptionsWithError helper method to construct GetContext API output options with error
func DefaultGetContextOutputOptionsWithError(version core.RuntimeVersion, contextName string) *framework.GetContextOutputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.GetContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Error: fmt.Sprintf("context %v not found", contextName),
		}
	case core.Version0254:
		return &framework.GetContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.Version0254,
			},
			Error: fmt.Sprintf("could not find context \"%v\"", contextName),
		}
	}
	return nil
}

// DefaultSetCurrentContextInputOptions helper method to construct SetCurrentContext API input options
func DefaultSetCurrentContextInputOptions(version core.RuntimeVersion, contextName string) *framework.SetCurrentContextInputOptions {
	return &framework.SetCurrentContextInputOptions{
		RuntimeAPIVersion: &core.RuntimeAPIVersion{
			RuntimeVersion: version,
		},
		ContextName: contextName,
	}
}

// DefaultGetCurrentContextInputOptions helper method to construct GetCurrentContext API input options
func DefaultGetCurrentContextInputOptions(version core.RuntimeVersion) *framework.GetCurrentContextInputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.GetCurrentContextInputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Target: framework.TargetK8s,
		}
	case core.Version0254:
		return &framework.GetCurrentContextInputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.Version0254,
			},
			ContextType: framework.CtxTypeK8s,
		}
	}
	return nil
}

// DefaultGetCurrentContextOutputOptions helper method to construct GetCurrentContext API output options
func DefaultGetCurrentContextOutputOptions(version core.RuntimeVersion, contextName string) *framework.GetCurrentContextOutputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.GetCurrentContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.VersionLatest,
			},
			ContextOpts: &framework.ContextOpts{
				Name:   contextName,
				Target: framework.TargetK8s,
				GlobalOpts: &framework.GlobalServerOpts{
					Endpoint: common.DefaultEndpoint,
				},
			},
			ValidationStrategy: core.ValidationStrategyStrict,
		}
	case core.Version0254:
		return &framework.GetCurrentContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.Version0254,
			},
			ContextOpts: &framework.ContextOpts{
				Name: contextName,
				Type: framework.CtxTypeK8s,
				GlobalOpts: &framework.GlobalServerOpts{
					Endpoint: common.DefaultEndpoint,
				},
			},
		}
	}
	return nil
}

// DefaultGetCurrentContextOutputOptionsWithError helper method to construct GetCurrentContext API output options with error
func DefaultGetCurrentContextOutputOptionsWithError(version core.RuntimeVersion) *framework.GetCurrentContextOutputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.GetCurrentContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Error: fmt.Sprintf("no current context set for target \"%v\"", framework.TargetK8s),
		}
	case core.Version0254:
		return &framework.GetCurrentContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Error: fmt.Sprintf("no current context set for type \"%v\"", framework.CtxTypeK8s),
		}
	}
	return nil
}

// DefaultRemoveCurrentContextInputOptions helper method to construct RemoveCurrentContext API input options
func DefaultRemoveCurrentContextInputOptions(version core.RuntimeVersion) *framework.RemoveCurrentContextInputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.RemoveCurrentContextInputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Target: framework.TargetK8s,
		}
	}
	return nil
}

// DefaultRemoveCurrentContextOutputOptionsWithError helper method to construct RemoveCurrentContext API output option
func DefaultRemoveCurrentContextOutputOptionsWithError(version core.RuntimeVersion) *framework.RemoveCurrentContextOutputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.RemoveCurrentContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Error: fmt.Sprintf("no current context set for target \"%v\"", framework.TargetK8s),
		}
	case core.Version0254:
		return &framework.RemoveCurrentContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Error: fmt.Sprintf("no current context set for type \"%v\"", framework.CtxTypeK8s),
		}
	}
	return nil
}

// DefaultDeleteContextInputOptions helper method to construct DeleteContext API input options
func DefaultDeleteContextInputOptions(version core.RuntimeVersion, contextName string) *framework.DeleteContextInputOptions {
	return &framework.DeleteContextInputOptions{
		RuntimeAPIVersion: &core.RuntimeAPIVersion{
			RuntimeVersion: version,
		},
		ContextName: contextName,
	}
}

// DefaultDeleteContextOutputOptionsWithError helper method to construct DeleteContext API output options
func DefaultDeleteContextOutputOptionsWithError(version core.RuntimeVersion, contextName string) *framework.DeleteContextOutputOptions {
	switch version {
	case core.VersionLatest, core.Version0280:
		return &framework.DeleteContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: version,
			},
			Error: fmt.Sprintf("context %v not found", contextName),
		}
	case core.Version0254:
		return &framework.DeleteContextOutputOptions{
			RuntimeAPIVersion: &core.RuntimeAPIVersion{
				RuntimeVersion: core.Version0254,
			},
			Error: fmt.Sprintf("could not find context \"%v\"", contextName),
		}
	}
	return nil
}
