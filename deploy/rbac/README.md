# Templates RBAC (multi-tenant)

Use estes arquivos como base para dar acesso a um **usuário específico** (identidade do JWT) a recursos em um namespace. A API k8s-agent não conhece tenants; o isolamento é feito pelo RBAC do cluster.

**Substitua os placeholders em cada arquivo:**

| Placeholder   | Descrição |
|---------------|-----------|
| `NAMESPACE`   | Namespace onde o usuário terá permissão (ex.: `meu-projeto`, `microcontainers`) |
| `USER_NAME`   | Nome do usuário que virá no JWT (claim `sub`). Deve ser o mesmo na RoleBinding e no ClusterRole/ClusterRoleBinding |
| `ROLE_NAME`   | Nome da Role e do RoleBinding (ex.: `tenant-pods-access`) |

Para a ClusterRoleBinding de impersonation, o **subject** deve ser a identidade que **roda a API** (ServiceAccount ou User do kubeconfig). Veja `clusterrolebinding-impersonate.yaml`.

**Ordem sugerida:** aplicar Role e RoleBinding (namespaced), depois ClusterRole e ClusterRoleBinding (cluster-scoped).

**Múltiplos tenants:** repita o conjunto (Role + RoleBinding + ClusterRole impersonate + ClusterRoleBinding) para cada `USER_NAME`/namespace. No ClusterRole de impersonation você pode incluir vários `resourceNames` para vários usuários e um único ClusterRoleBinding para a API.
