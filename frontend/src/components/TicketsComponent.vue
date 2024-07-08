<script setup lang="ts">
import { type Tickets } from "@/models/Event"
import PrettyDate from "./PrettyDate.vue"
const props = defineProps<{
  tickets: Tickets
}>()
</script>

<template>
  <div class="wrapper">
    <div class="wrapper-1">
      <PrettyDate padding="0" :date="tickets.event.date" />
      <div class="wrapper-2">
        <h1>{{ tickets.event.name }}</h1>
        <div>{{ tickets.event.venue }}</div>
        <div class="dunno">
          <h1> Tickets </h1>
          <div class="seats">
            <div class="seat" v-for="s in tickets.seats">
              {{ s }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wrapper {
  padding: 0.5rem;
}

.seats {
  display: flex;
  gap: 0.2rem;
  flex-wrap: wrap;
}

.dunno { 
  display: flex;
  gap: .2rem;
}
.seat { 
  border: solid black 1px;
  border-radius: .1rem;
  padding: 0 .2rem;
}

.wrapper-1 {
  display: grid;
  column-gap: 4px;
  grid-template-columns: 1fr 5fr;
  padding: 4px;
  outline: solid 1px blue;
  border-radius: 0.2rem;
}

.wrapper-2 {
  /* background-color: blue; */
  border-radius: 0.2rem;
  /* color: white; */
  padding: 0.2rem 0.5rem;
}
</style>
