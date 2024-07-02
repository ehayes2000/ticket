<script lang="ts">
  import { RouterLink } from "vue-router";
</script>

<template>
  <div id="navBar">
    <div id="leftContent">
      <RouterLink class="homeLink" to="/"> Tickets </RouterLink>
    <nav id="pages">
      <RouterLink to="/events"> Events </RouterLink>
    </nav>
  </div>
    <nav id="meta">
      <RouterLink to="/about"> About </RouterLink>
    </nav>
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

  #leftContent { 
    display: flex;
    flex-direction: row;  
    flex-wrap: nowrap;    
    column-gap: 2rem;
  }

  .homeLink { 
    font-weight:bold;
    font-size: 16
  }
</style>