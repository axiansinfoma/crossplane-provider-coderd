// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	aiprovider "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/agents/aiprovider"
	defaultmodel "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/agents/defaultmodel"
	mcpserver "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/agents/mcpserver"
	model "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/agents/model"
	systemprompt "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/agents/systemprompt"
	license "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/deployment/license"
	oauth2providersettings "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/deployment/oauth2providersettings"
	provisionerkey "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/deployment/provisionerkey"
	workspaceproxy "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/deployment/workspaceproxy"
	group "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/iam/group"
	organization "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/iam/organization"
	organizationgroupsync "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/iam/organizationgroupsync"
	organizationsyncsettings "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/iam/organizationsyncsettings"
	user "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/iam/user"
	providerconfig "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/providerconfig"
	template "github.com/axiansinfoma/crossplane-provider-coderd/internal/controller/cluster/template/template"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		aiprovider.Setup,
		defaultmodel.Setup,
		mcpserver.Setup,
		model.Setup,
		systemprompt.Setup,
		license.Setup,
		oauth2providersettings.Setup,
		provisionerkey.Setup,
		workspaceproxy.Setup,
		group.Setup,
		organization.Setup,
		organizationgroupsync.Setup,
		organizationsyncsettings.Setup,
		user.Setup,
		providerconfig.Setup,
		template.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		aiprovider.SetupGated,
		defaultmodel.SetupGated,
		mcpserver.SetupGated,
		model.SetupGated,
		systemprompt.SetupGated,
		license.SetupGated,
		oauth2providersettings.SetupGated,
		provisionerkey.SetupGated,
		workspaceproxy.SetupGated,
		group.SetupGated,
		organization.SetupGated,
		organizationgroupsync.SetupGated,
		organizationsyncsettings.SetupGated,
		user.SetupGated,
		providerconfig.SetupGated,
		template.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		aiprovider.SetupWebhookWithManager,
		defaultmodel.SetupWebhookWithManager,
		mcpserver.SetupWebhookWithManager,
		model.SetupWebhookWithManager,
		systemprompt.SetupWebhookWithManager,
		license.SetupWebhookWithManager,
		oauth2providersettings.SetupWebhookWithManager,
		provisionerkey.SetupWebhookWithManager,
		workspaceproxy.SetupWebhookWithManager,
		group.SetupWebhookWithManager,
		organization.SetupWebhookWithManager,
		organizationgroupsync.SetupWebhookWithManager,
		organizationsyncsettings.SetupWebhookWithManager,
		user.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		template.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
