# Its looks like semver, but its not, because we cannot release patch on prev major.minor release
MAJOR_VERSION=1
MAJOR_LAST_COMMIT_HASH="873c3256312376bbc712fd180f50f1a0eff12482"

MINOR_LAST_COMMIT_HASH=$(git rev-list --invert-grep -i --grep="fix" ${MAJOR_LAST_COMMIT_HASH}..HEAD -n 1)
MINOR_VERSION=$(git rev-list --invert-grep -i --grep="fix" ${MAJOR_LAST_COMMIT_HASH}..HEAD --count)

PATCH_VERSION=$(git rev-list ${MINOR_LAST_COMMIT_HASH}..HEAD --count)

APP_VERSION=v${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}

# get current commit hash for tag
COMMIT_HASH=$(git rev-parse HEAD)

# POST a new ref to repo via Github API
curl -s -X POST https://api.github.com/repos/${GITHUB_REPOSITORY}/git/refs \
	-H "Authorization: token ${GITHUB_TOKEN}" \
	-d @- << EOF
{
  "ref": "refs/tags/${APP_VERSION}",
  "sha": "${COMMIT_HASH}"
}
EOF
