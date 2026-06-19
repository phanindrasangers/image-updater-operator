# ImagePolicy reference

`ImagePolicy` (`images.saphire.com/v1alpha1`) defines a repository to scan and the rule that selects a tag.

## Spec

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `imageRepository` | string | yes | | Repository to scan, no tag. e.g. `docker.io/library/nginx`, `ghcr.io/org/app`, `<acct>.dkr.ecr.<region>.amazonaws.com/app`. Bound containers are set to `<imageRepository>:<selectedTag>`. |
| `policy` | object | yes | | The selection rule. Set exactly one of `semver`, `regex`, `numerical`, `alphabetical`. |
| `filterTags.pattern` | string (regex) | no | | Pre-filter applied to the tag list before the rule runs. |
| `interval` | duration | no | `5m` | Scan cadence. Clamped to a 30s minimum. |
| `updateMode` | enum | no | `Automatic` | `Automatic`, `Approval`, or `DryRun`. Overridable per workload. |
| `registryRef.secretName` | string | no | | `kubernetes.io/dockerconfigjson` Secret in the policy's namespace. Omit to use ambient credentials. |
| `registryRef.insecure` | bool | no | `false` | Allow plain HTTP. Use only for trusted in-cluster registries. |
| `suspend` | bool | no | `false` | Pause scanning and updates for this policy. |

## Status

| Field | Description |
|-------|-------------|
| `latestTag` | Tag selected at the last successful scan. |
| `latestImage` | Full `repository:tag` applied to bound containers. |
| `lastScanTime` | Timestamp of the last scan attempt. |
| `scannedTags` | Number of tags the registry returned. |
| `conditions[Ready]` | `True` after a successful scan; `False` with a reason on error (`AuthError`, `ScanError`, `PolicyError`, `NoMatch`). |

## Selection strategies

Set exactly one strategy under `policy`.

### semver

Selects the highest tag satisfying a semver range.

```yaml
policy:
  semver:
    range: ">=1.0.0 <2.0.0"
```

### numerical

Orders purely numeric tags and picks the first or last.

```yaml
policy:
  numerical:
    order: desc   # desc = highest, asc = lowest
```

### regex

Matches tags against a pattern, optionally extracting a sort key from a capture group. Combine with `numeric` to sort the key as a number.

```yaml
filterTags:
  pattern: '^v\d+$'
policy:
  regex:
    pattern: '^v(\d+)$'
    extract: "$1"     # use capture group 1 as the sort key
    numeric: true     # compare it as a number, not a string
    order: desc
```

`extract` only produces the sort key; the tag actually selected is the original (e.g. `v2`, not `2`). With multiple capture groups, reference them as `$1`, `$2`, or combine like `"$1.$2"`.

### alphabetical

Lexical ordering, useful for date-stamped tags that sort correctly as strings.

```yaml
policy:
  alphabetical:
    order: desc
```

## filterTags

A pre-filter that runs before the strategy, narrowing the candidate set:

```yaml
filterTags:
  pattern: '^nightly-'   # only consider tags beginning with "nightly-"
```

## Examples

```yaml
# Latest stable 1.x
policy: { semver: { range: ">=1.0.0 <2.0.0" } }

# Highest numeric build tag
policy: { numerical: { order: desc } }

# Newest date-stamped nightly
filterTags: { pattern: '^nightly-' }
policy:
  regex: { pattern: '^nightly-(\d{8})$', extract: "$1", numeric: true, order: desc }
```
