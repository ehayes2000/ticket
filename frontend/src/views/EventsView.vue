<script setup lang="ts">
  import SaveableEvent from "@/components/SaveableEvent.vue"
  import { useEventStore, getEvents } from '@/store/events';

  let events = useEventStore();
  getEvents().then(data => { 
    events.events = [...events.events, ...data];
  })
</script> 

<template>   
  <div class="event-list-wrapper">
    <div class="event-list">
      <SaveableEvent v-for="e in events.events" :event="e" />
    </div>
  </div> 
</template>

<style scoped> 
  .event-list-wrapper { 
    display: flex;
    justify-content: center;
  }
  .event-list {   
    justify-content: center;
  }
</style>