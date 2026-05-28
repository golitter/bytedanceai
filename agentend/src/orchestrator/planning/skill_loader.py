from __future__ import annotations

import re
from pathlib import Path

import yaml


def _parse_frontmatter(text: str) -> dict | None:
    """Parse YAML frontmatter (--- ... ---) from text. Returns dict or None."""
    match = re.match(r"^---\s*\n(.*?)\n---", text, re.DOTALL)
    if not match:
        return None
    try:
        return yaml.safe_load(match.group(1))
    except yaml.YAMLError:
        return None


def _strip_frontmatter(text: str) -> str:
    """Remove YAML frontmatter and return body text."""
    return re.sub(r"^---\s*\n.*?\n---\s*\n?", "", text, count=1, flags=re.DOTALL)


def discover_skills(builtin_dir: str | Path) -> list[dict]:
    """Scan builtin_dir subdirectories for SKILL.md files with valid frontmatter.

    Returns list of dicts with 'name' and 'description' keys.
    Skips directories without SKILL.md or without a 'name' field.
    """
    builtin = Path(builtin_dir)
    if not builtin.is_dir():
        return []

    skills: list[dict] = []
    for child in sorted(builtin.iterdir()):
        if not child.is_dir():
            continue
        skill_md = child / "SKILL.md"
        if not skill_md.is_file():
            continue
        fm = _parse_frontmatter(skill_md.read_text(encoding="utf-8"))
        if not fm or "name" not in fm:
            continue
        skills.append(
            {
                "name": fm["name"],
                "description": fm.get("description", ""),
            }
        )
    return skills


def load_skill_l2(skill_name: str, builtin_dir: str | Path) -> str:
    """Load full SKILL.md body (excluding frontmatter) for a given skill."""
    skill_md = Path(builtin_dir) / skill_name / "SKILL.md"
    if not skill_md.is_file():
        return ""
    return _strip_frontmatter(skill_md.read_text(encoding="utf-8")).strip()


def load_skill_resource(skill_name: str, resource_path: str, builtin_dir: str | Path) -> str:
    """Load a resource file from a skill's directory.

    resource_path must not contain '..' (path traversal check).
    Reads from builtin_dir/skill_name/resource_path.
    """
    if ".." in resource_path:
        return "Error: invalid resource path"

    full_path = Path(builtin_dir) / skill_name / resource_path
    try:
        return full_path.read_text(encoding="utf-8")
    except FileNotFoundError:
        return f"Error: resource file not found: {resource_path}"
    except Exception as e:
        return f"Error: {e}"
