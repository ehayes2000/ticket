import { reactive } from "vue";

console.log("DEFINE AUTH STORE");
export const authStore = reactive({
  isLoggedIn: false,
})