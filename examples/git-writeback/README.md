# Git write-back manual test

This verifies the operator selects a newer tag and commits it to a Git repo
instead of patching the live Deployment.

## 1. Prepare the config repo

Pick a GitHub repo you own (create an empty one if needed), then commit a copy
of `deployment.yaml` to it at the same path the annotation names:

```sh
# in a checkout of your config repo
mkdir -p examples/git-writeback
cp /path/to/deployment.yaml examples/git-writeback/deployment.yaml
git add examples/git-writeback/deployment.yaml && git commit -m "seed deployment" && git push
```

Edit the `git-repo` annotation in `deployment.yaml` to point at that repo.

## 2. Create the Git credential Secret

Use a GitHub PAT with `repo` scope as the password. Create it in the same
namespace as the workload (default):

```sh
kubectl create secret generic git-https \
  --from-literal=username=phanindrasangers \
  --from-literal=password='<github-pat>'
```

## 3. Apply policy and workload, run the operator

```sh
kubectl apply -f imagepolicy.yaml
kubectl apply -f deployment.yaml
make install   # CRDs, if not already installed
make run        # or: go run ./cmd/main.go --metrics-bind-address=0
```

## 4. Observe

The live Deployment stays at `phanindrasangers/nginx:v1`. The operator commits
the selected tag (`v2`) to your config repo.

```sh
kubectl describe deploy web | sed -n '/Events:/,$p'   # ImageCommitted event
kubectl get imagepolicy app-stable -o yaml             # resolved tag in status
git -C /path/to/config-repo pull --quiet && \
  git -C /path/to/config-repo log --oneline -1         # the operator's commit
```

Confirm the live object did not change:

```sh
kubectl get deploy web -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'
# still docker.io/phanindrasangers/nginx:v1
```
