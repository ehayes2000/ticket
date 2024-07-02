<script setup lang="ts">
  import EventListItem from "@/components/EventListItem.vue";
  import { useEventStore, getEvents } from '@/store/events';

  let events = useEventStore();
  getEvents().then(data => { 
    events.events = [...events.events, ...data];
  })
</script> 

<template>   
  <div>
    <EventListItem v-for="e in events.events" :event="e" />
  </div>
</template>