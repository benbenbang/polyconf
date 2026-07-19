# polyconf

Configuration schemas for CLI tools monorepo and published in every format
the tools accept.

Each tool has one canonical [JSON Schema](https://json-schema.org/) definition plus
annotated reference files in TOML and YAML. A config file points at the schema URL and the
editor gives validation, autocompletion, and inline docs for free.

## Schemas

| Tool | Prefix | Location | Status |
| --- | --- | --- | --- |
| [Consilium](https://github.com/benbenbang/consilium) — AI commit-message generator | `csl` | `src/schemas/consilium/` | Available |
| ccs-picker |  | `src/schemas/ccs-picker/` | Planned |
| pcg | — | `src/schemas/pcg/` | Planned |

Every schema set ships three synced files:

| File | Purpose |
| --- | --- |
| `*-config.schema.json` | JSON Schema (draft-07) — **the source of truth** |
| `*_schema.toml` | Fully-commented TOML reference / starter config |
| `*_schema.yaml` | Fully-commented YAML reference / starter config |

The TOML and YAML files list every option with its default and a short description, all
commented out — copy one, uncomment what you need, and drop the rest.

## Usage

### Reference the schema from a config file

Add a schema directive to the top of the config file so the editor (VS Code, Zed,
JetBrains, Neovim, etc.) can validate it and offer completions.

TOML — via [Taplo](https://taplo.tamasfe.dev/):

```toml
#:schema https://raw.githubusercontent.com/benbenbang/polyconf/refs/tags/0.1.0/src/schemas/consilium/csl-config.schema.json

[commit]
require_scope = true
max_subject_length = 72
```

YAML — via the [YAML Language Server](https://github.com/redhat-developer/yaml-language-server):

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/benbenbang/polyconf/refs/tags/0.1.0/src/schemas/consilium/csl-config.schema.json

commit:
  require_scope: true
  max_subject_length: 72
```

> Pin to a release tag (e.g. `refs/tags/0.1.0`) for stability, or use `refs/heads/main` to
> track the latest.

### Start from a template

Copy the annotated reference for the tool and trim it down:

```sh
# Consilium, as TOML
curl -o .csl.toml \
  https://raw.githubusercontent.com/benbenbang/polyconf/main/src/schemas/consilium/csl_schema.toml
```

## Consilium schema at a glance

`csl-config.schema.json` describes Consilium's full configuration surface. Top-level
sections:

| Section | What it controls |
| --- | --- |
| `commit` | Commit rules — scope requirements, subject length, allowed types, scope modes (`allowed-list` / `regex` / `mappings`), validation patterns, diff exclusions |
| `llm` | Provider and model selection, credentials, `auth` for cloud providers (AWS / Azure / GCloud), temperature, token limits, prompt customization |
| `stream` | Diff-processing mode — `spinner` (single call) vs `streaming` (per-file review; 3–5× cost) and the line-count threshold |
| `ui` | Terminal theme and progress-spinner behavior |
| `context` | Optional repository-context gathering |
| `include` | Conditional inclusion of additional config fragments |

Supported LLM providers: `anthropic`, `openai`, `ollama`, `grok`, `gemini`,
`gemini-cli`, `claude-code`, `mistral`, `bedrock`.

The [JSON Schema](src/schemas/consilium/csl-config.schema.json) is authoritative; the
[TOML](src/schemas/consilium/csl_schema.toml) and
[YAML](src/schemas/consilium/csl_schema.yaml) references are the readable walkthrough with
defaults.

## Layout

```
polyconf/
├── src/
│   └── schemas/
│       └── consilium/
│           ├── csl-config.schema.json   # source of truth
│           ├── csl_schema.toml          # annotated TOML reference
│           └── csl_schema.yaml          # annotated YAML reference
└── scripts/                             # repo tooling helpers
```

When a schema changes, keep the three formats in sync — the JSON Schema leads, and the
TOML/YAML references mirror its options and defaults.
