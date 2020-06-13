import Vue from "vue";
import Vuex, { StoreOptions } from "vuex";
import { RootState } from "@/types";

import { alert } from "./modules/alert";
import { fuelTank } from "./modules/fuelTank";
import { truck } from "./modules/truck";

Vue.use(Vuex);

const store: StoreOptions<RootState> = {
  modules: {
    alert,
    fuelTank,
    truck
  }
};

export default new Vuex.Store<RootState>(store);
