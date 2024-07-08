<script setup lang="ts">
  import { useRouter } from "vue-router";
  import { authStore } from '@/store/auth';

  const router = useRouter()
  const logout = async () => {
    await fetch("/api/logout", {
      method:"GET", 
    })
    authStore.isLoggedIn = false;
    router.push("/login");
  }

</script>


<template>
  <div class="wrapper">
  <button @click="logout">
    Logout
  </button>
  </div>
</template>

<style scoped>
  .wrapper { 
    display: flex;
    align-items: center;
    justify-content: center;
  }
  button {
    color: blue;
    background: none;
    border: solid blue 1px;
    border-radius: .2rem;
    padding: 0 1rem 0 1rem;
    margin: 0;
    font: inherit;
    cursor: pointer;
    outline: none;
    transition: 0.14s ease;
  }

  
</style>