# linkmngr-cli

Go CLI for the LinkMngr API (built with Cobra).

## Project Status

- Active CLI for LinkMngr API operations.
- Suitable for local usage, scripting, and CI workflows.

## Requirements

- Go `1.22+` (if building from source)
- A valid LinkMngr API token

## Install

Build from source:

```bash
git clone <your-repo-url>
cd linkmngr-cli
go build -o linkmngr ./cmd/linkmngr
```

Or install with Go:

```bash
go install github.com/usama/linkmngr-cli/cmd/linkmngr@latest
```

## Builds and Releases

### Local Cross-Platform Builds

Use the included script to build release binaries for:
- `darwin/amd64`
- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`
- `windows/amd64`

Run:

```bash
./scripts/build-all.sh
```

Optional version override:

```bash
VERSION=v1.0.0 ./scripts/build-all.sh
```

Artifacts are written to `dist/` with per-file SHA256 checksum files.

### GitHub Releases (Automated)

This repository includes a release workflow at:
- `.github/workflows/release.yml`

Behavior:
- Trigger: push a tag matching `v*` (example: `v1.0.0`)
- Runs `go test ./...`
- Builds all platform binaries via `./scripts/build-all.sh`
- Uploads `dist/*` artifacts to the GitHub Release
- Generates release notes automatically

Release steps:

```bash
git tag v1.0.0
git push origin v1.0.0
```

After the workflow finishes, binaries will be attached to the GitHub Release for that tag.

## Quick Start

Set token:

```bash
./linkmngr auth login <your-token>
```

Optional: set API base URL (default is `https://api.linkmngr.com/v1`):

```bash
./linkmngr auth set-base-url https://api.linkmngr.com/v1
```

Run a command:

```bash
./linkmngr link list --page 1
```

## Configuration

- Config file path: `~/.linkmngr/config.json`
- Default base URL: `https://api.linkmngr.com/v1`
- Environment overrides:

```bash
export LINKMNGR_TOKEN=...
export LINKMNGR_BASE_URL=https://api.linkmngr.com/v1
```

Precedence:
- `LINKMNGR_TOKEN` overrides token from config file.
- `LINKMNGR_BASE_URL` overrides base URL from config file.

## Global Flags

All commands support:
- `-o, --output` with `table` (default) or `json`

Examples:

```bash
./linkmngr auth status --output json
./linkmngr link list -o table
```

## Command Reference

Primary resource commands are singular and production-default:
- `link` (alias: `links`)
- `brand` (alias: `brands`)
- `domain` (alias: `domains`)
- `page` (alias: `pages`)

Top-level commands:
- `version`
- `auth`
- `link` (`links` alias)
- `brand` (`brands` alias)
- `analytics`
- `domain` (`domains` alias)
- `page` (`pages` alias)
- `api`

### version

```bash
./linkmngr version
```

### auth

Subcommands:
- `login <token>` (alias: `set-token`)
- `set-base-url <url>`
- `status` (alias: `whoami`)
- `logout` (alias: `revoke`)

Examples:

```bash
./linkmngr auth login <token>
./linkmngr auth set-token <token>
./linkmngr auth set-base-url https://api.linkmngr.com/v1
./linkmngr auth status
./linkmngr auth whoami
./linkmngr auth logout
./linkmngr auth revoke
```

### link

Subcommands:
- `list` (alias: `ls`)
- `get <link-id>` (alias: `view`)
- `create <destination>`
- `stats <link-id>`

`link list` usage:

```bash
./linkmngr link list [--page <n>] [--brand-id <id>] [--domain <domain>]
```

Flags:
- `-p, --page` (default `1`)
- `--brand-id`
- `--domain`

`link get` usage:

```bash
./linkmngr link get <link-id>
./linkmngr link view <link-id>
```

`link create` usage:

```bash
./linkmngr link create <destination> [--domain <domain>] [--slug <slug>] [--brand-id <id>]
```

Flags:
- `--domain`
- `--slug`
- `--brand-id`

`link stats` usage:

```bash
./linkmngr link stats <link-id> --start <ISO8601> --end <ISO8601> [--time-unit <unit>] [--group-by <group>]
```

Required flags:
- `--start`
- `--end`

Optional flags:
- `--time-unit` (default `day`): `hour`, `day`, `week`, `month`, `year`
- `--group-by`: `device`, `device_type`, `country`, `browser`, `platform`, `referrer`

Examples:

```bash
./linkmngr link list
./linkmngr link ls --page 2 --brand-id 12
./linkmngr link create https://example.com --domain linkmn.gr --slug spring-sale --brand-id 12
./linkmngr link stats 123 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --time-unit hour --group-by country -o table
```

### brand

Subcommands:
- `list` (alias: `ls`)
- `get <brand-id>` (alias: `view`)
- `domain-check <brand-id> <domain>` (alias: `check-domain`)

`brand list` usage:

```bash
./linkmngr brand list [--page <n>]
```

Flags:
- `-p, --page` (default `1`)

`brand get` usage:

```bash
./linkmngr brand get <brand-id>
./linkmngr brand view <brand-id>
```

`brand domain-check` usage:

```bash
./linkmngr brand domain-check <brand-id> <domain>
./linkmngr brand check-domain <brand-id> <domain>
```

Examples:

```bash
./linkmngr brand list
./linkmngr brand get 12
./linkmngr brand domain-check 12 linkmn.gr
```

### analytics

Usage:

```bash
./linkmngr analytics --start <ISO8601> --end <ISO8601> [--time-unit <unit>] [--group-by <group>] [--brand-id <id>]
```

Required flags:
- `--start`
- `--end`

Optional flags:
- `--brand-id`
- `--time-unit` (default `day`): `hour`, `day`, `week`, `month`, `year`
- `--group-by`: `device`, `device_type`, `country`, `browser`, `platform`, `referrer`

Examples:

```bash
./linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00
./linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --time-unit day --group-by platform --brand-id 12 -o table
```

### domain

Subcommands:
- `list` (alias: `ls`)

Usage:

```bash
./linkmngr domain list
./linkmngr domain ls -o table
```

### page

Subcommands:
- `list` (alias: `ls`)
- `get <page-id>` (alias: `view`)
- `stats <page-id>`
- `hits <page-id>`

`page list` usage:

```bash
./linkmngr page list [--page <n>] [--brand-id <id>] [--domain <domain>] [--custom-domain-id <id>] [--slug <slug>] [--search <text>]
```

Flags:
- `-p, --page` (default `1`)
- `--brand-id`
- `--domain`
- `--custom-domain-id`
- `--slug`
- `--search`

`page get` usage:

```bash
./linkmngr page get <page-id>
./linkmngr page view <page-id>
```

`page stats` usage:

```bash
./linkmngr page stats <page-id> --start <ISO8601> --end <ISO8601> [--time-unit <unit>] [--group-by <group>]
```

Required flags:
- `--start`
- `--end`

Optional flags:
- `--time-unit` (default `day`): `hour`, `day`, `week`, `month`, `year`
- `--group-by`: `device`, `device_type`, `country`, `browser`, `platform`, `referrer`

`page hits` usage:

```bash
./linkmngr page hits <page-id>
```

Examples:

```bash
./linkmngr page list --brand-id 12 --search "product launch"
./linkmngr page get 44
./linkmngr page stats 44 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --group-by country -o table
./linkmngr page hits 44
```

### api

Subcommands:
- `request <method> <path>`

`api request` usage:

```bash
./linkmngr api request <method> <path> [--data <json> | --data-file <file> | --set key=value ...]
```

Flags:
- `--data` raw JSON request body
- `--data-file` path to JSON request body file
- `--set` repeatable `key=value` body fields

Rules:
- Use either `--data` or `--data-file`, not both.
- `--set` values are sent as strings.
- `<path>` accepts `/links` and `links`.

Examples:

```bash
./linkmngr api request GET /links
./linkmngr api request DELETE /links/123
./linkmngr api request POST /links --data '{"destination":"https://example.com","domain":"linkmn.gr"}'
./linkmngr api request POST /links --set destination=https://example.com --set domain=linkmn.gr --set slug=campaign-1
./linkmngr api request PATCH /brands/12 --data-file ./payload.json
```

## Advanced Usage

Pipe JSON to `jq`:

```bash
./linkmngr link list | jq '.items[] | {id, link, clicks}'
./linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 | jq
```

Reusable date window:

```bash
START="2026-03-01T00:00:00+00:00"
END="2026-03-03T00:00:00+00:00"

./linkmngr link stats 123 --start "$START" --end "$END" --group-by country
./linkmngr page stats 44 --start "$START" --end "$END" --group-by platform
```

## Troubleshooting

Missing token:

```text
missing API token; set LINKMNGR_TOKEN or run `linkmngr auth login <token>`
```

Common fixes:
- Set token using `auth login`.
- Export `LINKMNGR_TOKEN`.
- Ensure ID arguments are positive integers.
- For raw API requests, use valid JSON and correct `--set key=value` format.

## Contributing

Contributions are welcome via pull requests.

Suggested local checks before opening a PR:

```bash
go test ./...
go vet ./...
```

Please include:
- clear description of behavior changes
- updated docs/examples for CLI changes
- tests for new behavior where possible

## Security

- Never commit API tokens or credentials.
- Report security issues privately to maintainers instead of opening a public issue.

## Support

Open an issue with:
- exact command used
- expected behavior
- actual output
- environment details (OS, Go version, CLI version)

## License

This project is licensed under the MIT License.
See [LICENSE](./LICENSE).
