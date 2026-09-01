<script setup lang="ts">
import { computed } from 'vue'

// A live shell on one of the demo VMs, wrapped in the theme's terminal
// chrome. Red top border and label for rootful, green for rootless, the same
// convention the tmux panes carry. The shell lands in ~/crl via
// scripts/vm-shell.sh, so ./scripts/demo.sh and ./nsdemo work as typed in the
// speaker notes.
const props = withDefaults(defineProps<{
  vm?: 'primary' | 'cgroupv1' | 'hardened'
  root?: boolean
  fontSize?: number
  autoConnect?: boolean
}>(), {
  vm: 'primary',
  root: false,
  fontSize: 15,
  autoConnect: true,
})

const cmd = computed(() => `./scripts/vm-shell.sh ${props.vm}${props.root ? ' root' : ''}`)
const label = computed(() => (props.root ? 'ROOTFUL (uid 0)' : 'ROOTLESS'))
</script>

<template>
  <div class="demo-terminal" :class="root ? 'dt-rootful' : 'dt-rootless'">
    <TerminalWindow :label="label" :command="vm">
      <LiveTerminal :cmd="cmd" :font-size="fontSize" :auto-connect="autoConnect" />
    </TerminalWindow>
  </div>
</template>

<style scoped>
.demo-terminal {
  height: 100%;
  border-top: 2px solid;
}

.dt-rootful  { border-color: var(--ac-error); }
.dt-rootless { border-color: var(--ac-success); }

.dt-rootful  :deep(.chrome-label) { color: var(--ac-error); }
.dt-rootless :deep(.chrome-label) { color: var(--ac-success); }

/* The theme pads .terminal-body for static slot content; the live terminal
   brings its own geometry, so pull the padding back in. */
.demo-terminal :deep(.terminal-body) {
  padding: var(--ac-space-2);
  overflow: hidden;
}
</style>
