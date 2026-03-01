#!/bin/bash

# Version bump script for TriageProf
# Usage: ./scripts/version-bump.sh [major|minor|patch]

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 [major|minor|patch]"
    exit 1
fi

TYPE=$1
CURRENT_VERSION=$(cat VERSION)

# Parse version
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Bump version
case $TYPE in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Invalid version type: $TYPE"
        exit 1
        ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"

# Update VERSION file
echo "$NEW_VERSION" > VERSION

# Create git tag
git tag "v$NEW_VERSION"

echo "Version bumped from $CURRENT_VERSION to $NEW_VERSION"
echo "Created git tag: v$NEW_VERSION"
echo "Push the tag with: git push origin v$NEW_VERSION"
