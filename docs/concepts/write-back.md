# Write-back modes

A single annotation, `write-back`, chooses how a selected image is applied. The two modes are mutually exclusive per workload.

## Live (default)

```yaml
image-updater.saphire.com/write-back: live   # or omit entirely
```

The operator patches the running workload's container image in place. Kubernetes then rolls the workload. This is the simplest setup and is right when you are **not** running a GitOps controller.

- The live object is mutated directly.
- The operator records the last image it wrote in a `last-updated.<container>` annotation for idempotency.
- A `live patch: updated workload image ...` log line and an `ImageUpdated` event are emitted.

## Git

```yaml
image-updater.saphire.com/write-back: git
```

The operator never touches the live workload. Instead it clones a Git repository, edits the image reference in your YAML, commits, and pushes. A GitOps controller (Argo CD, Flux) then syncs the change to the cluster.

```text
new tag ─► operator commits to Git ─► Argo CD / Flux syncs ─► workload updated
```

This is the production GitOps pattern: Git stays the source of truth, and the cluster state is never mutated out of band.

!!! warning "Git mode does not update the cluster by itself"
    In Git mode the live workload stays on its current image until your GitOps controller applies the committed change. If you have no sync controller, the running pods will not change. Use live mode if you want the operator to update the cluster directly.

Git write-back supports three target types (Helm values, Kustomize, and plain manifests), customizable commit messages, and HTTPS or SSH credentials. See the [Git write-back reference](../reference/git-writeback.md).

## Choosing a mode

| | Live | Git |
|---|------|-----|
| Updates the cluster directly | yes | no (via GitOps) |
| Source of truth | cluster | Git |
| Needs a GitOps controller | no | yes |
| Audit trail | events | Git history + events |
| Best for | clusters without GitOps | Argo CD / Flux users |
