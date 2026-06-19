# GitOps write-back

This guide walks through having the operator commit image updates to Git instead of patching the cluster, so a GitOps controller rolls them out. If you just want the operator to update the cluster directly, use the [Quickstart](../getting-started/quickstart.md) (live mode) instead.

New to the two modes? Read [Write-back modes](../concepts/write-back.md) first.

## Prerequisites

Before you start, make sure you have:

- [x] The operator [installed](../getting-started/installation.md).
- [x] A **config repository** in Git holding the manifests or Helm values that define your workload's image (this is what the operator edits). Prefer a small, dedicated repo over your application's source repo.
- [x] **Push credentials** for that repo: a personal access token (HTTPS) or an SSH key. You will store these in a Secret in step 1.
- [x] A **GitOps controller** (Argo CD or Flux) syncing that repo to your cluster. Without one, the commit still happens but the running workload will not change. See [why](../concepts/write-back.md#git).

## 1. Create the Git credentials Secret

The Secret lives in the **same namespace as the workload** and is referenced later by the `git-secret` annotation. Its keys depend on the clone URL scheme.

=== "HTTPS (token)"

    ```sh
    kubectl create secret generic git-https \
      --from-literal=username=git \
      --from-literal=password=<personal-access-token>
    ```

=== "SSH (key)"

    ```sh
    kubectl create secret generic git-ssh \
      --from-file=identity=$HOME/.ssh/id_ed25519 \
      --from-file=known_hosts=$HOME/.ssh/known_hosts
    ```

!!! note
    For a **public** repo you can omit the Secret for the read/scan, but pushing the commit still requires credentials, so you will normally create one anyway.

## 2. Create the ImagePolicy

Same as any policy: it defines what to scan and how to pick a tag. The `imageRepository` must match the image referenced in your config repo's YAML.

```yaml
apiVersion: images.saphire.com/v1alpha1
kind: ImagePolicy
metadata:
  name: app-stable
  namespace: default
spec:
  imageRepository: ghcr.io/org/app
  interval: 5m
  updateMode: Automatic
  policy:
    semver:
      range: ">=1.0.0 <2.0.0"
```

## 3. Annotate the workload for Git write-back

The annotations tell the operator to commit instead of patch, which repo and branch to use, which Secret to authenticate with, and which file to edit.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  namespace: default
  annotations:
    # bind the container to the policy
    image-updater.saphire.com/policy.app: app-stable
    # commit to Git instead of patching live
    image-updater.saphire.com/write-back: git
    image-updater.saphire.com/git-repo: https://github.com/org/app-config.git
    image-updater.saphire.com/git-branch: main
    image-updater.saphire.com/git-secret: git-https      # the Secret from step 1
    # what to edit in the repo
    image-updater.saphire.com/write-back-target: helm:env/prod/values.yaml
    image-updater.saphire.com/helm.image-name.app: image.repository
    image-updater.saphire.com/helm.image-tag.app: image.tag
spec:
  # ... pod template ...
```

The `write-back-target` chooses the file layout. Swap it for `kustomize:<dir>` or `manifest:<file>` to edit those instead. Full options are in the [Git write-back annotations](../reference/git-writeback.md) reference.

## 4. Verify the commit

Once the policy selects a tag, the operator clones, edits, commits, and pushes.

```sh
# event on the workload, includes the commit SHA
kubectl describe deploy web | sed -n '/Events:/,$p'   # look for ImageCommitted

# the operator's commit in your repo
git -C <your-config-repo-checkout> pull --quiet
git -C <your-config-repo-checkout> log --oneline -1
```

The live Deployment stays on its current image. That is expected: in Git mode the operator only writes to Git.

## 5. Let GitOps roll it out

Your Argo CD or Flux app, pointed at the config repo, detects the commit and applies the new image to the cluster. That step closes the loop and updates the running pods.

## Customizing the commit

Set a commit message template and author with annotations:

```yaml
image-updater.saphire.com/git-commit-message: |
  ci: bump {{.Container}} to {{.Tag}}

  {{range .Changes}}- {{.File}}: {{.Image}}
  {{end}}
image-updater.saphire.com/git-author-name: ci-bot
image-updater.saphire.com/git-author-email: ci-bot@example.com
```

See the [Git write-back annotations](../reference/git-writeback.md#commit-message) reference for the full template context.

## Troubleshooting

- **Commit happened but pods unchanged** — expected in Git mode; your GitOps controller applies it. See [Troubleshooting](../troubleshooting.md#it-is-in-git-mode).
- **`CloneError` / `AuthError` / `PushError` events** — check the repo URL, the Secret keys, and that the token can push. See [Troubleshooting](../troubleshooting.md).
- **`OOMKilled`** — the clone exceeded the memory limit; use a smaller repo or raise `resources.limits.memory`.
