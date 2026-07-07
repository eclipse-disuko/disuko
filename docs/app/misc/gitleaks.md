# install
install pre-commit https://pre-commit.com/#installation

# check installation
pre-commit --version

# configure pre-commit hook with gitleaks
Create a [.pre-commit-config.yaml](../../../.pre-commit-config.yaml) file in the root directory of this Git repository.

```
repos:
  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.30.0  # Specify the desired version of Gitleaks
    hooks:
      - id: gitleaks
```

# install hook scripts:
pre-commit install

