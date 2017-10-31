#!/bin/bash

python pkg/acsengine/Get-AzureConstants.py
git status | grep pkg/acsengine/azureconst.go
exit_code=$?
if [ $exit_code -gt "0" ]; then
  exit 0
else
  echo "File was modified, failing test"
  exit 1
fi 
