#!/bin/bash
set -e
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
CONFIGS=${SCRIPT_DIR}/configs/buildkit

DEPS=uidmap

BUILDKIT_VERSION=0.12.4
INSTALL_DIR=${HOME}/.buildkit-rootless
CACHE_DIR=${HOME}/.cache/buildkit
USER_SYSTEMD=${HOME}/.config/systemd/user
BUILDKIT_BUNDLE=/tmp/buildkit-v${BUILDKIT_VERSION}.tgz

mkdir -p ${INSTALL_DIR} ${USER_SYSTEMD} ${CACHE_DIR}

curl -Lo ${BUILDKIT_BUNDLE} https://github.com/moby/buildkit/releases/download/v${BUILDKIT_VERSION}/buildkit-v${BUILDKIT_VERSION}.linux-amd64.tar.gz
tar -C ${INSTALL_DIR} -xvzf ${BUILDKIT_BUNDLE}

rm -rf ${BUILDKIT_BUNDLE}

cp ${CONFIGS}/buildkit.service ${USER_SYSTEMD}/

systemctl --user enable buildkit.service
systemctl --user start buildkit.service
