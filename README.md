# gh-repo-sync

gh-repo-sync syncs all repositories that you forked.
By setting configuration file, you can manage the repositories you want to sync.

## Installation

```bash
gh extensions install sivchari/gh-repo-sync
```

## Usage

gh-repo-sync requires a configuration file to sync repositories.
The configuration file is a YAML file that contains the repositories you want to sync.

The default name is gh-repo-sync.yaml.

This file is expected to be at top level of that you want to sync all repositories.

(e.g.)

The work tree is like this:

```shell
.
├── gh-repo-sync.yaml
├── repo1
├── repo2
├── sub
│   └── repo3
```

Then, when you want to sync repo1, repo2, and repo3, you should write the configuration file like this:

```yaml
repositories:
  - repo1
  - repo2
  - sub/repo3
```

After setting the configuration file, you can sync repositories by running the following command:

```bash
gh repo-sync
```
