#!/usr/bin/env python3
"""Generate private deployment credentials without replacing existing files."""

import os
from pathlib import Path
import secrets
import sys


def main():
    root = Path(__file__).resolve().parent
    env_path = root / ".env"
    config_path = root / "cpa/config.yaml"
    for path in (env_path, config_path):
        if path.exists() or path.is_symlink():
            print(f"Refusing to overwrite existing configuration: {path.name}", file=sys.stderr)
            return 1

    values = {
        "CPA_MANAGER_ADMIN_KEY": "cpamp_" + secrets.token_hex(32),
        "CPA_MANAGEMENT_KEY": secrets.token_hex(32),
        "CPA_API_KEY": "sk-" + secrets.token_hex(32),
        "NEW_API_SESSION_SECRET": secrets.token_hex(32),
    }
    env = (root / ".env.example").read_text(encoding="utf-8")
    for key, value in values.items():
        marker = key + "=\n"
        if env.count(marker) != 1:
            raise ValueError(f"Expected exactly one empty setting: {key}")
        env = env.replace(marker, key + "=" + value + "\n")

    config = (root / "cpa/config.example.yaml").read_text(encoding="utf-8")
    replacements = {
        '  secret-key: ""': f'  secret-key: "{values["CPA_MANAGEMENT_KEY"]}"',
        "api-keys: []": f'api-keys:\n  - "{values["CPA_API_KEY"]}"',
    }
    for marker, value in replacements.items():
        if config.count(marker) != 1:
            raise ValueError("Unexpected CPA configuration template")
        config = config.replace(marker, value)

    created = []
    try:
        for path, content in ((env_path, env), (config_path, config)):
            descriptor = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
            created.append(path)
            with os.fdopen(descriptor, "w", encoding="utf-8") as output:
                output.write(content)
    except OSError:
        for path in created:
            path.unlink()
        raise

    print("Created .env and cpa/config.yaml with private random credentials.")
    print("Start services with: docker compose up -d --wait --wait-timeout 180")
    return 0


if __name__ == "__main__":
    sys.exit(main())
