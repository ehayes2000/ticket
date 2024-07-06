
<script setup lang="ts">
  import { ref, onMounted } from "vue";
  import {
    type Concert, 
    type Game, 
  } from "@/models/Event";
  import { getMyEvents } from "@/store/myEvents"
  import EventListItem from "@/components/EventListItem.vue";

  
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

  getMyEvents().then(e => {
    
    mySavedEvents.value = e
  })

</script>

<template>
  <div class="home-layout">
    <div class="my-items">
      <h1> Saved Events </h1>
      <EventListItem v-for="e in mySavedEvents" :event="e" />
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
</style>
