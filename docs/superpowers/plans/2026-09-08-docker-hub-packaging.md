# Docker Hub Packaging Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Harden the existing Dockerfile (nonroot, stripped binary), add a `.dockerignore`, and document `docker run`/`docker build`/`docker buildx` usage in the README so the k8s-agent MCP server can be published to Docker Hub as `leocrispindev/k8s-mcp` and run without cloning/building the repo.

**Architecture:** No code changes — this is packaging and documentation only. Three independent, sequential changes: (1) `.dockerignore` to shrink the build context, (2) `Dockerfile` runtime-stage hardening (nonroot distroless base + stripped binary), (3) a new "Docker" section in `README.md` covering `docker pull`, `docker run` for all three existing auth modes, and a manual multi-arch `docker buildx` publish command.

**Tech Stack:** Docker (multi-stage build, `gcr.io/distroless/static-debian12:nonroot`), `docker buildx` for multi-arch, Go 1.25 (builder stage only, unchanged).

---

## Reference: spec

Full design rationale is in `docs/superpowers/specs/2026-09-08-docker-hub-packaging-design.md`. This plan implements it as-is; no open questions remain.

## Reference: current file contents

**`Dockerfile`** (current, unmodified):
```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-mcp ./cmd/main.go

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/k8s-mcp .

EXPOSE 8080

ENTRYPOINT ["/app/k8s-mcp"]
```

**`README.md`** relevant region — the new Docker section is inserted between the end of `## Running Remotely (HTTPS)` (ends at the `mcpServers` JSON block, line 229) and the start of `## Authentication` (line 231):

```
...
### 3. Connect from Cursor (remote)

```json
{
  "mcpServers": {
    "k8s-agent": {
      "url": "https://k8s-agent.example.com/mcp",
      "headers": {
        "Authorization": "Bearer <cloud-provider-token>"
      }
    }
  }
}
```
                                          <-- new "## Docker" section goes here
## Authentication
...
```

---

### Task 1: Add `.dockerignore`

**Files:**
- Create: `.dockerignore`

- [ ] **Step 1: Write the file**

Create `.dockerignore` at the repo root with:

```
.git
docs/
*.md
.vscode/
wordpress-k8s.yaml
```

- [ ] **Step 2: Verify it's picked up by the build**

