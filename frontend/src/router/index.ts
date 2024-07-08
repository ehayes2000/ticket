import { createRouter, createWebHistory } from 'vue-router'
import { authStore } from '@/store/auth'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import("../views/EventsView.vue")
    },
    {
      path: "/myTickets",
      name: "myTickets",
      component: () => import("../views/MyTicketsView.vue")
    },
    {
      path: "/events",
      name: "events",
      component: () => import("../views/EventsView.vue")
    },
    {
      path: "/login",
      name: "login",
      component: () => import("../views/LoginView.vue")  
    },
    {
      path: "/logout",
      name: "logout",
      component: () => import("../views/LogoutView.vue")
    }
  ]
})

export default router
