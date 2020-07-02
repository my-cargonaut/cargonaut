import vehiclesAPI from "@/api/vehicles";

const vehicles = {
  namespaced: true,

  state: {
    loading: false,
    vehicles: []
  },

  mutations: {
    SET_LOADING(state, loading) {
      state.loading = loading;
    },
    SET_VEHICLES(state, vehicles) {
      state.vehicles = vehicles;
    }
  },

  actions: {
    list({ commit }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        vehiclesAPI
          .list()
          .then(response => {
            const vehicles = response.data;
            commit("SET_VEHICLES", vehicles);
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    create({ commit }, vehicle) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        vehiclesAPI
          .create(vehicle)
          .then(response => {
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    update({ commit }, { id, vehicle }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        vehiclesAPI
          .update(id, vehicle)
          .then(response => {
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    delete({ commit }, id) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        vehiclesAPI
          .delete(id)
          .then(response => {
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    }
  },

  getters: {
    loading: state => (state.loading ? state.loading : false),
    vehicles: state => (state.vehicles ? state.vehicles : [])
  }
};

function getAlert(e) {
  const alert = {
    type: "error",
    message: "Something went wrong!",
    title: true
  };
  if (e.response && e.response.data && e.response.data.error) {
    alert.message = e.response.data.error;
  }
  return alert;
}

export default vehicles;
