---
title: Provider Coderd
layout: hextra-home
---

{{< hextra/hero-headline >}}
Manage Coder with Crossplane
{{< /hextra/hero-headline >}}

<p class="hx-mt-6 hx-mb-6">
A Crossplane provider for <a href="https://coder.com">Coder</a>, generated with
<a href="https://github.com/crossplane/upjet">Upjet</a> from the
<a href="https://github.com/coder/terraform-provider-coderd">coderd Terraform provider</a>.
Organizations, users, groups, templates, licenses, workspace proxies and AI agent
configuration, all as Kubernetes resources.
</p>

{{< hextra/hero-button text="Get started" link="docs/getting-started" >}}

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="No Terraform at runtime"
    subtitle="The coderd provider is linked into the controller and driven in-process. No Terraform CLI, no plugin in the image, no subprocess per reconcile."
  >}}
  {{< hextra/feature-card
    title="Crossplane v2, both scopes"
    subtitle="Every kind is served twice: cluster-scoped under *.coderd.crossplane.io and namespaced under *.coderd.m.crossplane.io."
  >}}
  {{< hextra/feature-card
    title="Generated, and kept in step"
    subtitle="A weekly job bumps the upstream provider, regenerates everything that derives from it and opens a pull request with a breaking-change report."
  >}}
{{< /hextra/feature-grid >}}
