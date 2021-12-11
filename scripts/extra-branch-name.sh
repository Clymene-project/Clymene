#!/bin/bash

echo "Note: GITHUB_REF=${GITHUB_REF}"

if [[ "${GITHUB_REF}" == "refs/heads"* ]]; then
    echo "Note: This is a push to a local branch -> using branch name"
    BRANCH=${GITHUB_REF#refs/heads/}
    BRANCH_SLUG=$(echo $BRANCH | iconv -t ascii//TRANSLIT | sed -r s/[^a-zA-Z0-9]+/-/g | sed -r s/^-+\|-+$//g | tr A-Z a-z)
else
    if [[ "${GITHUB_REF}" == "refs/pull/"* ]]; then
        # usually the format for PRs is: refs/pull/1234/merge
        echo "Note: This is a Pull Request -> using PR ID"
        tmp=${GITHUB_REF#refs/pull/}
        # remove the last "/merge"
        # Branch name is basically the PR id
        BRANCH=PR-${tmp%/merge}
        # And Slug is "PR-${PRID}"
        BRANCH_SLUG=${BRANCH}
    else
        echo "::error This is neither a push, nor a PR, probably something else... Exiting"
        exit 1
    fi
fi
GIT_SHA="$(git rev-parse --short HEAD)"

# print GIT_SHA, BRANCH and BRANCH_SLUG (make sure they are also set in needs.prepare_ci_run.outputs !!!)
echo "##[set-output name=BRANCH;]$(echo ${BRANCH})"
echo "##[set-output name=BRANCH_SLUG;]$(echo ${BRANCH_SLUG})"
echo "##[set-output name=GIT_SHA;]$(echo ${GIT_SHA})"