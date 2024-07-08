<script setup lang="ts">
  import { RouterLink } from "vue-router";
</script>

<template>
  <div id="navBar">
    <div class="nav-group">
      <RouterLink class="boldify" to="/events"> Events </RouterLink>
    <nav id="pages">
      <RouterLink to="/myTickets"> My Tickets </RouterLink>
    </nav>
  </div>
    <div class="nav-group">
      <nav id="meta">
        <RouterLink to="/about"> About </RouterLink>
      </nav>
      <nav>
        <RouterLink to="/login"> Login </RouterLink>
      </nav>
    </div>
  </div>
</template>

<style scoped>

  * { 
    color: white
  }

  #navBar { 
    display: flex;
    background-color: blue;
    margin: 0;
    justify-content: space-between;
    padding: 1rem 2rem 1rem 2rem;
  }

  .nav-group { 
    display: flex;
    flex-direction: row;  
    flex-wrap: nowrap;    
    column-gap: 2rem;
  }

  .boldify { 
    font-weight:bold;
    font-size: 16
  }
</style>