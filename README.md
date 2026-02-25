![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![MCP](https://img.shields.io/badge/MCP-Model%20Context%20Protocol-blueviolet)
![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?logo=kubernetes&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Open Source](https://img.shields.io/badge/Open%20Source-%E2%9D%A4-red)

# Kubernetes MCP Server

A Kubernetes agent that exposes cluster operations as **MCP (Model Context Protocol)** tools. It enables AI assistants and MCP-compatible clients to manage Kubernetes resources through a secure, multi-tenant interface powered by JWT authentication and Kubernetes RBAC impersonation.

## Architecture

```
┌──────────────┐       ┌──────────────────┐       ┌────────────────────┐
│  MCP Client  │──────▶│  MCP Server      │──────▶│  Kubernetes API    │
│  (AI / CLI)  │ HTTP  │  (Gin + MCP-Go)  │ RBAC  │  (Impersonation)   │
│              │ stdio │                  │       │                    │
└──────────────┘       └──────────────────┘       └────────────────────┘
                              │
                        JWT Validation
                        User Extraction
                        Group Mapping
```

The server follows clean architecture with clear separation between:

- **Core** — Business logic services and DTOs for each Kubernetes resource
- **Infrastructure** — MCP tool registration, Kubernetes client adapters, JWT auth

## Available MCP Tools

The server exposes **37 tools** across 7 Kubernetes resource types. Every tool call requires a valid JWT `token` for authentication.

### Namespace

| Tool | Description |
|------|-------------|
| `list_namespaces` | List all accessible namespaces |
| `get_namespace` | Get details of a specific namespace |
| `create_namespace` | Create a new namespace |
| `update_namespace` | Update namespace labels and annotations |
| `delete_namespace` | Delete a namespace |

### Pod

| Tool | Description |
|------|-------------|
| `list_pods` | List all pods in a namespace |
| `get_pod` | Get details of a specific pod |
| `create_pod` | Create a new pod |
| `update_pod` | Update pod labels and/or image |
| `delete_pod` | Delete a pod |
| `restart_pod` | Restart a pod (delete for controller recreation) |

### Deployment

| Tool | Description |
|------|-------------|
| `list_deployments` | List all deployments in a namespace |
| `get_deployment` | Get details of a specific deployment |
| `create_deployment` | Create a new deployment |
| `update_deployment` | Update replicas, images, or labels |
| `delete_deployment` | Delete a deployment |
| `get_rollout_status` | Get rollout status of a deployment |
| `toggle_pause_deployment` | Pause or resume a deployment rollout |

### ConfigMap

| Tool | Description |
|------|-------------|
| `list_configmaps` | List all configmaps in a namespace |
| `get_configmap` | Get details of a specific configmap |
| `create_configmap` | Create a new configmap |
| `update_configmap` | Update configmap data and/or labels |
| `delete_configmap` | Delete a configmap |

### Service

| Tool | Description |
|------|-------------|
| `list_services` | List all services in a namespace |
| `get_service` | Get details of a specific service |
| `create_service` | Create a new Kubernetes service |
| `update_service` | Update service type, selector, or labels |
| `delete_service` | Delete a service |

### Ingress

| Tool | Description |
|------|-------------|
| `list_ingresses` | List all ingresses in a namespace |
| `get_ingress` | Get details of a specific ingress |
| `create_ingress` | Create a new ingress |
| `update_ingress` | Update ingress class, labels, or annotations |
| `delete_ingress` | Delete an ingress |

### NetworkPolicy

| Tool | Description |
|------|-------------|
| `list_network_policies` | List all network policies in a namespace |
| `get_network_policy` | Get details of a specific network policy |
| `create_network_policy` | Create a new network policy |
| `update_network_policy` | Update network policy labels |
| `delete_network_policy` | Delete a network policy |

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | Yes | — | Secret key for JWT validation (HS256) |
| `PORT` | No | `8080` | HTTP server port |

### Kubernetes Connection

Choose **one** of the following methods:

#### Kubeconfig (default)

Uses `~/.kube/config` automatically. Override with:

```bash
export KUBECONFIG="/path/to/your/kubeconfig"
```

#### Remote API Server

```bash
export K8S_API_SERVER="https://k8s.example.com:6443"
export K8S_BASE_TOKEN="your-service-account-token"
export K8S_CA_CERT_PATH="/path/to/ca.crt"
export K8S_INSECURE="false"  # set to "true" to skip TLS verification
```

#### In-Cluster

```bash
export K8S_IN_CLUSTER="true"
```

## Quick Start with Cursor

The fastest way to get the k8s-agent running locally inside **Cursor** as an MCP server.

### Prerequisites

- Go 1.25+
- Access to a Kubernetes cluster (local via `~/.kube/config`, or remote)
- A `JWT_SECRET` (any string for development, e.g. `"test"`)

### 1. Clone and build

```bash
git clone https://github.com/<your-username>/kube-context.git
cd k8s-agent-new
go build -o k8s-mcp cmd/main.go
```

This produces a `k8s-mcp` binary in the project root.

### 2. Configure Cursor

Create or edit the file **`.cursor/mcp.json`** (in your home directory or project root) and add the k8s-agent server:

```json
{
  "mcpServers": {
    "k8s-mcp": {
      "command": "/absolute/path/to/kube-context/k8s-mcp",
      "args": ["--stdio"],
      "env": {
        "JWT_SECRET": "your-secret-key"
      }
    }
  }
}
```

> **Tip:** Replace `/absolute/path/to/k8s-agent-new/k8s-mcp` with the actual path to the built binary. For example:
> ```
> /Users/me/workspace/k8s-agent-new/k8s-mcp
> ```

### 3. Restart Cursor

After saving `mcp.json`, restart Cursor (or reload the window). The k8s-agent tools will appear automatically in the MCP tool list. You can verify by asking the AI assistant to call `list_namespaces`.

### 4. Passing the authentication token

> **Important:** The k8s-agent uses JWT-based authentication. The token is **not** sent via HTTP headers — it must be passed as a **`token` parameter** in every MCP tool call. This is how authentication works locally: each tool expects a `token` field in its input arguments.

When using Cursor, include the token directly in your prompt so the AI assistant sends it on every call. For example:

```
Use the following token for all k8s-agent tool calls:
eyJhbGciOiJIUzI1NiIs...
```

The AI assistant will then include `"token": "eyJhbGciOiJIUzI1NiIs..."` as a parameter in every tool invocation automatically.

You can generate a token using the included `jwt-gen` CLI:

```bash
go run cmd/jwt-gen/main.go
```

### 5. (Optional) Auto-inject the token

To avoid passing the JWT token manually on every tool call, you can add it to the environment and create a Cursor rule to inject it automatically:

**a)** Add `K8S_TOKEN` to the MCP config:

```json
{
  "mcpServers": {
    "k8s-mcp": {
      "command": "/absolute/path/to/k8s-agent-new/k8s-mcp",
      "args": ["--stdio"],
      "env": {
        "JWT_SECRET": "your-secret-key",
        "K8S_TOKEN": "eyJhbGciOiJIUzI1NiIs..."
      }
    }
  }
}
```

**b)** Create a Cursor rule at `.cursor/rules/k8s-agent.mdc`:

```markdown
---
description: Rules for k8s-agent MCP tools
globs: *
---

When using any k8s-agent MCP tool, always pass the token from the
environment variable K8S_TOKEN. Example:

  token: ${K8S_TOKEN}

Never ask the user for the token — use the environment variable directly.
```

---

## Running Locally (standalone)

If you prefer running the server **outside** of Cursor (e.g. for testing with `curl`):

### HTTP Mode (default)

```bash
export JWT_SECRET="your-secret-key"

go run cmd/main.go
```

The server starts at `http://localhost:8080` with:

- **MCP endpoint** — `POST /mcp`
- **Health check** — `GET /health`

### Stdio Mode

For MCP clients that communicate over stdin/stdout:

```bash
export JWT_SECRET="your-secret-key"

go run cmd/main.go --stdio
```

## Running Remotely (HTTPS)

For production deployments, place the server behind a reverse proxy with TLS termination.

### 1. Build the binary

```bash
go build -o k8s-agent cmd/main.go
```

### 2. Deploy with a reverse proxy

Use **nginx**, **Traefik**, **Caddy**, or any reverse proxy that handles TLS. Example nginx snippet:

```nginx
server {
    listen 443 ssl;
    server_name k8s-agent.example.com;

    ssl_certificate     /etc/ssl/certs/your-cert.pem;
    ssl_certificate_key /etc/ssl/private/your-key.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. Configure environment

```bash
export JWT_SECRET="<strong-production-secret>"
export K8S_API_SERVER="https://k8s-api.internal:6443"
export K8S_BASE_TOKEN="<service-account-token>"
export K8S_CA_CERT_PATH="/etc/k8s/ca.crt"
export PORT="8080"
```

## JWT Structure and Kubernetes Impersonation

The k8s-agent **does not** communicate with the Kubernetes API as a single privileged identity. Instead, it uses [**User Impersonation**](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#user-impersonation) — a native Kubernetes mechanism that allows a service account to perform API calls *on behalf of* the user described in the JWT.

This means the permissions of each request are determined by the **Kubernetes RBAC rules bound to the impersonated user/groups**, not by the service account running the agent. This is the foundation for multi-tenant security: different users see only what their RBAC policies allow.

### How it works

1. The client sends a JWT `token` in every MCP tool call.
2. The k8s-agent validates the signature (HS256) and expiration.
3. Claims are extracted and mapped to Kubernetes impersonation headers:

| JWT Claim | Kubernetes Header | Description |
|-----------|-------------------|-------------|
| `sub` (or `user`, `username`, `preferred_username`, `email`) | `Impersonate-User` | The identity used for RBAC evaluation. First non-empty value wins, checked in the listed order. |
| `groups` | `Impersonate-Group` | List of groups. Kubernetes RBAC can bind roles to groups for broad access control. |
| `extra` (or `extras`) | `Impersonate-Extra-<key>` | Arbitrary key-value pairs forwarded as extra impersonation attributes. |
| `exp` | *(validation only)* | Standard JWT expiration (Unix timestamp). The token is rejected if expired. |

4. The Kubernetes API server then evaluates the request as if it came directly from the impersonated user, applying all RBAC rules normally.

### Required JWT payload

The **only mandatory field** is the user identity — at least one of: `sub`, `user`, `username`, `preferred_username`, or `email`. Everything else is optional but recommended for fine-grained RBAC.

#### Minimal payload

```json
{
  "sub": "jane@example.com",
  "exp": 1893456000
}
```

This impersonates user `jane@example.com` with no groups and no extras.

#### Full payload

```json
{
  "sub": "jane@example.com",
  "groups": ["team-backend", "developers"],
  "extra": {
    "department": ["engineering"],
    "cost-center": ["cc-1234"]
  },
  "exp": 1893456000
}
```

This impersonates `jane@example.com` as a member of `team-backend` and `developers`, with additional extra attributes forwarded to the API server.

### JWT header

The token must use the **HS256** algorithm and be signed with the secret defined in `JWT_SECRET`:

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Why impersonation?

| Concern | How impersonation addresses it |
|---------|-------------------------------|
| **Multi-tenancy** | Each user gets only the permissions their RBAC policies grant — no shared super-admin. |
| **Auditability** | Kubernetes audit logs record the impersonated user, not the service account, making it clear *who* performed each action. |
| **Least privilege** | The agent's service account only needs the `impersonate` verb; actual resource permissions are defined per user/group via standard RBAC. |
| **No credential leakage** | Users never receive a kubeconfig or service account token — the JWT is the only credential, and it can have short TTLs. |

### RBAC setup for impersonation

The service account (or kubeconfig user) running the k8s-agent needs permission to impersonate. Template manifests are provided in `deploy/rbac/`:

```yaml
# ClusterRole — allows impersonating a specific user
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: impersonate-jane
rules:
  - apiGroups: [""]
    resources: ["users"]
    resourceNames: ["jane@example.com"]
    verbs: ["impersonate"]
  - apiGroups: [""]
    resources: ["groups"]
    verbs: ["impersonate"]
```

```yaml
# ClusterRoleBinding — binds the role to the agent's service account
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: impersonate-jane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: impersonate-jane
subjects:
  - kind: ServiceAccount
    name: k8s-agent
    namespace: default
```

Then, bind regular Roles/ClusterRoles to `jane@example.com` (or her groups) to control what resources she can actually access.

## Authentication — Sending the Token

Every MCP tool call requires a valid JWT in the `token` parameter. The k8s-agent **only validates** tokens — it never generates them. Token generation is handled by an external service or CLI.

### Sending the Token Locally (HTTP)

When running the MCP server locally in HTTP mode (`http://localhost:8080/mcp`), the token is sent **inside the MCP tool call arguments**, not as an HTTP header.

Every tool expects a `token` field in its input. For example, to list namespaces:

```json
{
  "method": "tools/call",
  "params": {
    "name": "list_namespaces",
    "arguments": {
      "token": "eyJhbGciOiJIUzI1NiIs..."
    }
  }
}
```

Full `curl` example against the local server:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d "{
    \"jsonrpc\": \"2.0\",
    \"id\": 1,
    \"method\": \"tools/call\",
    \"params\": {
      \"name\": \"list_namespaces\",
      \"arguments\": {
        \"token\": \"$TOKEN\"
      }
    }
  }"
