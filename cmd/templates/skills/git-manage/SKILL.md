---
name: git-manage
description: Git repository operations (status, commit, branch, log, diff, stash, remote).
tags: [git, dev, cross-platform]
---
# Git Management

Cross-platform Git operations. All commands use standard `git` CLI.

## Status & Info

Working tree status:
```
exec: git -C /path/to/repo status
```

Short status:
```
exec: git -C /path/to/repo status -sb
```

Current branch:
```
exec: git -C /path/to/repo branch --show-current
```

Remote URL:
```
exec: git -C /path/to/repo remote -v
```

## Log & History

Recent commits (oneline):
```
exec: git -C /path/to/repo log --oneline -20
```

Detailed log with graph:
```
exec: git -C /path/to/repo log --oneline --graph --decorate -30
```

Commits by author:
```
exec: git -C /path/to/repo log --author="NAME" --oneline -20
```

Commits in date range:
```
exec: git -C /path/to/repo log --after="2026-01-01" --before="2026-02-01" --oneline
```

File change history:
```
exec: git -C /path/to/repo log --oneline --follow -- path/to/file
```

## Diff

Unstaged changes:
```
exec: git -C /path/to/repo diff
```

Staged changes:
```
exec: git -C /path/to/repo diff --cached
```

Diff between branches:
```
exec: git -C /path/to/repo diff main..feature-branch --stat
```

Diff specific file:
```
exec: git -C /path/to/repo diff -- path/to/file
```

## Stage & Commit

Stage files:
```
exec: git -C /path/to/repo add path/to/file1 path/to/file2
```

Stage all changes:
```
exec: git -C /path/to/repo add -A
```

Commit:
```
exec: git -C /path/to/repo commit -m "COMMIT_MESSAGE"
```

Amend last commit message:
```
exec: git -C /path/to/repo commit --amend -m "NEW_MESSAGE"
```

## Branch

List branches:
```
exec: git -C /path/to/repo branch -a
```

Create and switch:
```
exec: git -C /path/to/repo checkout -b NEW_BRANCH
```

Switch branch:
```
exec: git -C /path/to/repo checkout BRANCH_NAME
```

Delete branch:
```
exec: git -C /path/to/repo branch -d BRANCH_NAME
```

## Remote & Sync

Fetch:
```
exec: git -C /path/to/repo fetch --all
```

Pull:
```
exec: git -C /path/to/repo pull origin BRANCH
```

Push:
```
exec: git -C /path/to/repo push origin BRANCH
```

Push new branch:
```
exec: git -C /path/to/repo push -u origin BRANCH
```

## Stash

Save stash:
```
exec: git -C /path/to/repo stash push -m "DESCRIPTION"
```

List stashes:
```
exec: git -C /path/to/repo stash list
```

Apply latest stash:
```
exec: git -C /path/to/repo stash pop
```

## Tags

List tags:
```
exec: git -C /path/to/repo tag -l
```

Create tag:
```
exec: git -C /path/to/repo tag -a v1.0.0 -m "Release v1.0.0"
```

Push tag:
```
exec: git -C /path/to/repo push origin v1.0.0
```

## Search

Search commit messages:
```
exec: git -C /path/to/repo log --grep="KEYWORD" --oneline
```

Search code changes (pickaxe):
```
exec: git -C /path/to/repo log -S "function_name" --oneline
```

Blame (who changed each line):
```
exec: git -C /path/to/repo blame path/to/file
```

## GitHub CLI (if `gh` is installed)

Create PR:
```
exec: gh pr create --title "TITLE" --body "DESCRIPTION" --repo owner/repo
```

List PRs:
```
exec: gh pr list --repo owner/repo
```

List issues:
```
exec: gh issue list --repo owner/repo
```

## Notes

- Replace `/path/to/repo` with actual repo path, or omit `-C` if already in the repo.
- Works on macOS, Linux, and Windows (Git Bash / WSL).
- `gh` CLI is optional but useful for GitHub-specific operations.
