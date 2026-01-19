#!/usr/bin/env sh
# Copyright (C) 2024-2026 Bryce Thuilot
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the FSF, either version 3 of the License, or (at your option) any later version.
# See the LICENSE file in the root of this repository for full license text or
# visit: <https://www.gnu.org/licenses/gpl-3.0.html>.


echo "beginning git-lost-and-found search" 1>&2

if ! git status -s 1> /dev/null;
then
    echo "ERROR: current directory is not apart of git repository" 1>&2
    exit 1
fi

GIT_DIR="$(git rev-parse --git-dir || (echo "ERROR: unable to retrieve git dir"; exit 1))"

if ! git fsck --no-progress --lost-found 1>/dev/null ;
then
    echo "ERROR: git returned error when looking for dangling commits" 1>&2
    exit 1
fi


DANGLING_REFS="${GIT_DIR}/refs/dangling"
mkdir -p "$DANGLING_REFS"

for found_commit in "${GIT_DIR}/lost-found/commit/"*; do
    if [ -f "$found_commit" ]; then
	commit="$(basename "$found_commit")"
	echo "found dangling commit '$commit'" 1>&2
	cp "$found_commit" "$DANGLING_REFS/$commit"
    fi
done

rm -rf "${GIT_DIR}/lost-found/commit/"* || (echo "unable to cleanup lost-found directory"; exit 1)

echo "complete" 2>&1
