# Git write-back

When a workload sets `write-back: git`, the operator clones the repository, edits the image reference in your YAML, commits, and pushes. It never patches the live workload. Edits are made through the YAML node tree, so unrelated content and comments are preserved. There are no marker comments in your source.

## Targets

The `write-back-target` annotation is `kind:path`, where `kind` is one of:

=== "Helm values"

    ```yaml
    image-updater.saphire.com/write-back-target: helm:env/prod/values.yaml
    image-updater.saphire.com/helm.image-tag.app: image.tag
    image-updater.saphire.com/helm.image-name.app: image.repository
    ```

    Sets dotted keys in a values file. `helm.image-tag.<container>` is required and receives the tag; `helm.image-name.<container>` is optional and receives the repository. A numeric path segment indexes a list, so array-form values work:

    ```yaml
    # values.yaml -> images.0.tag
    image-updater.saphire.com/helm.image-tag.app: images.0.tag
    ```

=== "Kustomize"

    ```yaml
    image-updater.saphire.com/write-back-target: kustomize:overlays/prod
    ```

    Upserts an entry in the kustomization's `images` list whose `name` matches the policy repository, setting its `newTag`. A missing entry is added.

=== "Plain manifest"

    ```yaml
    image-updater.saphire.com/write-back-target: manifest:apps/web/deployment.yaml
    ```

    Rewrites every container `image:` field that references the policy repository, across all documents in the file.

## Commit message

The commit message is a Go template. When `git-commit-message` is unset, the default is:

```text
chore(images): update <name>

<container>: <old> -> <new> (<file>)
```

Customize it with the annotation:

```yaml
image-updater.saphire.com/git-commit-message: |
  ci: bump {{.Container}} to {{.Tag}}

  {{range .Changes}}- {{.File}}: {{.Image}}
  {{end}}
```

### Template context

| Field | Description |
|-------|-------------|
| `.Name`, `.Namespace`, `.Kind` | The workload being updated. |
| `.Changes` | List of edits, each with `.File`, `.Container`, `.Repository`, `.Tag`, `.Image`, `.OldImage`. |
| `.Container`, `.Repository`, `.Tag`, `.Image`, `.OldImage` | Convenience fields mirroring the first change, for the common single-image case. |

An invalid template records a `CommitTemplateError` warning and falls back to the default rather than blocking the write-back.

## Credentials

The credential `Secret` lives in the workload's namespace and is named by `git-secret`. Its keys depend on the clone URL scheme.

=== "HTTPS"

    ```sh
    kubectl create secret generic git-https \
      --from-literal=username=git \
      --from-literal=password=<personal-access-token>   # or key "token"
    ```

=== "SSH"

    ```sh
    kubectl create secret generic git-ssh \
      --from-file=identity=$HOME/.ssh/id_ed25519 \
      --from-file=known_hosts=$HOME/.ssh/known_hosts
    ```

Omit `git-secret` entirely for public HTTPS repositories (read works, but pushing still needs credentials).

## Events

| Reason | Meaning |
|--------|---------|
| `ImageCommitted` | A change was committed and pushed; includes the commit SHA. |
| `CloneError` | The repository could not be cloned. |
| `AuthError` | Credentials could not be loaded or were rejected. |
| `PushError` | The push was rejected (often a stale ref; retried next reconcile). |
| `WriteBackError` | The target file or key could not be edited. |
| `WriteBackMisconfigured` | Required annotations are missing or malformed. |

## Recommendation

Prefer a small, dedicated config repository as the write-back target rather than your application's source repo. Clones stay fast and cheap, the operator runs comfortably at low memory, and automated commits do not pollute your source history.