Run:
```bash
docker build -t k8s-mcp-context-check . 2>&1 | grep -i "sending build context"
```
Expected: a line like `Sending build context to Docker daemon  XXXkB` where XX is small (well under 1MB — the repo has no large assets, so this mainly confirms the build didn't error and the ignored files aren't the bulk of what's sent). This also implicitly confirms the existing `Dockerfile` still builds correctly with the new `.dockerignore` in place.

- [ ] **Step 3: Commit**

```bash
git add .dockerignore
git commit -m "Add .dockerignore to trim Docker build context

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 2: Harden the Dockerfile (nonroot base + stripped binary)

**Files:**
- Modify: `Dockerfile`

- [ ] **Step 1: Update the build stage to strip debug symbols**

In `Dockerfile`, change:
```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-mcp ./cmd/main.go
```
to:
```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o k8s-mcp ./cmd/main.go
```

- [ ] **Step 2: Switch the runtime stage to the nonroot distroless image**

In `Dockerfile`, change:
```dockerfile
# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12
```
to:
```dockerfile
# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot
```

The rest of the file (`WORKDIR /app`, `COPY --from=builder`, `EXPOSE 8080`, `ENTRYPOINT`) is unchanged.

The full resulting `Dockerfile` should read:

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o k8s-mcp ./cmd/main.go

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/k8s-mcp .

EXPOSE 8080

ENTRYPOINT ["/app/k8s-mcp"]
```

- [ ] **Step 3: Build the image**

Run:
```bash
docker build -t leocrispindev/k8s-mcp:test .
```
Expected: build completes with `Successfully tagged leocrispindev/k8s-mcp:test` (or buildkit's equivalent final "naming to ... done" line). No errors.

- [ ] **Step 4: Verify the image runs as a non-root user**

Run:
```bash
docker inspect --format='{{.Config.User}}' leocrispindev/k8s-mcp:test
```
Expected: `nonroot` (the user distroless's `:nonroot` tag bakes into the image config — our Dockerfile doesn't set `USER` itself, it inherits this from the base image).

- [ ] **Step 5: Verify the HTTP server starts and `/health` responds**

Run:
```bash
docker run --rm -d --name k8s-mcp-smoketest -p 8080:8080 leocrispindev/k8s-mcp:test
sleep 1
curl -sf http://localhost:8080/health
docker logs k8s-mcp-smoketest
docker stop k8s-mcp-smoketest
```
Expected: `curl` prints `{"status":"ok"}`, and `docker logs` shows the `Kubernetes Agent MCP Server` banner with `MCP endpoint: POST http://localhost:8080/mcp`. Note: the container starts and serves `/health` even without a valid Kubernetes connection configured, because `adapter.InitClientProvider()` in `cmd/main.go` only builds a client from local defaults (kubeconfig path lookup) — it does not fail startup just because no cluster is reachable. This step only proves the container runs; it does not exercise Kubernetes connectivity (see Task 4).

- [ ] **Step 6: Commit**

```bash
git add Dockerfile
git commit -m "Harden Dockerfile: run as nonroot, strip debug symbols

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 3: Document Docker usage in the README

**Files:**
- Modify: `README.md:229-231` (insert a new `## Docker` section between the end of `## Running Remotely (HTTPS)` and the start of `## Authentication`)

- [ ] **Step 1: Insert the new section**

In `README.md`, find this exact text (end of the `## Running Remotely (HTTPS)` section, currently lines 216–231):

```markdown
### 3. Connect from Cursor (remote)

```json
{
  "mcpServers": {
    "k8s-agent": {
      "url": "https://k8s-agent.example.com/mcp",
      "headers": {
        "Authorization": "Bearer <cloud-provider-token>"
      }
    }
  }
}
```

## Authentication
```

Replace it with (adds the new `## Docker` section in between, keeps everything else identical):

```markdown
### 3. Connect from Cursor (remote)

```json
{
  "mcpServers": {
    "k8s-agent": {
      "url": "https://k8s-agent.example.com/mcp",
      "headers": {
        "Authorization": "Bearer <cloud-provider-token>"
      }
    }
  }
}
```

## Docker

A pre-built image is published to Docker Hub as [`leocrispindev/k8s-mcp`](https://hub.docker.com/r/leocrispindev/k8s-mcp) — no clone or build step required.

```bash
docker pull leocrispindev/k8s-mcp
```

The image runs the HTTP transport by default (`ENTRYPOINT ["/app/k8s-mcp"]`, no `--stdio` flag) and listens on port `8080`. Pick the `docker run` example that matches how you connect to your cluster — these map directly to the [Kubernetes Connection](#kubernetes-connection) methods below.

**Kubeconfig** (mount your existing `~/.kube/config` read-only):

```bash
docker run --rm -p 8080:8080 \
  -v ~/.kube/config:/home/nonroot/.kube/config:ro \
  -e KUBECONFIG=/home/nonroot/.kube/config \
  leocrispindev/k8s-mcp
```

**Remote API server** (connect directly to a cluster endpoint, mounting its CA cert):

```bash
docker run --rm -p 8080:8080 \
  -e K8S_API_SERVER="https://k8s-api.internal:6443" \
  -e K8S_CA_CERT_PATH=/ca.crt \
  -v /path/to/ca.crt:/ca.crt:ro \
  leocrispindev/k8s-mcp
```

**In-cluster** (only meaningful when the container itself runs as a Pod with a bound ServiceAccount — not applicable to a plain `docker run`):

```yaml
env:
  - name: K8S_IN_CLUSTER
    value: "true"
```

### Building and publishing the image

The image is published manually (multi-arch, via `buildx`):

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t leocrispindev/k8s-mcp:latest \
  -t leocrispindev/k8s-mcp:<version> \
  --push .
```

## Authentication
```

- [ ] **Step 2: Verify the README renders correctly**

Run:
```bash
grep -n "^## " README.md
```
Expected output includes, in this order, a new `## Docker` line between `## Running Remotely (HTTPS)` and `## Authentication`:
```
...
## Running Remotely (HTTPS)
## Docker
## Authentication
...
```

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "Document Docker usage: pull, run, and multi-arch publish

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 4: End-to-end validation against a real cluster

This task has no code/doc changes — it's a manual verification that the hardened, documented image actually works against a live Kubernetes cluster, using the exact `docker run` command just added to the README.

**Files:** none

- [ ] **Step 1: Get a short-lived cloud-provider token and start the container**

Using a real (or local, e.g. `kind`/`minikube`) cluster you have `~/.kube/config` access to:

```bash
TOKEN="$(aws eks get-token --cluster-name <your-cluster> --output json | jq -r '.status.token')"
# For a local kind/minikube cluster with no cloud IAM, any valid bearer token
# your cluster's RBAC accepts will do, or skip -H Authorization entirely if
# your kubeconfig context has no auth requirement.

docker run --rm -d --name k8s-mcp-e2e -p 8080:8080 \
  -v ~/.kube/config:/home/nonroot/.kube/config:ro \
  -e KUBECONFIG=/home/nonroot/.kube/config \
  leocrispindev/k8s-mcp:test
```

- [ ] **Step 2: Call `list_namespaces` through the MCP HTTP endpoint**

```bash
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_namespaces","arguments":{}}}'
```
Expected: a JSON-RPC response whose `result` contains a text-content block with a JSON body listing real namespace names from your cluster (e.g. `default`, `kube-system`). This confirms the full token round-trip — `tokenFromContext` → `authenticate` → `BearerTokenRoundTripper` → live Kubernetes API call — still works from inside the nonroot container with the kubeconfig-mount volume path documented in the README.

- [ ] **Step 3: Clean up**

```bash
docker stop k8s-mcp-e2e
docker rmi leocrispindev/k8s-mcp:test k8s-mcp-context-check
```

No commit for this task — it's verification only, confirming Tasks 1–3 together produce a working image.

---

## Self-Review Notes

- **Spec coverage:** Dockerfile hardening (nonroot + ldflags) → Task 2. `.dockerignore` → Task 1. README Docker section (pull, run for all 3 auth modes, buildx publish) → Task 3. End-to-end validation from the spec's "Testing / validation" section → Task 4. No spec section is without a task.
- **Placeholders:** none — every step has literal commands, file contents, or exact diffs.
- **Consistency:** the `/home/nonroot/.kube/config` mount path is used identically in Task 3 (README) and Task 4 (manual validation); the image tag `leocrispindev/k8s-mcp:test` from Task 2 is reused (not renamed) in Task 4.
