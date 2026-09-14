---
title: Getting started
weight: 1
---

## Requirements

- Crossplane **v2** on a Kubernetes cluster.
- A Coder deployment and an API token. Most resource types need a token with
  elevated permissions.

## Install the provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-coderd
spec:
  package: ghcr.io/axiansinfoma/provider-coderd:v0.1.0
```

```console
kubectl apply -f https://raw.githubusercontent.com/axiansinfoma/crossplane-provider-coderd/main/examples/install.yaml
```

Wait for it to become healthy:

```console
kubectl wait provider.pkg provider-coderd --for condition=Healthy --timeout 5m
```

## Provide credentials

The provider reads its credentials from a `Secret` whose keys mirror the
configuration schema of the coderd Terraform provider. `url` and `token` are
required; `default_organization_id` and `headers` are optional.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: coderd-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "url": "https://coder.example.com",
      "token": "REPLACE_ME"
    }
```

Point a `ProviderConfig` at it:

```yaml
apiVersion: coderd.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: coderd-creds
      namespace: crossplane-system
      key: credentials
```

See [ProviderConfig](../providerconfig) for the namespaced flavour.

## Create a resource

```yaml
apiVersion: iam.coderd.crossplane.io/v1alpha1
kind: Organization
metadata:
  name: engineering
spec:
  forProvider:
    name: engineering
    displayName: Engineering
```

```console
kubectl apply -f organization.yaml
kubectl wait organization.iam.coderd.crossplane.io engineering --for condition=Ready --timeout 5m
```

If it does not become ready, the reason is on the resource itself:

```console
kubectl describe organization.iam.coderd.crossplane.io engineering
```

The
[`examples/`](https://github.com/axiansinfoma/crossplane-provider-coderd/tree/main/examples)
directory has a manifest per kind, in both the cluster-scoped and the namespaced
flavour.
