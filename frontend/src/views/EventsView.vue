<script setup lang="ts">
  import { type Game, type Concert } from "@/models/Event";
  import { onMounted, ref } from "vue";
  import { getEvents } from '@/store/events';
  import EventItem from "@/components/EventItem.vue"
  import SaveButton from '@/components/SaveButton.vue';
  import BuyTicketsButton from "@/components/BuyTicketsButton.vue";

  let events = ref<(Game | Concert)[]>([]); 

  onMounted(async () => { 
    const newEvents = await getEvents();
    events.value = newEvents;  
  })
</script> 

<template>   
  <div class="event-list-wrapper">
    <div class="event-list">
      <EventItem v-for="e in events" :key="e.id" :event="e">
        <div class="event-controls">
          <BuyTicketsButton :eventId="e.id"/>
          <SaveButton :eventId="e.id"/>
        </div> 
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
  .event-controls { 
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-content: center;
    row-gap: .5rem;
    padding: .5rem;
  }
</style>