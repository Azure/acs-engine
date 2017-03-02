#!/usr/bin/env bash

set -x

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
ROOT="${DIR}/.."

# This script is meant to be used to help out when performing a rebase.
# Particularly, when two people have editted the templates, the generated file
# needs to be updated based on the merge of your parts edits. This script attempts
# to generate, add the templates.go file, fix any of the expected testdata output,
# and then attempts to push the rebase forward.

# If anything fails, we bail out and leave the user to repair.

# This will fail to push the rebase forward if there are other conflicts to be resolved.
# It's recommended you re-run this script again after resolving them to ensure that the
# generated file and expected test output reflects your merge resolutions.

# Go ahead and generate and make sure it's git added
cd "${ROOT}"
go generate -v ./...
git add "${ROOT}/pkg/acsengine/templates.go"

# Fixup all the testdata (this requires manual review in the PR!)
cd "${ROOT}"
make ci || true

cd "${ROOT}/pkg/acsengine/testdata"
for f in **/*.json ; do
	mv $f{.err,}
	git add $f
done

cd "${ROOT}"
make ci

# if we're in a rebase, try to advance
if [[ -e "${ROOT}/.git/rebase-merge/done" ]]; then
	# try to advance the rebase
	git rebase --continue
	# if we successfully advanced...
	if [[ "$?" != "0" ]]; then
		# and are still in a rebase...
		if [[ -e "${ROOT}/.git/rebase-merge/done" ]]; then
			# let's do it again
			exec "${BASH_SOURCE}"
		fi
	fi
fi
