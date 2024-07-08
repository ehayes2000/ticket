<script setup lang="ts"> 
  import { ref, onMounted } from "vue"
  import { type Tickets, type Game, type Concert} from "@/models/Event"
  import MyEvents from "@/components/MyEvents.vue"
  import MyTickets from "@/components/MyTickets.vue"

  const myTicketss = ref<Array<Tickets>>([])
  onMounted(async () => { 
    await fetchMyTickets();
  })
  const fetchMyTickets = async () => {
    const newTicketss = await fetch("/api/getAllTickets", {
      method: "GET",
      credentials: "same-origin",
    }).then(response => { 
      if (!response.ok){
        throw new Error(`failed to get tickets ${response}`)
      }
      return response.json()
    })
    .then(data => {
      if (!data){
        console.error("no data: ", data)
        throw new Error("no data")
      }
      console.log("we get a bit of data", data)
      const typifiedData: Tickets[] = data.map((e: any)=> { 
        const typified: Tickets = e as Tickets
        typified.event.date = new Date(typified.event.date)
        return typified
      });
      return typifiedData
    })
    .catch(e => {
      console.error("something went wrong getting tickets:", e)
    })
    if (newTicketss) {
      console.log("WE GET TICKETS", newTicketss)
      myTicketss.value = newTicketss;
    }
  }

</script>

<template>
  <div class="wrapper">
    <div class="my-stuff">
      <h1> My Saved Events</h1>
      <MyEvents @ticketsBought="fetchMyTickets"/>
    </div>
    <div class="my-stuff">
      <h1> My Tickets </h1>
      <MyTickets :ticketss="myTicketss"/>
    </div>
  </div>
</template>

<style scoped> 
  h1 { 
    padding-top: 1rem;
    padding-left: 1rem;
  }
  .wrapper { 
    display: grid;
    grid-template-columns: 1fr 1fr;
    column-gap: 1rem;
  }

  .my-stuff { 
    outline: 1px solid #e0e0e0;
    border-radius: .5rem;   
    padding: .5rem;
  }

</style>
