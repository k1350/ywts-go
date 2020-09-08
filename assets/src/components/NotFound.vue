<template>
  <article :class="$style.container">
    <h1 :class="$style.title">404 Not Found</h1>
    <p :class="$style.message">お探しの画面は見つかりませんでした。</p>
    <p :class="$style.button">
      <router-link :to="{ name: 'home' }" class="button">戻る</router-link>
    </p>
  </article>
</template>

<script>
export default {
  name: "NotFound",
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

.message {
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
