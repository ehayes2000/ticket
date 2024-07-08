<script setup lang="ts"> 
  import { useRouter } from "vue-router";
  import { authStore } from "@/store/auth";
  import { onMounted, ref } from "vue";
  
  const username = ref("")
  const password = ref("");
  const router = useRouter();
  
  onMounted(() => {
    if (authStore.isLoggedIn){
      router.push("/logout")
    }
  })

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
      authStore.isLoggedIn = true;
      router.push("/logout");
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
      if (!response.ok)
        throw new Error(`HTTP Error ${response.status}`);
      authStore.isLoggedIn = true;
      router.push("/logout");
    })
    .catch(e => {
      console.error("Request failed with error:", e);
    });
  }
</script>

<template> 
  <div class="wrapper"> 
    <div class="login-view">
      <div class="header">
        <h1> Create an Account </h1>  
        or
        <h1> Login </h1>
      </div> 
      <input type="text" v-model="username" placeholder="Username"/>
      <input type="password" v-model="password" placeholder="Password"/>
      <div class="controls">
        <button @click="handleLogin"> Login </button>
        <button @click="createAccount"> Create Account</button>
      </div>
    </div>
  </div>
</template>

<style scoped> 
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

  button:hover { 
    color: white;
    background-color: blue;
  }

  .header { 
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .controls { 
    display: flex;
    justify-content: space-between;
    column-gap: 1rem;
  }
  .wrapper { 
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .login-view { 
    display: flex;
    flex-direction: column;     
    row-gap: .5rem;
  }

</style>
