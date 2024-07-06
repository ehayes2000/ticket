<script setup lang="ts">
  import { formatDate } from "@/util/formatDate";
  import { type Game, type Concert} from "@/models/Event";
  import PrettyDate from "@/components/PrettyDate.vue";
  import EventTag from "@/components/EventTag.vue";

  const props = defineProps<{
    event: Game | Concert
  }>();
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
    </div>

    <slot>
    </slot>
  </div>
</template>

<style scoped>
  .event { 
    display: grid;
    grid-template-columns: 1fr 4fr 1fr;
    justify-content: justify-between;
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
    max-width: 4fr;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

</style>