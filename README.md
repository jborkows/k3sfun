# k3sfun

Small home lab configuration running on k3s.

## Repository Structure

Subprojects live in separate branches, allowing a single GitHub runner to serve multiple projects:

| Branch | Description |
|--------|-------------|
| `master` | K3s infrastructure and deployment configurations |
| `shoppinglist` | Go-based shopping list web application |

## Working with Subprojects

Subprojects are managed as git worktrees:

```bash
# List worktrees
git worktree list

# Add a worktree for a subproject
git worktree add .worktrees/<name> <branch>
```
