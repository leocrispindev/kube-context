# Docker Hub Packaging — Design

## Problem

Today the only way to run the k8s-agent MCP server is to clone the repository and
build it locally (`go build` or `go run`). This is friction for anyone who just
wants to run the server — especially for the HTTP/remote deployment case, where
a published container image is the natural distribution format.

## Goal

Package the server as a Docker image, hardened and documented well enough that
it can be published to Docker Hub as `leocrispindev/k8s-mcp` and pulled/run
without cloning the repo.

## Non-goals

- No CI/CD automation (GitHub Actions) for building/publishing — publishing is
  manual for now.
- No change to the three existing Kubernetes auth modes (kubeconfig, in-cluster,
  remote API server) — the image must support all three generically, none is
  privileged over another.
- No fix for the missing `LICENSE` file referenced by the README badge — tracked
  separately, out of scope here.

## Design

### 1. Dockerfile hardening

Keep the existing two-stage build (`golang:1.25-alpine` builder →
distroless runtime). Two changes to the runtime stage:

- Switch the base image from `gcr.io/distroless/static-debian12` (implicit
  `latest`, runs as root) to `gcr.io/distroless/static-debian12:nonroot`
  (uid/gid 65532). This is the only meaningful hardening the image needs —
  `static-debian12` already ships CA certificates and tzdata, which the
  client-go HTTPS client needs, so no other runtime dependency changes.
- Build the binary with `-ldflags="-s -w"` to strip debug symbols and shrink
  the image.

Alternatives considered:
- **Alpine runtime stage** — keeps a shell for in-container debugging, but
  larger attack surface and image size. Rejected: the server has no
  interactive debugging use case; `kubectl exec` isn't relevant since this
  isn't stateful.
- **scratch** — smallest possible image, but requires manually copying CA
  certificates and tzdata, and has zero tooling for support/diagnosis.
  Rejected: distroless nonroot already gets the security benefit without the
  extra maintenance burden.

No `HEALTHCHECK` instruction is added — when deployed to Kubernetes, liveness
and readiness probes are defined in the K8s manifest (hitting `GET /health`),
which is the actual orchestrator in play; a Docker-level `HEALTHCHECK` would
be redundant and distroless has no shell/wget to implement one with anyway.

### 2. Build context hygiene

Add a `.dockerignore` file so `COPY . .` in the builder stage doesn't pull in
unrelated repo content — `.git`, `docs/`, `*.md`, `.vscode/`, and the example
`wordpress-k8s.yaml` manifest. This shrinks the build context and keeps the
image reproducible regardless of what other files exist in the working tree.

### 3. Documentation

Add a new "Docker" section to the existing `README.md` (not a separate file —
single source of truth; since publishing is manual, this same content gets
pasted as the Docker Hub "Full Description" when the image is published).

Contents:

- `docker pull leocrispindev/k8s-mcp` as the no-clone alternative to the
  existing "Clone and build" quick start.
- `docker run` examples for all three auth modes:
  - Kubeconfig: `-v ~/.kube/config:/home/nonroot/.kube/config:ro -e KUBECONFIG=/home/nonroot/.kube/config`
  - In-cluster: `-e K8S_IN_CLUSTER=true` (documented for context; only
    meaningful when the container actually runs as a Pod with a
    ServiceAccount, not via plain `docker run`)
  - Remote API server: `-e K8S_API_SERVER=... -v /path/to/ca.crt:/ca.crt:ro -e K8S_CA_CERT_PATH=/ca.crt`
- A manual multi-arch build/publish command using `docker buildx`:
  `docker buildx build --platform linux/amd64,linux/arm64 -t leocrispindev/k8s-mcp:<version> --push .`

## Testing / validation

- `docker build -t leocrispindev/k8s-mcp:test .` succeeds.
- `docker run --rm leocrispindev/k8s-mcp:test` starts the HTTP server and
  `GET /health` responds `{"status":"ok"}` (verified via
  `docker run -p 8080:8080 ...` + `curl`).
- Confirm the container runs as non-root (`docker run --rm <image> id` is not
  available with no shell — instead inspect `docker inspect` for the `User`
  field, or rely on the base image's documented nonroot UID).
- `docker run --rm -e KUBECONFIG=/home/nonroot/.kube/config -v ~/.kube/config:/home/nonroot/.kube/config:ro -p 8080:8080 leocrispindev/k8s-mcp:test` against a real/local cluster, then call `list_namespaces` through the MCP HTTP endpoint to confirm the token round-trip still works end-to-end.
