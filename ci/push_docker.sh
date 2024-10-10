#!/bin/bash

DOCKER_REPO="semitechnologies/weaviate"

function release() {
  # for multi-platform build
  docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
  docker buildx create --use

  tag_latest="${DOCKER_REPO}:latest"
  tag_exact=
  tag_preview=

  git_hash=$(echo "$GITHUB_SHA" | cut -c1-7)

  echo "GITHUB_REF_TYPE: $GITHUB_REF_TYPE"
  echo "GITHUB_REF_NAME: $GITHUB_REF_NAME"
  echo "GITHUB_HEAD_REF: $GITHUB_HEAD_REF"
  echo "GITHUB_BASE_REF: $GITHUB_BASE_REF"
  echo "PR_TITLE: $PR_TITLE"

  weaviate_version="$(jq -r '.info.version' < openapi-specs/schema.json)"
  if [  "$GITHUB_REF_TYPE" == "tag" ]; then
        if [ "$GITHUB_REF_NAME" != "v$weaviate_version" ]; then
            echo "The release tag ($GITHUB_REF_NAME) and Weaviate version (v$weaviate_version) are not equal! Can't release."
            return 1
        fi
        tag_exact="${DOCKER_REPO}:${weaviate_version}"
  else
    pr_title="$(echo -n "$PR_TITLE" | tr '[:upper:]' '[:lower:]' | tr -c -s '[:alnum:]' '-' | sed 's/-$//g')"
    if [ "$pr_title" == "" ]; then
      prefix="$(echo -n $GITHUB_BASE_REF | sed 's/\//-/g')"
      tag_preview="${DOCKER_REPO}:${prefix}-${git_hash}"
      weaviate_version="${prefix}-${git_hash}"
    else
      tag_preview="${DOCKER_REPO}:preview-${pr_title}-${git_hash}"
      weaviate_version="preview-${pr_title}-${git_hash}"
    fi
  fi

  args=("--build-arg=GITHASH=$git_hash" "--build-arg=DOCKER_IMAGE_TAG=$weaviate_version" "--platform=linux/amd64,linux/arm64" "--target=weaviate" "--push")
  if [ -n "$tag_exact" ]; then
    # exact tag on master
    args+=("-t=$tag_exact")
    args+=("-t=$tag_latest")
  fi
  if [ -n "$tag_preview" ]; then
    # preview tag on PR builds
    args+=("-t=$tag_preview")
  fi

  docker buildx build "${args[@]}" .

  if [ -n "$tag_preview" ]; then
    echo "PREVIEW_TAG=$tag_preview" >> "$GITHUB_OUTPUT"
  elif [ -n "$tag_exact" ]; then
    echo "PREVIEW_TAG=$tag_exact" >> "$GITHUB_OUTPUT"
  fi
}

release
