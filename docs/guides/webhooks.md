# Instant updates with webhooks

By default the operator scans on each policy's `interval`. To react immediately when a registry pushes a new tag, point the registry's webhook at the operator's receiver to trigger an out-of-band rescan.

!!! info "Prerequisites"
    - The operator [installed](../getting-started/installation.md), with at least one `ImagePolicy`.
    - The webhook receiver enabled and exposed (step below) so your registry can reach it from outside the cluster: an Ingress or a `LoadBalancer`/`NodePort` Service.
    - Optionally, a shared token to authenticate incoming requests.

## The receiver

The receiver listens on `--webhook-receiver-bind-address` (default `:9090`). Enable it and expose it with the chart:

```sh
helm upgrade image-updater oci://ghcr.io/phanindrasangers/charts/image-updater-operator \
  --version 0.2.0 -n image-updater-system --reuse-values \
  --set webhookReceiver.enabled=true \
  --set webhookReceiver.service.enabled=true
```

Then route external traffic to the Service via an Ingress or `LoadBalancer`.

## Endpoints

| Path | Source |
|------|--------|
| `POST /webhook/dockerhub` | Docker Hub |
| `POST /webhook/harbor` | Harbor |
| `POST /webhook/generic` | Anything; body `{"repository":"<host>/<repo>"}` or `{"image":"<host>/<repo>:<tag>"}` |

A received event triggers an immediate rescan of any `ImagePolicy` whose `imageRepository` matches, without waiting for the interval.

## Authentication

Set the `WEBHOOK_RECEIVER_TOKEN` environment variable to require a bearer token:

```sh
curl -X POST http://<host>:9090/webhook/generic \
  -H "Authorization: Bearer $WEBHOOK_RECEIVER_TOKEN" \
  -d '{"repository":"ghcr.io/org/app"}'
```

## Example: Docker Hub

In your Docker Hub repository, add a webhook pointing at `https://<your-host>/webhook/dockerhub`. On each push, the operator rescans the matching policy and applies the update per its mode and write-back method.
