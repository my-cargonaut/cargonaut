import Vue from "vue";
import VueRouter from "vue-router";
import Calculate from "@/views/Calculate.vue";
import Trucks from "@/views/Trucks.vue";
import FuelTanks from "@/views/FuelTanks.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "calculate",
    component: Calculate,
    meta: {
      Title: "CARGONAUT - Calculate"
    }
  },
  {
    path: "/trucks",
    name: "trucks",
    component: Trucks,
    meta: {
      Title: "CARGONAUT - Trucks"
    }
  },
  {
    path: "/tanks",
    name: "fuel-tanks",
    component: FuelTanks,
    meta: {
      Title: "CARGONAUT - Fuel Tanks"
    }
  },
  {
    path: "*",
    name: "not-found",
    component: () =>
      import(/* webpackChunkName: "about" */ "@/views/NotFound.vue"),
    meta: {
      Title: "CARGONAUT - Not Found"
    }
  }
];

const router = new VueRouter({
  mode: "hash", // history
  base: process.env.BASE_URL,
  routes
});

export default router;
