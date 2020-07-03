import vehiclesAPI from "@/api/vehicles";

const vehicles = {
  namespaced: true,

  state: {
    loading: false,
    vehicles: [],
    vehicle: {}
  },

  mutations: {
    SET_LOADING(state, loading) {
      state.loading = loading;
    },
    SET_VEHICLES(state, vehicles) {
      state.vehicles = vehicles;
    },
    SET_VEHICLE(state, vehicle) {
      state.vehicle = vehicle;
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

    get({ commit }, id) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        vehiclesAPI
          .get(id)
          .then(response => {
            const vehicle = response.data;
            commit("SET_VEHICLE", vehicle);
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
    vehicles: state => (state.vehicles ? state.vehicles : []),
    vehicle: state => (state.vehicle ? state.vehicle : {})
  }
};

export default vehicles;

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
