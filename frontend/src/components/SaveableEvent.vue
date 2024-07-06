<!-- <script setup lang="ts"> 
  import { type Game, type Concert } from "@/models/Event";
  import EventListItem from "./EventItem.vue";
  import IconBuyTickets from "@/components/icons/IconBuyTickets.vue";
  import IconSaveEvent from "@/components/icons/IconSaveEvent.vue";

  const props = defineProps<{
    event: Game | Concert
  }>();
  function saveEvent() { 
    fetch(`/api/saveEvent?eventId=${props.event.id}`, {
      method: "POST",
      credentials: "same-origin"
    })
  }
</script>

<template> 
  <div class="wrapper">
    <EventListItem :event="props.event" />
    <div class="save-actions">
      <!-- <button class="tooltip"> 
       <span class="tooltiptext"> Buy Tickets </span>
         buy tickets
      </button> -->
      <button class="tooltip">
        <span class="tooltiptext"> Buy Tickets </span>
        <IconBuyTickets />
      </button>
      <button @click="saveEvent" class="tooltip">
        <span class="tooltiptext"> Save Event </span>
          <IconSaveEvent />
      </button>
    </div>
  </div>
</template>

<style scoped> 
  .save-actions { 
    display: flex;
    flex-direction: column;
    align-self: center;
    row-gap: .24rem;;
    padding-left: 1rem;
  }

  .wrapper {
    display: flex
  }
  button { 
	  background: none;
	  color: inherit;
	  border: none;
	  padding: 0;
	  font: inherit;
	  cursor: pointer;
	  outline: inherit;
  }

  .tooltip {
  position: relative;
  display: inline-block;
  }

/* Tooltip text */
.tooltip .tooltiptext {
  visibility: hidden;
  width: 120px;
  background-color: gray;
  color: #fff;
  text-align: center;
  padding: 2px 0;
  border-radius: 6px;
  top: -10px;
  left: 150%;
  /* Position the tooltip text - see examples below! */
  position: absolute;
  z-index: 1;
}
.tooltip .tooltiptext::after {
  content: " ";
  position: absolute;
  top: 50%;
  right: 100%; /* To the left of the tooltip */
  margin-top: -5px;
  border-width: 5px;
  border-style: solid;
  border-color: transparent gray transparent transparent;
}
.tooltip .tooltiptext {
  opacity: 0;
  transition: opacity 1s;
}

.tooltip:hover .tooltiptext {
  opacity: 1;
}
/* Show the tooltip text when you mouse over the tooltip container */
.tooltip:hover .tooltiptext {
  visibility: visible;
}
</style>  -->