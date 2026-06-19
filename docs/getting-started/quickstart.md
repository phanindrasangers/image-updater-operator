# Quickstart

This walks through a minimal live-patch setup: a policy that tracks the latest 1.x release of nginx, and a Deployment that opts in.

Make sure the operator is [installed](installation.md) first.

## 1. Create an ImagePolicy

The policy defines what to scan and how to pick a tag.

```yaml title="imagepolicy.yaml"
apiVersion: images.saphire.com/v1alpha1
kind: ImagePolicy
metadata:
  name: nginx-stable
  namespace: default
spec:
  imageRepository: docker.io/library/nginx
  interval: 5m
  updateMode: Automatic
  policy:
    semver:
      range: ">=1.0.0 <2.0.0"
```

```sh
kubectl apply -f imagepolicy.yaml
```

## 2. Annotate a workload

Bind a container to the policy by name. The annotation key suffix is the container name.

```yaml title="deployment.yaml"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  namespace: default
  annotations:
    image-updater.saphire.com/policy.app: nginx-stable
spec:
  replicas: 1
  selector:
    matchLabels: { app: web }
  template:
    metadata:
      labels: { app: web }
    spec:
      containers:
        - name: app
          image: docker.io/library/nginx:1.0.0
```

```sh
kubectl apply -f deployment.yaml
```

With no `write-back` annotation, the default is **live**: the operator patches the running Deployment.

## 3. Watch it work

```sh
kubectl get imagepolicy nginx-stable -o wide
kubectl get deploy web -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'
kubectl describe deploy web | sed -n '/Events:/,$p'
```

Within an interval the policy resolves `status.latestTag`, the operator patches the container image, and Kubernetes rolls a new pod. You will see an `ImageUpdated` event on the Deployment.

## What changed

- The `ImagePolicy` status now carries the selected tag and a `Ready` condition.
- The Deployment's container image was bumped to the highest matching 1.x tag.
- A log line `live patch: updated workload image ...` was emitted by the controller.

## Next

- Switch to committing changes to Git instead of patching: [Write-back modes](../concepts/write-back.md) and [Git write-back](../reference/git-writeback.md).
- Tune selection: [ImagePolicy reference](../reference/imagepolicy.md).
- Require approval before updates, or run in dry-run: [Workload annotations](../reference/annotations.md).
