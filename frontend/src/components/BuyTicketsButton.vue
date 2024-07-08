
<script setup lang="ts"> 
  const emit = defineEmits(["ticketsBought"]);

  import IconBuyTickets from "@/components/icons/IconBuyTickets.vue"
  const props = defineProps<{
    eventId: number,
  }>();


  const buyTickets = async (_: any) => { 
    const good: boolean = await fetch(`/api/buyTickets?eventId=${props.eventId}&nSeats=1`, { 
      method: "POST",
      credentials: "same-origin",
    }).then(response => { 
      return response.ok;
    })
    .catch(e => { 
      console.error("error buying tickets");
      return false;
    })
    if (!good) { 
      return;
    }
    emit("ticketsBought");
  }
</script>

<template> 
    <button class="tooltip" @click="buyTickets"> 
      <span class="tooltiptext"> BuyTickets </span>
      <IconBuyTickets/>
    </button>
</template>

<style scoped> 
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
</style>