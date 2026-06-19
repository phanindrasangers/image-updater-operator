# Approvals and dry-run

By default the operator updates as soon as a newer tag is selected (`Automatic`). For sensitive workloads you can require a human to approve each update, or watch what would change without applying anything.

## Prerequisites

- [x] The operator [installed](../getting-started/installation.md).
- [x] A working policy and an annotated workload, as in the [Quickstart](../getting-started/quickstart.md).

The mode is set on the `ImagePolicy` (`spec.updateMode`) and can be overridden per workload with the `update-mode` annotation.

## Require approval

Set the policy (or the workload) to `Approval`. The operator records the candidate tag and emits an `ApprovalRequired` event instead of applying anything.

```yaml
spec:
  updateMode: Approval
```

Or per workload:

```yaml
metadata:
  annotations:
    image-updater.saphire.com/policy.app: app-stable
    image-updater.saphire.com/update-mode: Approval
```

When a new tag is available, check the event for the exact candidate:

```sh
kubectl describe deploy web | sed -n '/Events:/,$p'
# ApprovalRequired: container "app" has candidate 1.4.0 pending approval
#   (set image-updater.saphire.com/approve.app: "1.4.0")
```

Approve it by setting the annotation it names. The update then proceeds on the next reconcile, in whichever write-back mode the workload uses.

```sh
kubectl annotate deploy web \
  image-updater.saphire.com/approve.app=1.4.0 --overwrite
```

The approval is tag-specific: a later candidate requires a new approval, so you never silently roll past a version you vetted.

## Dry-run

Set `DryRun` to surface available updates as `UpdateAvailable` events while never touching the workload or Git. Useful for evaluating a policy before trusting it.

```yaml
metadata:
  annotations:
    image-updater.saphire.com/policy.app: app-stable
    image-updater.saphire.com/update-mode: DryRun
```

```sh
kubectl get events --field-selector involvedObject.name=web
# UpdateAvailable: container "app" could update to ghcr.io/org/app:1.4.0 (dry-run)
```

## Mode summary

| Mode | Behavior | Event |
|------|----------|-------|
| `Automatic` | Apply immediately | `ImageUpdated` / `ImageCommitted` |
| `Approval` | Hold until the candidate tag is approved | `ApprovalRequired` |
| `DryRun` | Never apply; report only | `UpdateAvailable` |

These work the same in both live and Git write-back modes.
