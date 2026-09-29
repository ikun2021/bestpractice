# Using zero-skills with Codex

Codex supports the [Agent Skills specification](https://agentskills.io/) and can load this repository directly as a skill. This gives Codex progressive access to the go-zero guidance in `SKILL.md` and the referenced pattern files without copying the guidance into `AGENTS.md`.

## Installation

### Project-level installation (recommended)

Install the skill for everyone working in a repository:

```bash
mkdir -p .agents/skills
git clone https://github.com/zeromicro/zero-skills.git .agents/skills/zero-skills
```

Codex scans `.agents/skills` from the current working directory up to the repository root. You can also add the repository as a Git submodule:

```bash
git submodule add https://github.com/zeromicro/zero-skills.git .agents/skills/zero-skills
```

### Personal installation

Install the skill once for use across all repositories:

```bash
mkdir -p ~/.agents/skills
git clone https://github.com/zeromicro/zero-skills.git ~/.agents/skills/zero-skills
```

Codex detects skill changes automatically. Restart Codex if the skill does not appear after installation.

## Verify and invoke the skill

In Codex CLI or the IDE extension:

1. Run `/skills` and confirm that `zero-skills` is listed.
2. Mention `$zero-skills` to invoke it explicitly.
3. Or describe a go-zero task normally; Codex can select the skill implicitly when the task matches its description.

Example prompts:

```text
$zero-skills Create a user management REST API with CRUD operations.
```

```text
Add a cached MySQL model to this go-zero service and follow the repository's go-zero skill.
```

```text
Diagnose this go-zero timeout and apply the documented resilience patterns: ...
```

## Skills and AGENTS.md serve different purposes

Use the installed skill for reusable go-zero knowledge. Use `AGENTS.md` for short, repository-specific instructions such as build commands, package boundaries, naming conventions, or required checks.

For example:

```markdown
# Repository instructions

- Use the zero-skills skill for go-zero implementation patterns.
- Run `go test ./...` after changing Go code.
- Keep generated `goctl` files in `internal/` unchanged unless the task requires regeneration.
```

Codex reads `AGENTS.md` automatically, but you do not need to duplicate this skill's pattern documentation there.

## Updating

For a cloned installation:

```bash
git -C .agents/skills/zero-skills pull --ff-only
```

For a submodule installation:

```bash
git submodule update --remote .agents/skills/zero-skills
```

## Troubleshooting

### The skill is not listed

1. Confirm that `.agents/skills/zero-skills/SKILL.md` exists (or `~/.agents/skills/zero-skills/SKILL.md` for a personal installation).
2. Run Codex inside the repository containing `.agents/skills`.
3. Run `/skills` again or restart Codex.
4. Check that the skill directory contains the repository's `SKILL.md`, not an extra nested `zero-skills/zero-skills` directory.

### The skill is not selected automatically

Invoke it explicitly with `$zero-skills`, then include the go-zero task in the same prompt. Automatic selection depends on how closely the request matches the skill description.

## Official documentation

- [Build skills in Codex](https://developers.openai.com/codex/skills/)
- [Agent Skills specification](https://agentskills.io/)
