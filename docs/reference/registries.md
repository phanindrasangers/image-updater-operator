# Registry authentication

Set `spec.imageRepository` to the full repository (no tag) and provide credentials one of two ways:

1. A `kubernetes.io/dockerconfigjson` Secret referenced by `spec.registryRef.secretName` (the same shape as an imagePullSecret).
2. The ambient default keychain (host docker config, cloud credential helpers, workload identity) when no `secretName` is set. This is how short-lived cloud credentials are picked up.

The policy is identical across registries apart from `imageRepository` and credentials:

```yaml
apiVersion: images.saphire.com/v1alpha1
kind: ImagePolicy
metadata: { name: app-stable }
spec:
  imageRepository: <REPO>
  interval: 5m
  updateMode: Automatic
  registryRef:
    secretName: regcred     # omit to use ambient credentials
  policy:
    semver: { range: ">=1.0.0 <2.0.0" }
```

## Docker Hub

`imageRepository: docker.io/<org>/<repo>` (public repos need no secret).

```sh
kubectl create secret docker-registry regcred \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=<user> --docker-password=<token>
```

Use `https://index.docker.io/v1/` as the server, not the `hub.docker.com` web URL.

## GHCR

`imageRepository: ghcr.io/<org>/<repo>`. Use a PAT with `read:packages`.

```sh
kubectl create secret docker-registry regcred \
  --docker-server=ghcr.io \
  --docker-username=<github-user> --docker-password=<PAT>
```

## Nexus

`imageRepository: <nexus-host>:<port>/<repo>`, e.g. `nexus.example.com:8082/app`.

```sh
kubectl create secret docker-registry regcred \
  --docker-server=nexus.example.com:8082 \
  --docker-username=<user> --docker-password=<pass>
```

If Nexus is served over plain HTTP, set `spec.registryRef.insecure: true`.

## JFrog Artifactory

`imageRepository: <subdomain>.jfrog.io/<docker-repo>/<image>`.

```sh
kubectl create secret docker-registry regcred \
  --docker-server=<subdomain>.jfrog.io \
  --docker-username=<user> --docker-password=<identity-token>
```

## Amazon ECR

`imageRepository: <account>.dkr.ecr.<region>.amazonaws.com/<repo>`.

ECR tokens are short-lived (12h), so a static Secret is not ideal. Preferred options:

- On EKS, give the operator's ServiceAccount an IAM role (IRSA) with `ecr:GetAuthorizationToken`, `ecr:BatchGetImage`, and `ecr:DescribeImages`/`ListImages`. Leave `registryRef` unset; the ambient keychain resolves credentials from the pod's web identity.
- For a one-off, mint a token and store it as a Secret (refresh as needed):

```sh
TOKEN=$(aws ecr get-login-password --region <region>)
kubectl create secret docker-registry regcred \
  --docker-server=<account>.dkr.ecr.<region>.amazonaws.com \
  --docker-username=AWS --docker-password="$TOKEN"
```

## GCR / Artifact Registry and ACR

Leave `registryRef` unset and rely on workload identity (GKE) or managed identity (AKS); the ambient keychain handles `*.pkg.dev`, `gcr.io`, and `*.azurecr.io`.
