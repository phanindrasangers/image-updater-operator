# Contributing

Thanks for considering a contribution to image-updater-operator. This guide covers how to set up a dev environment, the expected workflow, and what a good pull request looks like.

## Code of Conduct

Participation in this project is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting started

Requirements: Go 1.25 (see [Develop and run](README.md#develop-and-run) in the README for the toolchain note), [kubebuilder](https://book.kubebuilder.io/), and a cluster for testing ([kind](https://kind.sigs.k8s.io/) works well).

```sh
git clone https://github.com/phanindrasangers/image-updater-operator.git
cd image-updater-operator
make manifests generate   # regenerate CRDs and deepcopy after API changes
make test                 # unit tests + envtest
make install              # install CRDs into the current cluster
make run                  # run the manager locally against the current kube context
```

See [TESTING.md](TESTING.md) for end-to-end test scripts covering live patch, Git write-back, and per-registry setup.

## Reporting bugs

Open a [GitHub issue](https://github.com/phanindrasangers/image-updater-operator/issues/new) with:

- What you expected to happen and what happened instead.
- The operator version (or commit), Kubernetes version, and relevant `ImagePolicy`/workload annotations (redact registry credentials, repo URLs, and tokens).
- Relevant controller logs (`kubectl logs -n <namespace> deploy/<operator>`) and events (`kubectl describe`).

Report suspected **security vulnerabilities** privately per [SECURITY.md](SECURITY.md), not as a public issue.

## Proposing features

Open an issue describing the use case before sending a large pull request, especially for anything touching the annotation contract, the CRD schema, or the Git write-back format. Small fixes, docs, and tests can go straight to a PR.

## Making a change

1. Fork the repo and create a branch from `main`.
2. Make your change. Keep it focused; unrelated cleanup belongs in a separate PR.
3. Add or update tests. `internal/policy`, `internal/gitwriteback`, and `internal/controller` all have unit tests; controller changes are also covered by envtest.
4. Run the full check locally before opening a PR:
   ```sh
   make manifests generate
   make test
   go vet ./...
   ```
5. If you changed the CRD, RBAC, or webhook config, run `make manifests` and commit the regenerated files under `config/`.
6. Open a pull request against `main` with a clear description of what changed and why.

## Commit messages

Write commit messages that explain *why*, not just what changed; the diff already shows what changed. No fixed prefix convention is enforced, but keeping the subject line under about 72 characters is appreciated.

## Documentation

User-facing docs live in [docs/](docs/) and are published with MkDocs; see [README.md](README.md) and the [docs site](https://phanindrasangers.github.io/image-updater-operator/) for the rendered version. Update the relevant page alongside any behavior change (new annotation, new selection strategy, new write-back target, etc.).

## Release process

Maintainers cut releases by pushing a semver tag; see [Releasing](README.md#releasing) in the README. Contributors do not need to do anything release-related.
