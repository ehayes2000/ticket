<script setup lang="ts">
  import { ref, onMounted } from "vue";
  import {
    type Concert, 
    type Game, 
  } from "@/models/Event";
  import { getMyEvents } from "@/store/myEvents"
  import EventItem from "@/components/EventItem.vue";
  import UnsaveButton from "@/components/UnsaveButton.vue";
  import BuyTicketsButton from "@/components/BuyTicketsButton.vue";

  const emit = defineEmits(["ticketsBought"]);
  const loading = ref(true);
  const mySavedEvents = ref<(Concert | Game)[]>([]);
  onMounted(async () => { 
    try { 
      mySavedEvents.value = await getMyEvents();
    } catch (e) { 
      console.log("gg go next")
    } finally { 
      loading.value = false;
    }
  })

  const removeEvent = (id: number) => { 
    console.log("filter for id", id)
    mySavedEvents.value = mySavedEvents.value.filter(e => {
      return e.id != id;
    })  
  }
  
  getMyEvents().then(e => {
    mySavedEvents.value = e
  })

</script>

<template>   
    <div class="event-list">
      <EventItem v-for="e in mySavedEvents" class="item" :key="e.id" :event="e">
        <div class="event-controls">
          <BuyTicketsButton @ticketsBought="emit('ticketsBought')" :eventId="e.id"/>
          <UnsaveButton @click="()=>removeEvent(e.id)" :eventId="e.id" />
        </div> 
      </EventItem>
    </div>
</template>

<style scoped> 
  .event-list {   
    justify-content: center;
  }
  .event-controls { 
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-content: center;
    row-gap: .5rem;
    padding-left: .5rem;
    padding-right: .5rem;
  }

.item {
  padding: 10px 0;
  position: relative;
}

.item:not(:last-child)::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background-color: #e0e0e0;
}
</style>
