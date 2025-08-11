# gh-pr-stats

📊 Export GitHub Pull Request statistics via the [GitHub CLI](https://github.com/cli/cli)

A `gh` CLI extension that lets you export PR stats aggregated by labels with various format (CSV, JSON, TSV).

![Screencast](https://github.com/user-attachments/assets/6f606cea-6284-4674-af32-bb1b718e261d)


## Prerequisite

- [GitHub CLI](https://github.com/cli/cli)

## Installation

```bash
$ gh extension install shufo/gh-pr-stats
```

## Output Example

Stats from [GitHub CLI](https://github.com/cli/cli) (`cli/cli`)

```bash
$ gh pr-stats cli/cli -f csv
Label,Open,Closed,Total,Open %,Average Time to close (days),Median Time to close (days)
*unlabeled*,7,1916,1923,0.36,9,1
external,20,959,979,2.04,22,3
dependencies,0,71,71,0.00,4,2
go,0,47,47,0.00,3,1
enhancement,0,26,26,0.00,96,48
github_actions,0,24,24,0.00,5,3
actions,0,17,17,0.00,3,2
docs,0,14,14,0.00,2,1
needs-user-input,0,13,13,0.00,41,6
discuss,1,9,10,10.00,117,43
blocked,0,8,8,0.00,141,54
core,0,5,5,0.00,20,13
codespaces,0,4,4,0.00,80,7
invalid,0,4,4,0.00,0,0
tech-debt,0,2,2,0.00,14,14
gh-pr,0,1,1,0.00,12,12
platform,0,1,1,0.00,50,50
packaging,0,1,1,0.00,7,7
gh-attestation,0,1,1,0.00,3,3
gh-api,0,1,1,0.00,43,43
needs-investigation,0,1,1,0.00,5,5
Total,27,3010,3037,0.89%,13,1
```

## Usage

- Basic usage

```bash
gh pr-stats
```

- Specific repository

```bash
gh pr-stats owner/repo
```

- Change output format. (default: table. Supports `json`, `csv` and `tsv`)

```bash
gh pr-stats --format json
gh pr-stats owner/repo --format csv
gh pr-stats owner/repo --format tsv
```

- Persist aggregated results to file

```bash
gh pr-stats -s stats.json
```

- Persist raw source data to file

```bash
gh pr-stats -o prs.json
```

- Verbose output

```bash
gh pr-stats --debug
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## Development

```bash
# build extension
$ make build
# or go-task build to watch source
# Install extension locally
$ gh extension install .
# Run 
$ gh pr-stats
```

## Testing

```bash
$ make test
# or go-task test
```

## LICENSE

MIT
