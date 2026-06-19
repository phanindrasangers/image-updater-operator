# Installation

The operator ships as a Helm chart and a container image, both published to GitHub Container Registry (GHCR).

## Requirements

- Kubernetes 1.25 or newer
- Helm 3.8 or newer (for OCI registry support)
- Cluster permission to install CRDs and a ClusterRole

## Install with Helm

The chart is published as an OCI artifact. Install it directly from GHCR:

```sh
helm install image-updater \
  oci://ghcr.io/phanindrasangers/charts/image-updater-operator \
  --version 0.2.0 \
  --namespace image-updater-system --create-namespace
```

This installs the CRD, the controller `Deployment`, and the RBAC it needs (a ClusterRole to read and patch workloads, plus a namespaced Role for leader-election leases).

### Install from source

Clone the repository and install the local chart:

```sh
git clone https://github.com/phanindrasangers/image-updater-operator.git
cd image-updater-operator
helm install image-updater helm-charts/image-updater-operator \
  --namespace image-updater-system --create-namespace
```

## Common values

Override these with `--set key=value` or a `-f values.yaml` file.

| Value | Default | Description |
|-------|---------|-------------|
| `image.repository` | `ghcr.io/phanindrasangers/image-updater-operator` | Operator image. |
| `image.tag` | chart `appVersion` | Image tag. |
| `replicaCount` | `1` | Replicas. Leader election keeps one active. |
| `leaderElection.enabled` | `true` | Required when running more than one replica. |
| `dashboard.enabled` | `true` | Serve the read-only dashboard inside the pod. |
| `dashboard.service.enabled` | `true` | Create a `ClusterIP` Service for the dashboard. |
| `webhookReceiver.enabled` | `true` | Serve the registry webhook receiver. |
| `webhookReceiver.service.enabled` | `false` | Expose the receiver with a Service (needed for external registries). |
| `resources.limits.memory` | `1Gi` | Memory limit. Git write-back clones repos in memory, so keep headroom for large repositories. |

!!! tip "Memory and Git write-back"
    Git write-back clones the target repository before editing. If you point it at a large repository, raise `resources.limits.memory`. Prefer a small, dedicated config repository for write-back to keep clones fast and cheap.

## Verify

```sh
kubectl get pods -n image-updater-system
kubectl get crd imagepolicies.images.saphire.com
```

The controller pod should be `Running` and the CRD present. Continue to the [Quickstart](quickstart.md).

## Uninstall

```sh
helm uninstall image-updater -n image-updater-system
kubectl delete crd imagepolicies.images.saphire.com   # removes all ImagePolicy objects
```
