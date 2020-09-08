import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";
import auth from "./modules/auth";
import error from "./modules/error";

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    auth,
    error,
  },
  strict: true,
  plugins: [
    createPersistedState({
      key: "ywts",
      paths: [
        "auth.signined",
        "auth.uid",
        "auth.name",
        "auth.email",
        "auth.csrfToken",
        "auth.message",
      ],
      storage: window.sessionStorage,
    }),
  ],
});
