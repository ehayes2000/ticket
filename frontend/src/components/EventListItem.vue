<script setup lang="ts">
  import { formatDate } from "@/util/formatDate";
  import { type Game, type Concert} from "@/models/Event";
  import PrettyDate from "@/components/PrettyDate.vue";
  import EventTag from "@/components/EventTag.vue";

  const props = defineProps<{
    event: Game | Concert
  }>();
  const event = props.event
</script>

<template>
  <div class="event">
    <PrettyDate :date="props.event.date"/>
    <div class="details"> 
      <div class="header">
        <h1> {{ props.event.name }} </h1>
        <h> <EventTag :kind="props.event.kind" /> </h>
      </div>
      <div class="more-details"> {{ props.event.venue }} </div>
      <div class="more-details"> {{ props.event.description }} </div>
      <div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .event { 
    display: flex;
    flex-direction: row;
    cursor: pointer;
  }

  .header { 
    display: flex;
    column-gap: .3rem;

  }
  .details {
    padding: .5rem;
  }
  .more-details { 
    max-width: 150px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>