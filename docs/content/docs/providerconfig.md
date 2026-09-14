---
title: ProviderConfig
weight: 2
---

## Credentials

Credentials are a JSON document in a `Secret` key. The fields are those of the
coderd Terraform provider's configuration schema:

| Field | Required | Description |
| --- | --- | --- |
| `url` | yes | The URL of your Coder deployment. |
| `token` | yes | A Coder API token. Most resource types need elevated permissions. |
| `default_organization_id` | no | Organization used by resources that do not name one. |
| `headers` | no | Extra HTTP headers to send with every request. |

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

## Cluster-scoped and namespaced

Every kind this provider serves exists twice, and the `ProviderConfig` a managed
resource refers to must come from the same family.

### Cluster-scoped — `coderd.crossplane.io`

Used by managed resources under `*.coderd.crossplane.io`.

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

### Namespaced — `coderd.m.crossplane.io`

Used by managed resources under `*.coderd.m.crossplane.io`. A namespaced
`ProviderConfig` is resolved within the managed resource's own namespace:

```yaml
apiVersion: coderd.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: crossplane-system
spec:
  credentials:
    source: Secret
    secretRef:
      name: coderd-creds
      namespace: crossplane-system
      key: credentials
```

A `ClusterProviderConfig` lets namespaced resources in any namespace share one
configuration:

```yaml
apiVersion: coderd.m.crossplane.io/v1beta1
kind: ClusterProviderConfig
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

## Selecting one

A managed resource picks its configuration with `spec.providerConfigRef`. It
defaults to the one named `default`, so the examples above need no reference at
all.

```yaml
spec:
  providerConfigRef:
    name: production
```
