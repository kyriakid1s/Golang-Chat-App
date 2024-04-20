import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("../components/Login.vue"),
    },
    {
      path: "/chats",
      name: "chats",
      component: () => import("../components/Chats.vue"),
    },
  ],
});

export default router;
