<template>
  <div id="app">
    <div :class="$style.container">
      <NavBar :id="$style.nav" />
      <div :id="$style.main">
        <template v-if="this.$store.state.error.message != null">
          <Error></Error>
        </template>
        <template v-else>
          <router-view />
        </template>
      </div>
      <Footer :id="$style.footer" />
    </div>
  </div>
</template>

<script>
import NavBar from "./components/NavBar";
import Footer from "./components/Footer";
import Error from "./components/Error";
export default {
  name: "app",
  components: {
    NavBar,
    Footer,
    Error,
  },
};
</script>

<style>
@import "../public/Skeleton-2.0.4/css/normalize.css";
@import "../public/Skeleton-2.0.4/css/skeleton.css";
body {
  font-family: "M PLUS 1p", sans-serif;
}

p {
  word-wrap: break-word;
}

textarea {
  min-height: 100px;
}

.material-icons {
  font-size: 2em;
}

.router-link-exact-active {
  font-weight: bold;
}

.card {
  background-color: #fff;
  box-shadow: 0 0 8px rgba(0, 0, 0, 0.16);
  padding: 1rem;
}

.add-card {
  background-color: #fff;
  border-radius: 4px;
  border: 1px solid #bbb;
  padding: 1rem;
  cursor: pointer;
  color: #1eaedb;
}

.add-card:hover {
  color: #0fa0ce;
  border-color: #888;
  outline: 0;
}

.error_message {
  color: #dc3545;
}

div .error_input {
  border-color: #dc3545;
}

.material-icons,
.icon-text {
  vertical-align: middle;
}
</style>

<style module>
.container {
  display: grid;
  grid-template-rows: max-content 1fr max-content;
  grid-template-columns: 1rem 1fr 1rem;
  row-gap: 1rem;
}

#nav {
  grid-row: 1 / 2;
  grid-column: 2 / 3;
}

#main {
  grid-row: 2 / 3;
  grid-column: 2 / 3;
}

#footer {
  grid-row: 3 / 4;
  grid-column: 2 / 3;
}
</style>
