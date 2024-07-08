<script setup lang="ts">
  import NavBar from "./components/NavBar.vue"
</script>

<template>
  <NavBar class="nav-bar"/>
  <div id="content">
    <RouterView />
  </div>
</template>

<style scoped>

.nav-bar {
  /* width: 100%; */
  /* font-size: 12px;
  text-align: center;
  margin-top: 2rem; */
}



</style>
