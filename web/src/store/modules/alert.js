const alert = {
  namespaced: true,

  state: {
    alert: {}
  },

  mutations: {
    SET(state, alert) {
      state.alert = alert;
    }
  },

  actions: {
    set({ commit }, alert) {
      commit("SET", alert);
    }
  },

  getters: {
    alert: state => state.alert
  }
};

export default alert;
