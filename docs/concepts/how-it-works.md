# How it works

The operator runs two cooperating controllers.

## The ImagePolicy controller

Each `ImagePolicy` names a repository and a selection rule. On an interval the controller:

1. Loads registry credentials (from a referenced Secret, or the ambient keychain).
2. Lists the repository's tags.
3. Applies the optional `filterTags` pre-filter, then the selection rule.
4. Records the winner in `status.latestTag` / `status.latestImage` and sets the `Ready` condition.

It scans at most once per `interval`. Reconciliation can be triggered far more often than the interval (status writes, cache resyncs), so the controller gates on the elapsed time since the last scan and only rescans when due, when the spec changes, or on first run. A failed scan is throttled the same way so a broken policy never hammers the registry.

## The workload controller

A separate controller watches annotated workloads and the policies they reference. When a policy's `status.latestImage` differs from what a bound container runs, it:

1. Checks the effective update mode (`Automatic`, `Approval`, `DryRun`).
2. Applies the change according to the workload's write-back method.

Only the top-level workload is acted on. A Deployment propagates its annotations to its ReplicaSet and Pods, so the controller skips controller-owned objects to avoid acting on the same logical workload twice.

## Opt-in model

Nothing is changed unless a workload carries a `policy.<container>` annotation. The container name in the annotation suffix binds a specific container to a specific `ImagePolicy`, so init containers and sidecars are addressed the same way as the main container.

## Supported workload kinds

| Kind | Live patch | Notes |
|------|------------|-------|
| Deployment | yes | |
| StatefulSet | yes | |
| DaemonSet | yes | |
| ReplicaSet | yes | Standalone only; managed ones are skipped. |
| CronJob | yes | |
| Pod | yes | Standalone only. |
| Job | reported only | Job pod templates are immutable; the update is reported as an event but not applied. |

## Update modes

- **Automatic** — apply the update as soon as it is selected.
- **Approval** — hold until a human approves the candidate tag with an annotation. The operator emits `ApprovalRequired` with the exact annotation to set.
- **DryRun** — never apply; emit an `UpdateAvailable` event so you can see what would change.

The mode is set on the policy and can be overridden per workload. See [Workload annotations](../reference/annotations.md).

## Status and events

- `ImagePolicy.status` exposes `latestTag`, `latestImage`, `lastScanTime`, `scannedTags`, and a `Ready` condition with a reason on failure.
- Workloads receive Kubernetes events: `ImageUpdated`, `ImageCommitted`, `UpdateAvailable`, `ApprovalRequired`, and warnings such as `CloneError` or `PushError`.
