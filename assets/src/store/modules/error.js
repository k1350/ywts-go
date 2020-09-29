const state = {
    message: null,
    status: null,
  };
  
  const mutations = {
    set400(state, payload) {
        state.message = "エラーが発生しました。";
        state.status = payload;
    },
    set500(state, payload) {
        state.message = "エラーが発生しました。しばらくたってからやり直してください。";
        state.status = payload;
    },
    remove() {
        state.message = null;
        state.status = null;
    }
  };
  
  export default {
    namespaced: true,
    state,
    mutations,
  };
  