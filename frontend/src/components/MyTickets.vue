<script setup lang="ts">
  import { type Tickets } from "@/models/Event"
  import TicketsComponent from "@/components/TicketsComponent.vue";

  const props = defineProps<{
    ticketss: Array<Tickets>
  }>();

</script>

<template> 
  <TicketsComponent class="item" v-for="(t, i) in ticketss" 
    :key="i" 
    :tickets="t"/> 
</template>

<style scope>
  .item { 
    row-gap: 1rem;
  }
</style>

