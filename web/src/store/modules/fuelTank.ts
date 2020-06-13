import { ActionTree, Module, MutationTree, GetterTree } from "vuex";
import { Alert, FuelTank, FuelTankState, RootState } from "@/types";
import FuelTankAPI from "@/services/FuelTankAPI";

const state: FuelTankState = {
  fuelTanks: [],
  loading: true
};

export const getters: GetterTree<FuelTankState, RootState> = {
  fuelTanks: (state: FuelTankState) => {
    return state.fuelTanks;
  },
  fuelTankByID: (state: FuelTankState) => (id: string) => {
    return state.fuelTanks.find(fuelTank => fuelTank.id === id);
  },
  fuelTanksForTruck: (state: FuelTankState) => (truckId: string) => {
    return state.fuelTanks.filter(fuelTank => fuelTank.truckId === truckId);
  },
  loading: (state: FuelTankState) => {
    return state.loading;
  }
};

const mutations: MutationTree<FuelTankState> = {
  setFuelTanks(state: FuelTankState, fuelTanks: FuelTank[]) {
    state.fuelTanks = fuelTanks;
  },
  setLoading(state: FuelTankState, loading: boolean) {
    state.loading = loading;
  }
};

export const actions: ActionTree<FuelTankState, RootState> = {
  listFuelTanks({ commit }) {
    commit("setLoading", true);
    FuelTankAPI.list()
      .then((fuelTanks: FuelTank[]) => {
        commit("setFuelTanks", fuelTanks);
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

export const fuelTank: Module<FuelTankState, RootState> = {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
};
