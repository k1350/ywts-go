<template>
  <article>
    <div id="firebaseui-auth-container"></div>
    <div v-if="this.message">{{ message }}</div>
    <div id="loader">Loading...</div>
  </article>
</template>

<script>
require("../firebaseConfig.js");

var firebase = require("firebase/app");
var firebaseui = require("firebaseui-ja");
firebase.auth().setPersistence(firebase.auth.Auth.Persistence.NONE);

import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const AuthRepository = RepositoryFactory.get("auth");

export default {
  name: "SignIn",
  data: function() {
    return {
      message: null,
    };
  },
  async created() {
    await AuthRepository.getCsrfToken().then((token) => {
      let t = JSON.parse(token.data);
      if ("uid" in t) {
        // ログイン済みなのでセッションからユーザー情報を復元する
        var payload = new Object();
        payload.csrfToken = t.csrfToken;
        payload.uid = t.uid;
        payload.name = t.name;
        payload.email = t.email;
        this.$store.commit("auth/signin", payload);
      } else {
        this.$store.commit("auth/setCsrfToken", t.csrfToken);
      }
    });
    if (this.$store.state.auth.signined) {
      this.$router.push("/");
    }
  },
  mounted() {
    if (this.$store.state.auth.message != null) {
      this.message = this.$store.state.auth.message;
      this.$store.commit("auth/removeMessage");
    }

    let self = this;
    var uiConfig = {
      callbacks: {
        signInSuccessWithAuthResult: function(authResult) {
          self.signin(authResult);
        },
        uiShown: function() {
          document.getElementById("loader").style.display = "none";
        },
      },
      credentialHelper: firebaseui.auth.CredentialHelper.NONE,
      signInSuccessUrl: "/",
      signInOptions: [
        {
          provider: firebase.auth.EmailAuthProvider.PROVIDER_ID,
          requireDisplayName: true,
        },
      ],
    };
    var ui = firebaseui.auth.AuthUI.getInstance();
    if (!ui) {
      ui = new firebaseui.auth.AuthUI(firebase.auth());
    }
    ui.start("#firebaseui-auth-container", uiConfig);
  },
  methods: {
    signin: function(authResult) {
      let self = this;
      firebase
        .auth()
        .currentUser.getIdToken(/* forceRefresh */ true)
        .then(function(idToken) {
          return AuthRepository.signin(authResult.user.uid, idToken);
        })
        .then(function(token) {
          return new Promise(function(resolve) {
            let t = JSON.parse(token.data);
            var payload = new Object();
            payload.csrfToken = t.csrfToken;
            payload.uid = authResult.user.uid;
            payload.name = authResult.user.displayName;
            payload.email = authResult.user.email;
            self.$store.commit("auth/signin", payload);
            resolve();
          });
        })
        .then(function() {
          firebase.auth().signOut();
        })
        .then(function() {
          self.$router.push("/");
        })
        .catch(function(error) {
          console.log(error);
          self.$store.commit("auth/setMessage", "ログインに失敗しました。");
          firebase.auth().signOut();
          self.$router.go({
            path: self.$router.currentRoute.path,
            force: true,
          });
        });
    },
  },
};
</script>
