# Roadmap

This roadmap tracks direction at a high level. For granular, in-progress work see the [GitHub Issues](https://github.com/phanindrasangers/image-updater-operator/issues) and [milestones](https://github.com/phanindrasangers/image-updater-operator/milestones).

## Shipped

- Registry-agnostic tag scanning (any Docker Registry v2 endpoint: Docker Hub, GHCR, ECR, GCR, ACR, Nexus, JFrog Artifactory).
- Selection strategies: semver ranges, regex with capture-group ordering, numerical, and alphabetical, with an optional tag pre-filter.
- Two write-back methods per workload: live patch, or Git commit (Helm values, Kustomize images list, or plain manifest), for GitOps controllers (Argo CD, Flux) to sync.
- Update modes: Automatic, Approval (with per-tag approval), and DryRun.
- Multi-container support (main, init, and sidecar containers), each bound to its own policy.
- Webhook receiver for immediate re-scans (Docker Hub, Harbor, generic) alongside interval polling.
- Customizable Git commit message (Go template) and committer identity.
- Built-in read-only dashboard showing every policy and monitored workload.
- Helm chart, published to GHCR (OCI) and listed on Artifact Hub.

## In progress / near term

- Signature and provenance verification (cosign/Sigstore) before applying a selected tag.
- Additional notification sinks beyond Kubernetes events (Slack, generic webhook-out).
- Per-policy rate limiting and backoff tuning for large fleets of shared policies.
- Expanded conformance test suite against more registry implementations.

## Under consideration

- Multi-cluster policy fan-out (one policy driving updates across clusters).
- A CLI for local testing of selection rules against a live registry without applying a `ImagePolicy`.
- Structured audit log export for write-back history, independent of Git log mining.

## How to influence this roadmap

Open an issue describing the use case. Roadmap items move up when there's a concrete need behind them, not just a nice-to-have.
