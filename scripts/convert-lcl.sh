#!/usr/bin/env bash

# DON'T RUN. This script requires PythonLocalizerTool which is not published yet.
# TODO: make PythonLocalizerTool public
set -eo pipefail

langdirs=loc/*

convert_lcl_to_po() {
  for dir in $langdirs
  do
    loc_lang=`basename "$dir"`
    translation_lang=`echo $loc_lang | tr - _`
    publish/PythonLocalizerTool lcltopo $dir translations/$translation_lang/LC_MESSAGES/ translations/en_US/LC_MESSAGES/en-US/metadata acsengine ""
    msgfmt -c -v -o translations/$translation_lang/LC_MESSAGES/acsengine.mo translations/$translation_lang/LC_MESSAGES/acsengine.po
  done
}

convert_lcl_to_po