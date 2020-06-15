import Vue from "vue";
import Vuex from "vuex";

import alert from "./modules/alert";
import auth from "./modules/auth";

Vue.use(Vuex);

export default new Vuex.Store({
  strict: true,

  modules: {
    alert,
    auth
  }
});
