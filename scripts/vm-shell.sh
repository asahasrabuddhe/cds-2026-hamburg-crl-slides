#!/usr/bin/env bash
#
# Entry point for the in-slide terminals (slidev-addon-live-terminal).
# Lands in ~/crl on a demo VM, as the demo user or as root.
#
#   vm-shell.sh [primary|cgroupv1|hardened] [root]
#
# Uses the ssh_config that ./qemu/vm.sh up/push writes in the code repo, so
# the VM must be up and pushed first (EVENT-DAY.md, T-90). ControlMaster keeps
# reconnects instant as slides with terminals mount and unmount.
set -euo pipefail

VARIANT="${1:-primary}"
ROLE="${2:-}"

CODE_REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../cds-2026-hamburg-crl-code" 2>/dev/null && pwd || true)"
SSH_CONFIG="$CODE_REPO/qemu/ssh_config"
if [[ -z "$CODE_REPO" || ! -f "$SSH_CONFIG" ]]; then
  echo "no VM ssh_config found."
  echo "from cds-2026-hamburg-crl-code run: make vm && make push"
  exit 1
fi

REMOTE='cd ~/crl && exec bash -l'
if [[ "$ROLE" == "root" ]]; then
  # sudo bash, not sudo -i: keep ~/crl as the working directory, matching
  # what stage.sh does in the rootful tmux pane.
  REMOTE='cd ~/crl && exec sudo bash'
fi

exec ssh -F "$SSH_CONFIG" \
  -o ControlMaster=auto \
  -o "ControlPath=$HOME/.ssh/crl-%h-%p" \
  -o ControlPersist=600 \
  -t "crl-${VARIANT}" "$REMOTE"
