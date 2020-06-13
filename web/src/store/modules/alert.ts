import { ActionTree, Module, MutationTree, GetterTree } from "vuex";
import { Alert, AlertState, RootState } from "@/types";

const state: AlertState = {
  alert: {} as Alert
};

export const getters: GetterTree<AlertState, RootState> = {
  alert: (state: AlertState) => {
    return state.alert;
  }
};

const mutations: MutationTree<AlertState> = {
  setAlert(state: AlertState, alert: Alert) {
    state.alert = alert;
  }
};

export const actions: ActionTree<AlertState, RootState> = {
  setAlert({ commit }, alert: Alert) {
    commit("setAlert", alert);
  }
};

export const alert: Module<AlertState, RootState> = {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
};
