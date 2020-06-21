import Vue from "vue";
import VueRouter from "vue-router";

import store from "@/store";

Vue.use(VueRouter);

const routes = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/Login.vue"),
    meta: {
      Title: "My Cargonaut - Login"
    }
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("@/views/Register.vue"),
    meta: {
      Title: "My Cargonaut - Register"
    }
  },
  {
    path: "/",
    name: "dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: {
      Title: "My Cargonaut - Dashboard",
      requiresAuth: true
    }
  },
  {
    path: "/users/:id",
    name: "profile",
    component: () => import("@/views/User.vue"),
    props: true,
    meta: {
      Title: "My Cargonaut - Profile",
      requiresAuth: true
    }
  },
  {
    path: "*",
    name: "not-found",
    component: () => import("@/views/NotFound.vue"),
    meta: {
      Title: "My Cargonaut - Not Found"
    }
  }
];

const router = new VueRouter({
  routes
});

router.beforeEach((to, from, next) => {
  if (!to.matched.some(record => record.meta.requiresAuth)) {
    next();
    return;
  }

  if (!store.getters["auth/isLoggedIn"]) {
    next("/login");
    return;
  }

  // Reresh if token expires in less than 12 hours.
  const twelveHours = 1000 * 60 * 60 * 12; // In Milliseconds
  let now = +new Date();
  if (store.getters["auth/tokenExpiry"] - now < twelveHours) {
    store.dispatch("auth/refresh");
  }
  next();
});

export default router;
