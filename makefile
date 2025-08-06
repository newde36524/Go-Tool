lint-fix: lint
	golangci-lint run --fix

.PHONY: amend
amend:
	git commit --amend --no-edit

# 同步tag
.PHONY: sync-tag
sync-tag:
	git tag -l | xargs git tag -d
	git pull

# 同步远程分支
.PHONY: fetch-prune
fetch-prune:
	git fetch --prune

.PHONY: force-with-lease
force-with-lease:
	git push origin $(git rev-parse --abbrev-ref HEAD) --force-with-lease

# 定义 rebase 目标，执行 rebase 操作 例子: make rebase dev-2.1
.PHONY: rebase
rebase:
	@echo "Starting rebase to target branch $(word 2,$(MAKECMDGOALS))"
	@$(SHELL) -c '\
		targetBranch=$(word 2,$(MAKECMDGOALS)); \
		if [ -z "$$targetBranch" ]; then \
			echo "The target branch cannot be empty"; \
			echo "example: \n\t make rebase dev-2.1"; \
			exit; \
		fi; \
		echo "target_branch:$$targetBranch"; \
		currentBranch=$$(git rev-parse --abbrev-ref HEAD); \
		echo "current_branch:$$currentBranch"; \
		git checkout "$$targetBranch" && git pull; \
		lastCommitHead=$$(git rev-parse HEAD); \
		git checkout "$$currentBranch"; \
		git rebase --onto "$$targetBranch" "$$lastCommitHead" "$$currentBranch" && git push origin $$currentBranch --force-with-lease && echo "Rebase completed successfully"; \
		'
.PHONY: %
%:
	@:

.PHONY: git_stash_push
git_stash_push:
	git stash push

.PHONY: git_stash_pop
git_stash_pop:
	git stash pop


