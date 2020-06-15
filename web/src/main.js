import Vue from "vue";

import App from "@/App.vue";
import client from "@/api/client";
import router from "@/router";
import store from "@/store";
import vuetify from "@/plugins/vuetify";

Vue.prototype.$http = client;

const token = localStorage.getItem("token") || null;
if (token) {
  Vue.prototype.$http.defaults.headers.common[
    "Authorization"
  ] = `Bearer ${token}`;
}

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
