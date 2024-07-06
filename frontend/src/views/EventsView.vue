<script setup lang="ts">
  import { type Game, type Concert } from "@/models/Event";
  import { onMounted, ref } from "vue";
  import { getEvents } from '@/store/events';
  import EventItem from "@/components/EventItem.vue"
  import SaveButton from '@/components/SaveButton.vue';

  let events = ref<(Game | Concert)[]>([]); 

  onMounted(async () => { 
    const newEvents = await getEvents();
    events.value = newEvents;  
  })
</script> 

<template>   
  <div class="event-list-wrapper">
    <div class="event-list">
      <EventItem v-for="e in events" :event="e"> 
        <SaveButton :eventId="e.id", />
      </EventItem>
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