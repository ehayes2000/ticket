
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
  <div class="home-layout">
    <div class="my-items">
      <h1> Saved Events </h1>
      <EventItem class="fade-out" v-for="e in mySavedEvents" :event="e">
        <div class="event-controls">
          <BuyTicketsButton :eventId="e.id"/>
          <UnsaveButton @click="()=>removeEvent(e.id)" :eventId="e.id" />
        </div>
      </EventItem>
    </div>
  </div>
</template>

<style scoped> 
  h1 { 
    color: blue
  }
  .home-layout {
    display: flex;
    justify-content: center;
    column-gap: 4rem;
  }
  .my-items { 
    outline: black solid 1px;  
    padding: 0 2rem 0 2rem ;
  }
  .event-controls { 
    display: flex;
    flex-direction: column;
    padding: .5rem;
    row-gap: .5rem;
    justify-content: center;
  }
</style>
