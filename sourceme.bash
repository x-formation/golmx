#!/usr/bin/env bash
#
# Helper utility script for sourcing vital environment variables needed
# for building LM-X API wrapper for Go languge.
# Copyright (C) 2002-2013 X-Formation. All rights reserved.
#
# The information and source code contained herein is the exclusive
# property of X-Formation and may not be disclosed or reproduced
# in whole or in part without explicit written authorization
# from the company.

TOP=$(dirname $(readlink -f "${BASH_SOURCE[0]}"))

lmx-env() {
  local GIT_TOPLEVEL=
  local SEARCH_PATH=
  local LMX_INC_FILE=
  local LMX_INC_PATH=
  local LMX_LIB_FILE=
  local LMX_LIB_PATH=
  local SEARCH_PATH=()

  [ -z "${LMXSDK_PATH}" ] && {
    SEARCH_PATH=(
      /usr
      /opt
    )

    type -a git &>/dev/null && git status &>/dev/null && {
      GIT_TOPLEVEL=$(git rev-parse --show-toplevel 2>/dev/null)
      SEARCH_PATH+=("${GIT_TOPLEVEL}")
      cd "${GIT_TOPLEVEL}"/../ 2>/dev/null && git status &>/dev/null && {
        GIT_TOPLEVEL=$(git rev-parse --show-toplevel 2>/dev/null)
        SEARCH_PATH+=("${GIT_TOPLEVEL}")
      }
      cd "${TOP}"
    }
  } || {
    SEARCH_PATH=(
      "${LMXSDK_PATH}"
    )
  }

  for path in ${SEARCH_PATH[@]}; do
    LMX_INC_FILE=$(find "${path}" -type f -regex '.*\/include\/lmx\.h' -print -quit 2>/dev/null)
    LMX_LIB_FILE=$(find "${path}" -type f -regex '.*\/liblmxclient\.a' -print -quit 2>/dev/null)
    [ ! -z "${LMX_INC_FILE}" ] && [ ! -z "${LMX_LIB_FILE}" ] && {
      LMX_INC_PATH=$(dirname "${LMX_INC_FILE}")
      LMX_LIB_PATH=$(dirname "${LMX_LIB_FILE}")
      echo "${C_INCLUDE_PATH}" | tr ':' '\n' | grep -x "${LMX_INC_PATH}" &>/dev/null || {
        C_INCLUDE_PATH=${LMX_INC_PATH}:${C_INCLUDE_PATH}
      }
      echo "${LIBRARY_PATH}" | tr ':' '\n' | grep -x "${LMX_LIB_PATH}" &>/dev/null || {
        LIBRARY_PATH=${LMX_LIB_PATH}:${LIBRARY_PATH}
      }
      break
    } || {
      LMX_INC_FILE=
      LMX_LIB_FILE=
    }
  done

  die_missing() {
    echo "error: unable to find ${1}"
    echo "error: consider exporting LMXSDK_PATH pointing to your LM-X SDK"
    return 1
  }

  [ -z "${LMX_INC_FILE}" ] && {
    die_missing 'lmx.h'
    return 1
  }

  [ -z "${LMX_LIB_FILE}" ] && {
    die_missing 'liblmxclient.a'
    return 1
  }


  echo "${GOPATH}" | tr ':' '\n' | grep -x "${TOP}" &>/dev/null || {
    GOPATH=${TOP}:${GOPATH}
  }

  type -a xmllicgen &>/dev/null || {
    XMLLICGEN_PATH=$(dirname $(find "${LMX_LIB_PATH}" -type f -name xmllicgen 2>/dev/null))
    test -x "${XMLLICGEN_PATH}/xmllicgen" && {
      echo "${PATH}" | tr ':' '\n' | grep -x "${XMLLICGEN_PATH}" &>/dev/null || {
        PATH=${PATH}:${XMLLICGEN_PATH}
      }
      true
    } || {
      XMLLICGEN_PATH=
    }
  }

  echo "# added \"${LMX_INC_PATH}\" to your C_INCLUDE_PATH env"
  echo "# added \"${LMX_LIB_PATH}\" to your LIBRARY_PATH env"
  test ! -z "${XMLLICGEN_PATH}" && {
    echo "# added \"${XMLLICGEN_PATH}\" to your PATH env"
  } || {
    echo "! unable to find your xmllicgen executable (needed for testing only)"
  }

  export C_INCLUDE_PATH LIBRARY_PATH GOPATH PATH
}

lmx-env

