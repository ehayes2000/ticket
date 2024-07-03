<script setup lang="ts"> 
  import { ref } from "vue";
  const username = ref("")
  const password = ref("");
  const API = import.meta.env.VITE_API;

  function handleLogin(_) {
    const formData = new FormData();
    if (!username.value || !password.value)
      return;
    formData.append("username", username.value);
    formData.append("password", password.value);
    fetch(`${API}/login`, {
      method: "POST",
      body: formData
    }).then(response => {
      console.log("resp,", response);
      if (!response.ok)
        throw new Error(`HTTP Error ${response.status}`);
      return response.json();
    }).then(data => {
      console.log(data);
    })
    .catch(e => {
      console.error("Request failed with error:", e);
    });
  }
</script>

<template> 
  <div>
    <input type="text" v-model="username"/>
    <input type="password" v-model="password"/>
    <button @click="handleLogin"> Login </button>
    
  </div>
</template>