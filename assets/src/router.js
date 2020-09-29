import Vue from "vue";
import VueRouter from "vue-router";
import Home from "./components/Home";
import SignIn from "./components/SignIn";
import AddItem from "./components/AddItem";
import EditItem from "./components/EditItem";
import Board from "./components/Board";
import AddBoard from "./components/AddBoard";
import EditBoard from "./components/EditBoard";
import InviteBoard from "./components/InviteBoard";
import Account from "./components/Account";
import NotFound from "./components/NotFound";
import store from "./store/index";

Vue.use(VueRouter);

let router = new VueRouter({
  mode: "history",
  routes: [
    {
      path: "/",
      name: "home",
      component: Home,
    },
    {
      path: "/signin",
      name: "signIn",
      component: SignIn,
      meta: { isPublic: true },
    },
    {
      path: "/boards/:id",
      name: "board",
      component: Board,
      props: true,
    },
    {
      path: "/boards/:id/add-item",
      name: "add",
      component: AddItem,
      props: true,
    },
    {
      path: "/boards/:id/edit-item/:iid",
      name: "edit",
      component: EditItem,
      props: true,
    },
    {
      path: "/add-board",
      name: "addBoard",
      component: AddBoard,
    },
    {
      path: "/edit-board/:id",
      name: "editBoard",
      component: EditBoard,
      props: true,
    },
    {
      path: "/invite-board/:code",
      name: "invite",
      component: InviteBoard,
      props: true,
    },
    {
      path: "/account",
      name: "account",
      component: Account,
    },
    { path: "*", component: NotFound },
  ],
});

router.beforeEach((to, from, next) => {
  store.commit("error/remove");
  // isPublic でない場合(=認証が必要な場合)、かつ、ログインしていない場合
  if (
    to.matched.some((record) => !record.meta.isPublic) &&
    !store.state.auth.signined
  ) {
    next({ path: "/signin", query: { redirect: to.fullPath } });
  } else if (
    to.matched.some((record) => record.meta.isPublic) &&
    store.state.auth.signined
  ) {
    next({ path: "/" });
  } else {
    next();
  }
});

export default router;
