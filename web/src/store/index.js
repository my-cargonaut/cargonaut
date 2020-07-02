import Vue from "vue";
import Vuex from "vuex";

import alert from "./modules/alert";
import auth from "./modules/auth";
import trips from "./modules/trips";
import users from "./modules/users";
import vehicles from "./modules/vehicles";

Vue.use(Vuex);

export default new Vuex.Store({
  strict: true,

  modules: {
    alert,
    auth,
    trips,
    users,
    vehicles
  }
});
