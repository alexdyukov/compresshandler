git config --global user.name 'update_deps robot'
git config --global user.email 'noreply@example.com'
git remote set-url origin https://x-access-token:${GH_TOKEN}@github.com/${REPO_NAME}

git checkout -b ${BRANCH_NAME}

git commit -am "Fix: $(date +%F) update dependencies" || exit 0

git push --set-upstream origin ${BRANCH_NAME} -f || exit 0

gh pr create -a ${REPO_OWNER} -b "$(go test ./...)" --fill || exit 0
