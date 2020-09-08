import Repository from "./Repository";
import store from "../store/index";

const resource = "/v1/boards";
export default {
  get() {
    return Repository.get(`${resource}`);
  },
  getBoard(id) {
    return Repository.get(`${resource}/${id}`);
  },
  addBoard(payload) {
    return Repository.post(`${resource}`, payload, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  deleteBoard(id) {
    return Repository.delete(`${resource}/${id}`, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  updateBoard(id, payload) {
    return Repository.put(`${resource}/${id}`, payload, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  getBoardId(code) {
    return Repository.get(`${resource}/${code}/id`);
  },
  getMembers(id) {
    return Repository.get(`${resource}/${id}/members`);
  },
  joinBoard(id, payload) {
    return Repository.post(`${resource}/${id}/members`, payload, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  deleteBoardMember(id, uid) {
    return Repository.delete(`${resource}/${id}/members/${uid}`, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  getItems(id) {
    return Repository.get(`${resource}/${id}/items`);
  },
  getItem(id, iid) {
    return Repository.get(`${resource}/${id}/items/${iid}`);
  },
  addItem(id, payload) {
    return Repository.post(`${resource}/${id}/items`, payload, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  updateItem(id, iid, payload) {
    return Repository.put(`${resource}/${id}/items/${iid}`, payload, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
  deleteItem(id, iid) {
    return Repository.delete(`${resource}/${id}/items/${iid}`, {
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": `${store.state.auth.csrfToken}`,
      },
    });
  },
};
