<script setup lang="ts">
import { EventKind } from "@/models/Event";

const props = defineProps<{
  kind: EventKind
}>();

const formatKind = (k: string) => {
  return k[0].toUpperCase() + k.slice(1).toLowerCase();
};

const styleSwitcher = (k: EventKind) => {
  if (k === EventKind.Game) {
    return "rgb(238, 99, 82)"
  }
  if (k === EventKind.Concert) {
    return "rgb(91, 192, 190)";
  }
  return ""; // Default case
};
</script>

<template>
  <div :style="{ borderRadius: '0.2rem', backgroundColor: styleSwitcher(props.kind) }">
    {{ formatKind(props.kind) }}
  </div>
</template>

<style scoped>
  * { 
    color: white;
    font-style: italic;
    font-size: medium;
    display: inline-block;
    padding:  0 .3rem 0 .3rem;
  }
</style>