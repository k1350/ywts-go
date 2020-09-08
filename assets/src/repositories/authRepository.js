import Repository from "./Repository";
import store from "../store/index";

const resource = "/v1/auth";
export default {
  signin(uid, idToken) {
    return Repository.post(
      `${resource}/signin`,
      { uid: uid },
      {
        headers: {
          Authorization: `Bearer ${idToken}`,
          "X-CSRF-Token": `${store.state.auth.csrfToken}`,
        },
      }
    );
  },
  getCsrfToken() {
    return Repository.post(`${resource}/token`);
  },
  signout() {
    return Repository.get(`${resource}/signout`);
  },
};
