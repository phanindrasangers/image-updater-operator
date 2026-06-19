# Image Updater Operator

A Kubernetes operator that keeps your workload container images current. It watches an external registry for new tags, selects one using a policy you define, and applies the update either by patching the live workload or by committing the change to Git for a GitOps controller to sync.

It is a generic, registry-agnostic alternative to tying image automation to a single GitOps tool: define an `ImagePolicy`, annotate a workload, and choose how updates land.

## Why

- **Any registry.** Anything that speaks the Docker Registry v2 API: Docker Hub, GHCR, ECR, GCR, ACR, Nexus, JFrog Artifactory.
- **Any workload.** Deployments, StatefulSets, DaemonSets, ReplicaSets, CronJobs, Jobs, and bare Pods.
- **Two write-back modes.** Patch the running workload directly, or commit to Git and let Argo CD or Flux apply it. Chosen per workload with one annotation.
- **Flexible selection.** Semver ranges, numeric ordering, regex extraction, or alphabetical, with an optional tag pre-filter.
- **Safe by default.** Per-policy `Automatic`, `Approval`, and `DryRun` modes, and an opt-in model: only annotated workloads are touched.
- **Observable.** A built-in read-only dashboard, Kubernetes events, and status conditions.

## How it fits together

```text
   ImagePolicy (CRD)                Workload (annotated)
   repository + rule        ┌──────────────────────────────┐
        │                   │  policy.<container>: <name>   │
        ▼                   │  write-back: live | git       │
   scan registry            └──────────────────────────────┘
        │                                  │
   select tag ──► status.latestImage ──────┤
                                           ▼
                          ┌────────────────────────────┐
                          │  live  → patch the workload │
                          │  git   → commit to a repo   │
                          └────────────────────────────┘
```

## Next steps

- [Installation](getting-started/installation.md) — install with Helm.
- [Quickstart](getting-started/quickstart.md) — a working policy and workload in minutes.
- [How it works](concepts/how-it-works.md) — the reconciliation model.
- [Write-back modes](concepts/write-back.md) — live patching vs Git write-back.

## Reference

- [ImagePolicy](reference/imagepolicy.md) — the CRD spec and selection strategies.
- [Workload annotations](reference/annotations.md) — every annotation the operator reads.
- [Git write-back](reference/git-writeback.md) — targets, commit templates, and credentials.
- [Registry authentication](reference/registries.md) — per-registry credential setup.