```

> **Note:** When using an MCP-compatible client (like Cursor with stdio mode), the client handles the MCP protocol automatically — you only need to configure the `JWT_SECRET` environment variable. The AI assistant will include the token in each tool call for you.

### Sending the Token via HTTPS (Remote / Production)

For production deployments behind a TLS reverse proxy (e.g. `https://k8s-agent.example.com/mcp`), the mechanism is **exactly the same** — the token travels inside the MCP tool call arguments:

```bash
curl -X POST https://k8s-agent.example.com/mcp \
  -H "Content-Type: application/json" \
  -d "{
    \"jsonrpc\": \"2.0\",
    \"id\": 1,
    \"method\": \"tools/call\",
    \"params\": {
      \"name\": \"list_pods\",
      \"arguments\": {
        \"token\": \"$TOKEN\",
        \"namespace\": \"production\"
      }
    }
  }"
```

The difference from local usage is:

| | Local (HTTP) | Remote (HTTPS) |
|---|---|---|
| **URL** | `http://localhost:8080/mcp` | `https://k8s-agent.example.com/mcp` |
| **TLS** | No (development only) | Yes (reverse proxy handles TLS) |
| **Token delivery** | `token` field in tool arguments | `token` field in tool arguments |
| **Token source** | `jwt-gen` CLI or local token-issuer | Your organization's auth service, Keycloak, or the token-issuer |

