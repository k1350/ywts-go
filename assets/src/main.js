import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store/index";
import sanitizeHTML from 'sanitize-html';

Vue.config.productionTip = false;
Vue.prototype.$sanitize = sanitizeHTML;

Vue.config.errorHandler = function(err, vm, info) { // eslint-disable-line
  let status = err.response.status;
  if (status >= 400 && status < 500) {
    store.commit("error/set400", status);
  } else if (status >= 500 && status < 600) {
    store.commit("error/set500", status);
  }
};

// 残りのエラーをキャッチ
window.addEventListener("error", (event) => {
  console.log("Captured in error EventListener", event.error);
});
window.addEventListener("unhandledrejection", (event) => {
  let status = event.reason.response.status;
  if (status >= 400 && status < 500) {
    store.commit("error/set400", status);
  } else if (status >= 500 && status < 600) {
    store.commit("error/set500", status);
  }
});

let app;

if (!app) {
  new Vue({
    router,
    store,
    render: (h) => h(App),
  }).$mount("#app");
}
