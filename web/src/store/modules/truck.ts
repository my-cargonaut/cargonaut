import { ActionTree, Module, MutationTree, GetterTree } from "vuex";
import { Alert, Truck, TruckState, RootState } from "@/types";
import TruckAPI from "@/services/TruckAPI";

const state: TruckState = {
  trucks: [],
  loading: true
};

export const getters: GetterTree<TruckState, RootState> = {
  trucks: (state: TruckState) => {
    return state.trucks;
  },
  loading: (state: TruckState) => {
    return state.loading;
  }
};

const mutations: MutationTree<TruckState> = {
  setTrucks(state: TruckState, trucks: Truck[]) {
    state.trucks = trucks;
  },
  setLoading(state: TruckState, loading: boolean) {
    state.loading = loading;
  }
};

export const actions: ActionTree<TruckState, RootState> = {
  listTrucks({ commit }) {
    commit("setLoading", true);
    TruckAPI.list()
      .then((trucks: Truck[]) => {
        commit("setTrucks", trucks);
      })
      .catch(e => {
        const alert = {
          type: "error",
          message: "Something went wrong!"
        } as Alert;
        if (e.response && e.response.data && e.response.data.error) {
          alert.message = e.response.data.error;
        }
        commit("alert/setAlert", alert, { root: true });
      })
      .finally(() => {
        commit("setLoading", false);
      });
  }
};

export const truck: Module<TruckState, RootState> = {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
};