> **Security:** In production, always use HTTPS so the JWT is encrypted in transit. The token is embedded in the request body, so TLS ensures it cannot be intercepted.

## MCP Client Configuration

### Cursor / AI Clients (stdio)

> See [Quick Start with Cursor](#quick-start-with-cursor) for a complete step-by-step guide.

There are two ways to configure Cursor:

#### Option A — Pre-built binary (recommended)

Build once with `go build -o k8s-mcp cmd/main.go`, then point Cursor to the binary:

```json
{
  "mcpServers": {
    "k8s-mcp": {
      "command": "/absolute/path/to/k8s-agent-new/k8s-mcp",
      "args": ["--stdio"],
      "env": {
        "JWT_SECRET": "your-secret-key"
      }
    }
  }
}
```

#### Option B — `go run` (no build step)

Cursor compiles and runs on every startup — slower but requires no manual build:

```json
{
  "mcpServers": {
    "k8s-agent": {
      "command": "go",
      "args": ["run", "cmd/main.go", "--stdio"],
      "env": {
        "JWT_SECRET": "your-secret-key"
      }
    }
  }
}
```

#### Auto-injecting the Token

Since every tool call requires a `token` argument, you can automate this by adding the token to the `env` block and creating a Cursor rule.

**1.** Add `K8S_TOKEN` to the MCP config:

```json
{
  "mcpServers": {
    "k8s-mcp": {
      "command": "/absolute/path/to/k8s-agent-new/k8s-mcp",
      "args": ["--stdio"],
      "env": {
        "JWT_SECRET": "your-secret-key",
        "K8S_TOKEN": "eyJhbGciOiJIUzI1NiIs..."
      }
    }
  }
}
```

**2.** Create a Cursor rule at `.cursor/rules/k8s-agent.mdc`:

```markdown
---
description: Rules for k8s-agent MCP tools
globs: *
---

When using any k8s-agent MCP tool, always pass the token from the
environment variable K8S_TOKEN. Example:

  token: ${K8S_TOKEN}

Never ask the user for the token — use the environment variable directly.
```

This way the AI assistant will automatically include the token in every tool call without manual input.

### HTTP Clients

Point your MCP client to:

```
POST http://localhost:8080/mcp
```

Or for remote deployments:

```
POST https://k8s-agent.example.com/mcp
```

## Project Structure

```
k8s-agent-new/
├── cmd/
│   ├── main.go                  # MCP server entry point
│   ├── jwt-gen/                 # CLI tool for JWT generation
│   └── token-issuer/            # HTTP service for JWT issuance
├── internal/
│   ├── core/
│   │   ├── dto/                 # Data Transfer Objects
│   │   ├── gateway/             # Interface definitions
│   │   └── service/             # Business logic per resource
│   └── infrastructure/
│       ├── adapter/             # Kubernetes client & impersonation
│       ├── auth/                # JWT validation
│       └── mcp/                 # MCP server & tool registration
│           └── tools/           # One file per resource type
├── deploy/
│   └── rbac/                    # Multi-tenant RBAC templates
├── go.mod
└── go.sum
```

## Contributing

This is an **open-source** project and contributions are welcome!

### Getting Started

1. **Fork** the repository
2. **Clone** your fork:
   ```bash
   git clone https://github.com/<your-username>/k8s-agent-new.git
   cd k8s-agent-new
   ```
3. **Create a branch** for your feature or fix:
   ```bash
   git checkout -b feature/my-feature
   ```
4. **Make your changes** following the project conventions
5. **Test** your changes against a Kubernetes cluster
6. **Commit** with clear, descriptive messages:
   ```bash
   git commit -m "add support for StatefulSet resources"
   ```
7. **Push** and open a **Pull Request**:
   ```bash
   git push origin feature/my-feature
   ```

### Conventions

- Follow standard Go project layout and idioms
- Keep the clean architecture separation (core vs infrastructure)
- Add DTOs for new resource types in `internal/core/dto/`
- Register new tools in `internal/infrastructure/mcp/tools/`
- Use Kubernetes impersonation — never bypass RBAC
- Write descriptive MCP tool names and descriptions

### Ideas for Contribution

- Add support for new Kubernetes resources (StatefulSet, CronJob, PersistentVolumeClaim, etc.)
- Improve error messages and validation
- Add unit and integration tests
- Create Helm charts or Kustomize manifests for deployment
- Write documentation and examples

## License

This project is open-source. See the [LICENSE](LICENSE) file for details.

