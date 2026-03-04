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
./linkmngr links list --page 1
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

Primary resource commands are singular (`link`, `brand`, `page`, `domain`). Plural forms (`links`, `brands`, `pages`, `domains`) are supported as aliases.

### version

```bash
./linkmngr version
```

### auth

Set token:

```bash
./linkmngr auth login <token>
```

Set base URL:

```bash
./linkmngr auth set-base-url <url>
```

Get authenticated user:

```bash
./linkmngr auth status
./linkmngr auth status -o table
```

Revoke token:

```bash
./linkmngr auth logout
```

Aliases:
- `auth login` also supports `auth set-token`
- `auth status` also supports `auth whoami`
- `auth logout` also supports `auth revoke`

### links

List links:

```bash
./linkmngr link list [--page <n>] [--brand-id <id>] [--domain <domain>]
```

Options:
- `--page`, `-p` (default `1`)
- `--brand-id`
- `--domain`

Examples:

```bash
./linkmngr link list
./linkmngr link list --page 2 --brand-id 12
./linkmngr link list --domain linkmn.gr -o table
```

Get one link:

```bash
./linkmngr link get <link-id>
```

Alias:
- `link view`

Create link:

```bash
./linkmngr link create <destination> [--domain <domain>] [--slug <slug>] [--brand-id <id>]
```

Options:
- `--domain`
- `--slug`
- `--brand-id`

Examples:

```bash
./linkmngr link create https://example.com
./linkmngr link create https://example.com --domain linkmn.gr
./linkmngr link create https://example.com --domain linkmn.gr --slug spring-sale --brand-id 12
```

Get link stats:

```bash
./linkmngr link stats <link-id> --start <ISO8601> --end <ISO8601> [--time-unit <unit>] [--group-by <group>]
```

Required:
- `--start`
- `--end`

Optional:
- `--time-unit` (default `day`): `hour`, `day`, `week`, `month`, `year`
- `--group-by`: `device`, `device_type`, `country`, `browser`, `platform`, `referrer`

Examples:

```bash
./linkmngr link stats 123 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00
./linkmngr link stats 123 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --time-unit hour --group-by country -o table
```

### brands

List brands:

```bash
./linkmngr brand list [--page <n>]
```

Options:
- `--page`, `-p` (default `1`)

Examples:

```bash
./linkmngr brand list
./linkmngr brand list --page 2 -o table
```

Get one brand:

```bash
./linkmngr brand get <brand-id>
```

Check domain setup:

```bash
./linkmngr brand domain-check <brand-id> <domain>
```

Examples:

```bash
./linkmngr brand get 12
./linkmngr brand domain-check 12 linkmn.gr
```

### analytics

```bash
./linkmngr analytics --start <ISO8601> --end <ISO8601> [--time-unit <unit>] [--group-by <group>] [--brand-id <id>]
```

Required:
- `--start`
- `--end`

Optional:
- `--time-unit` (default `day`): `hour`, `day`, `week`, `month`, `year`
- `--group-by`: `device`, `device_type`, `country`, `browser`, `platform`, `referrer`
- `--brand-id`

Examples:

```bash
./linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00
./linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --time-unit day --group-by platform --brand-id 12 -o table
```

### domains

List available domains:

```bash
./linkmngr domain list
./linkmngr domain list -o table
```

### pages

List pages:

```bash
./linkmngr page list [--page <n>] [--brand-id <id>] [--domain <domain>] [--custom-domain-id <id>] [--slug <slug>] [--search <text>]
```

Options:
- `--page`, `-p` (default `1`)
- `--brand-id`
- `--domain`
- `--custom-domain-id`
- `--slug`
- `--search`

Examples:

```bash
./linkmngr page list
./linkmngr page list --brand-id 12 --search "product launch" -o table
./linkmngr page list --custom-domain-id 3 --slug my-bio
```

Get one page:

```bash
./linkmngr page get <page-id>
```

Alias:
- `page view`

Get page stats:

```bash
./linkmngr page stats <page-id> --start <ISO8601> --end <ISO8601> [--time-unit <unit>] [--group-by <group>]
```

Required:
- `--start`
- `--end`

Optional:
- `--time-unit` (default `day`): `hour`, `day`, `week`, `month`, `year`
- `--group-by`: `device`, `device_type`, `country`, `browser`, `platform`, `referrer`

Examples:

```bash
./linkmngr page stats 44 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00
./linkmngr page stats 44 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --group-by country -o table
```

Get recent page hits:

```bash
./linkmngr page hits <page-id>
./linkmngr page hits 44 -o table
```

Notes:
- `page create` is intentionally not implemented in this CLI.
- Use `api request` for advanced or undocumented page endpoints.

### api

Raw API access:

```bash
./linkmngr api request <method> <path> [--data <json> | --data-file <file> | --set key=value ...]
```

Options:
- `--data` raw JSON body
- `--data-file` JSON file path
- `--set` repeatable `key=value` fields

Rules:
- Use either `--data` or `--data-file`, not both.
- `--set` values are passed as strings.
- Path accepts both `/links` and `links`.

Examples:

```bash
./linkmngr api request GET /links
./linkmngr api request GET /pages
./linkmngr api request DELETE /links/123
./linkmngr api request POST /links --data '{"destination":"https://example.com","domain":"linkmn.gr"}'
./linkmngr api request POST /links --set destination=https://example.com --set domain=linkmn.gr --set slug=campaign-1
./linkmngr api request PATCH /brands/12 --data-file ./payload.json
./linkmngr api request POST /pages --data-file ./page-create.json
./linkmngr api request PATCH /pages/44 --set title=NewTitle
```

## Advanced Usage

Pipe JSON to `jq`:

```bash
./linkmngr links list | jq '.items[] | {id, link, clicks}'
./linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 | jq
```

Reusable date window:

```bash
START="2026-03-01T00:00:00+00:00"
END="2026-03-03T00:00:00+00:00"

./linkmngr links stats 123 --start "$START" --end "$END" --group-by country
./linkmngr pages stats 44 --start "$START" --end "$END" --group-by platform
```

## Troubleshooting

Missing token:

```text
missing API token; set with `linkmngr auth login <token>` or LINKMNGR_TOKEN
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
