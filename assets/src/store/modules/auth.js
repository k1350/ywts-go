const state = {
  signined: false,
  csrfToken: null,
  uid: null,
  name: null,
  email: null,
  message: null,
};

const mutations = {
  signin(state, payload) {
    state.signined = true;
    state.csrfToken = payload.csrfToken;
    state.uid = payload.uid;
    state.name = payload.name;
    state.email = payload.email;
  },
  signout(state) {
    state.signined = false;
    state.csrfToken = null;
    state.uid = null;
    state.name = null;
    state.email = null;
  },
  setMessage(state, payload) {
    state.message = payload;
  },
  removeMessage(state) {
    state.message = null;
  },
  setCsrfToken(state, payload) {
    state.csrfToken = payload;
  },
};

export default {
  namespaced: true,
  state,
  mutations,
};
