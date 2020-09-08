<template>
  <article :class="$style.container">
    <h1 :class="$style.title">{{ this.$store.state.error.message }}</h1>
    <p :class="$style.status">
      エラーコード: {{ this.$store.state.error.status }}
    </p>
    <p :class="$style.button">
      <button @click="goHome">戻る</button>
    </p>
  </article>
</template>

<script>
export default {
  name: "Error",
  methods: {
    goHome() {
      this.$router.go({name: 'home', force: true});
    },
  },
};
</script>

<style module>
.container {
  display: grid;
  grid-auto-rows: max-content max-content max-content;
  grid-template-columns: 1fr minmax(300, 600) 1fr;
  justify-content: center;
}

.title {
  grid-row: 1 / 2;
  grid-column: 2 / 3;
}

.status {
  grid-row: 2 / 3;
  grid-column: 2 / 3;
  text-align: center;
}

.button {
  grid-row: 3 / 4;
  grid-column: 2 / 3;
  text-align: center;
}
</style>
