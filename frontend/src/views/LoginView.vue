<script setup lang="ts"> 
  import { ref } from "vue";
  const username = ref("")
  const password = ref("");
  const API = import.meta.env.VITE_API;
  
  function createAccount(_:any){
    const formData = new FormData();
    if (!username.value || !password.value)
      return;
    formData.append("username", username.value);
    formData.append("password", password.value);
    fetch(`/api/createAccount`, {
      method: "POST",
      credentials: "same-origin",
      body: formData
    }).then(response => {
      console.log("resp,", response);
      if (!response.ok)
        throw new Error(`HTTP Error ${response.status}`);
    })
    .catch(e => {
      console.error("Request failed with error:", e);
    });
  }
  function handleLogin(_: any) {
    const formData = new FormData();
    if (!username.value || !password.value)
      return;
    formData.append("username", username.value);
    formData.append("password", password.value);
    fetch(`/api/login`, {
      method: "POST",
      credentials: "same-origin",
      body: formData
    }).then(response => {
      console.log("resp,", response);
      if (!response.ok)
        throw new Error(`HTTP Error ${response.status}`);
    })
    .catch(e => {
      console.error("Request failed with error:", e);
    });
  }
  function saveEvent(_: any) { 
    fetch(`/api/saveEvent`, {
      method: "POST",
      credentials: "same-origin",
    }).then(response => {
      console.log("resp,", response);
      if (!response.ok)
        throw new Error(`HTTP Error ${response.status}`);
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
    <button @click="createAccount"> Create Account</button>
    <button @click="saveEvent"> TEST </button>

    
  </div>
</template>