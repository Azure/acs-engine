#!/usr/bin/env bash

exit_code=0

echo "==> Running dep check <=="

dep check || exit_code=1

if [ $exit_code -ne 0 ]; then
  echo "The dependency state is out of sync. Please run dep ensure."
else
  echo "dep check passed."
fi

exit $exit_code