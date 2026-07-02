# Security Policy

## Supported versions

The most recent minor release receives security fixes. Users are expected to upgrade to the latest release to receive patches.

## Reporting a vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Instead, report privately using [GitHub Security Advisories](https://github.com/phanindrasangers/image-updater-operator/security/advisories/new) for this repository, or email **phanindra.sangers@gmail.com** with:

- A description of the vulnerability and its impact.
- Steps to reproduce, including affected versions.
- Any known mitigation or workaround.

You should expect an initial response within **5 business days**. We will work with you to understand and validate the issue, develop and test a fix, and coordinate a disclosure timeline before any public release notes or advisory are published.

## Scope

This policy covers the operator code in this repository, its Helm chart, and its container image build. It does not cover the security posture of registries, Git providers, or GitOps controllers (Argo CD, Flux) that the operator integrates with; report those to their respective projects.

## Credentials handled by the operator

The operator reads registry credentials (`kubernetes.io/dockerconfigjson` Secrets or ambient cloud credentials) and Git credentials (HTTPS tokens or SSH keys from Kubernetes Secrets) to do its job. It never logs credential values and only holds them in memory for the duration of a reconcile. If you find a way for credentials to leak into logs, events, or Git commits, please report it as a security issue rather than a regular bug.
