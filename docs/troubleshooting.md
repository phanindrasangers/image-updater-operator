# Troubleshooting

## The workload is not updating

First check the policy resolved a tag:

```sh
kubectl get imagepolicy <name> -o jsonpath='{.status}' | jq
```

- `latestImage` empty or `Ready=False`: the scan is failing. Check the condition `reason` (`AuthError`, `ScanError`, `PolicyError`, `NoMatch`) and the controller logs.
- `latestImage` set but the workload is unchanged: see the cases below.

### It is in Git mode

In `write-back: git` mode the operator commits to Git and never patches the live workload. The running pods only change once your GitOps controller (Argo CD, Flux) applies the commit. If you have no sync controller, use `write-back: live` instead. See [Write-back modes](concepts/write-back.md).

Confirm the commit happened:

```sh
kubectl describe deploy <name> | sed -n '/Events:/,$p'   # look for ImageCommitted
```

### The mode is holding it

- `Approval`: the operator emits `ApprovalRequired` and waits until `approve.<container>` names the candidate tag.
- `DryRun`: it only emits `UpdateAvailable` and never applies.

### The container name does not match

`policy.<container>` must match a container name in the pod spec exactly. A mismatch emits `ContainerNotFound`.

### It is a Job

Job pod templates are immutable. The operator reports the available update as an event but cannot apply it; recreate the Job to pick up the new image.

## CrashLoopBackOff / OOMKilled

Git write-back clones the repository in memory. A large target repository can exceed the pod's memory limit and the container is `OOMKilled` mid-clone. Raise the limit:

```sh
helm upgrade image-updater ... --set resources.limits.memory=1Gi
```

Better, point write-back at a small, dedicated config repository. See the [Git write-back](reference/git-writeback.md) recommendation.

## PushError: cannot lock ref / stale ref

A push was rejected because the remote branch moved since the clone. This is transient and self-heals on the next reconcile, which clones fresh. If it persists, another writer is committing to the same branch concurrently.

## Scans happen too often / too rarely

The operator scans at most once per `interval` (minimum 30s). If you set a very small interval it is clamped. Each interval tick produces exactly one scan and one `scanned repository` log line. If you do not see scans, confirm the policy is not `suspend: true` and the controller holds the leader lease:

```sh
kubectl get lease -n image-updater-system
kubectl logs -n image-updater-system deploy/image-updater-image-updater-operator
```

## Registry authentication fails

- `AuthError` on the policy: the referenced Secret is missing or has the wrong shape. It must be `kubernetes.io/dockerconfigjson`. See [Registry authentication](reference/registries.md).
- For cloud registries (ECR, GCR, ACR), prefer workload identity over static Secrets and leave `registryRef` unset.

## Verifying the operator is healthy

```sh
kubectl get pods -n image-updater-system
kubectl logs -n image-updater-system deploy/image-updater-image-updater-operator --tail=50
```

The pod should be `Running` with low restarts, and the logs should show periodic `scanned repository` lines at your configured interval.
