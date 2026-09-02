<script setup lang="ts">
// Deck-level override of the theme's `end` layout. Identical to the theme
// version except the QR box shows a real image when a `qr` prop is given,
// falling back to the theme's [ QR ] placeholder otherwise.
defineProps<{
  url?: string
  handles?: Array<{ label: string; href?: string }>
  qr?: string
}>()
</script>

<template>
  <div class="slidev-layout layout-end">
    <div class="end-left">
      <span class="ac-label">// that's a wrap</span>
      <h1 class="ac-accent">Thank You</h1>
      <p class="ac-muted end-cta">Slides, code, and demo, all on GitHub:</p>
      <div v-if="handles && handles.length" class="chip-row">
        <a
          v-for="handle in handles"
          :key="handle.label"
          :href="handle.href || '#'"
          class="ac-chip ac-chip--active"
        >{{ handle.label }}</a>
      </div>
    </div>
    <div class="end-right">
      <div class="qr-placeholder" :class="{ 'qr-has-image': qr }">
        <img v-if="qr" :src="qr" alt="scan for slides" class="qr-img" />
        <span v-else class="qr-inner">[ QR ]</span>
      </div>
      <span class="qr-caption">scan for slides + repo</span>
    </div>
  </div>
</template>

<style scoped>
.layout-end {
  display: flex;
  flex-direction: row;
  align-items: center;
  height: 100%;
  padding: var(--ac-space-8);
}

.end-left {
  width: 60%;
  display: flex;
  flex-direction: column;
  gap: var(--ac-space-5);
  padding-right: var(--ac-space-7);
}

.end-cta {
  font-family: var(--ac-font-sans);
  font-size: 17px;
  margin: 0;
}

.chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: var(--ac-space-2);
}

.end-right {
  width: 40%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--ac-space-4);
}

.qr-placeholder {
  width: 200px;
  height: 200px;
  background: var(--ac-surface-high);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* White quiet zone behind the code so it scans against the dark slide. */
.qr-has-image {
  background: #ffffff;
  padding: 10px;
}

.qr-img {
  width: 100%;
  height: 100%;
  display: block;
  image-rendering: pixelated;
}

.qr-inner {
  font-family: var(--ac-font-mono);
  font-size: 14px;
  color: var(--ac-on-surface-muted);
}

.qr-caption {
  font-family: var(--ac-font-mono);
  font-size: 12px;
  color: var(--ac-on-surface-muted);
}
</style>
