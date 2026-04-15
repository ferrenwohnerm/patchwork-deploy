# patchwork-deploy

A lightweight CLI for managing incremental deployment configs across multiple environments without a full CD platform.

---

## Installation

```bash
go install github.com/yourorg/patchwork-deploy@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/patchwork-deploy.git
cd patchwork-deploy
go build -o patchwork-deploy .
```

---

## Usage

Define your environment configs in a `patchwork.yaml` file, then apply incremental patches per environment:

```bash
# Apply a config patch to the staging environment
patchwork-deploy apply --env staging --patch ./patches/v1.2.0.yaml

# Preview changes without applying
patchwork-deploy diff --env production --patch ./patches/v1.2.0.yaml

# List all applied patches for an environment
patchwork-deploy history --env production
```

**Example `patchwork.yaml`:**

```yaml
environments:
  - name: staging
    base: ./configs/base.yaml
  - name: production
    base: ./configs/base.yaml

patches_dir: ./patches
```

Patches are applied incrementally and tracked locally, giving you a lightweight audit trail without requiring a full CD pipeline.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)