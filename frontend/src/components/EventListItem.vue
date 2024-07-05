<script setup lang="ts">
  import { formatDate } from "@/util/formatDate";
  import { type BaseEvent } from "@/models/Event";
  const props = defineProps<{
    event: BaseEvent
  }>();
  const event = props.event
</script>

<template>
  <div class="event">
    <h1>{{ event.name }}</h1>
    <div>
    <div class="event-info">
      <h2> {{ event.venue }} </h2>
      <h2> {{ formatDate(event.date) }} </h2>
      <h2 :class="['event-kind', event.kind]"> {{ event.kind }}</h2>
    </div>
    </div>
    <body> {{ event.description }}</body>
  </div>
</template>

<style scoped>
  .event { 
    display: flex;
    flex-direction: column;
    cursor: pointer;
    padding: 1rem;
  }

  .event-info { 
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    flex-direction: row;
  }

  .event-kind { 
    text-align: right;
  }

</style>