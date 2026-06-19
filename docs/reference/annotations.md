# Workload annotations

Workloads opt in and configure behavior through annotations, all prefixed `image-updater.saphire.com/`. The prefix is omitted from the keys in the tables below.

## Binding and modes

| Annotation | Required | Default | Description |
|------------|----------|---------|-------------|
| `policy.<container>` | yes | | Bind a container to an `ImagePolicy` by name. The key suffix is the container name; init and sidecar containers are addressed the same way. e.g. `policy.app: nginx-stable`. |
| `update-mode` | no | policy's mode | Override the policy's `updateMode` for this workload: `Automatic`, `Approval`, or `DryRun`. |
| `approve.<container>` | when in Approval mode | | Approve a specific candidate tag for a container. Value is the tag to approve, e.g. `approve.app: "1.4.0"`. |

The operator also sets `last-updated.<container>` on the workload to record the last image it wrote (live mode only). Do not set this yourself.

## Write-back method

| Annotation | Required | Default | Description |
|------------|----------|---------|-------------|
| `write-back` | no | `live` | `live` patches the running workload; `git` commits to a repository. |

## Git write-back

Used only when `write-back: git`. See the [Git write-back reference](git-writeback.md) for full detail.

| Annotation | Required | Default | Description |
|------------|----------|---------|-------------|
| `git-repo` | for git | | Clone URL. HTTPS (`https://host/org/repo.git`) or SSH (`git@host:org/repo.git`). |
| `git-branch` | no | `main` | Branch to read and write. |
| `git-secret` | no | | Secret in the workload's namespace holding Git credentials. Omit for public HTTPS repos. |
| `write-back-target` | for git | | `helm:<file>`, `kustomize:<dir>`, or `manifest:<file>`, relative to the repo root. |
| `helm.image-tag.<container>` | for helm | | Dotted values key that receives the selected tag, e.g. `image.tag` or `images.0.tag`. A numeric segment indexes a list. |
| `helm.image-name.<container>` | no | | Dotted values key that receives the repository, e.g. `image.repository`. |
| `git-commit-message` | no | built-in | Go template for the commit message. |
| `git-author-name` | no | `image-updater-operator` | Committer name. |
| `git-author-email` | no | `image-updater@saphire.com` | Committer email. |

## Full example

```yaml
metadata:
  annotations:
    # bind the "app" container and require approval before applying
    image-updater.saphire.com/policy.app: app-stable
    image-updater.saphire.com/update-mode: Approval
    image-updater.saphire.com/approve.app: "1.4.0"

    # commit to Git instead of patching live
    image-updater.saphire.com/write-back: git
    image-updater.saphire.com/git-repo: https://github.com/org/app-config.git
    image-updater.saphire.com/git-branch: main
    image-updater.saphire.com/git-secret: git-https
    image-updater.saphire.com/write-back-target: helm:env/prod/values.yaml
    image-updater.saphire.com/helm.image-name.app: image.repository
    image-updater.saphire.com/helm.image-tag.app: image.tag
```
