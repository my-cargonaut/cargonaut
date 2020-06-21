import Vue from "vue";
import Vuex from "vuex";

import alert from "./modules/alert";
import auth from "./modules/auth";
import users from "./modules/users";

Vue.use(Vuex);

export default new Vuex.Store({
  strict: true,

  modules: {
    alert,
    auth,
    users
  }
});
