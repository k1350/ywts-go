<template>
  <ul :class="$style.container">
    <li :class="$style.link_nav" v-if="this.$store.state.auth.signined">
      <router-link to="/" exact>
        ホーム
      </router-link>
    </li>
    <li :class="$style.link_nav" v-if="this.$store.state.auth.signined">
      <a href="#" @click="signOut">ログアウト</a>
    </li>
  </ul>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const AuthRepository = RepositoryFactory.get("auth");

export default {
  name: "NavBar",
  methods: {
    signOut: function() {
      let self = this;
      AuthRepository.signout()
        .then(function() {
          new Promise(function(resolve) {
            self.$store.commit("auth/signout");
            resolve();
          });
        })
        .then(function() {
          self.$router.push("/signin");
        })
        .catch(function(error) {
          console.log(error);
        });
    },
  },
};
</script>

<style module>
.container {
  display: grid;
  grid-template-rows: 1fr;
  grid-template-columns: 10rem 10rem 1fr;
  padding: 1rem 0;
}
</style>
