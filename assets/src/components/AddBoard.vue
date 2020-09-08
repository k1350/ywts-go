<template>
  <article>
    <h1>新規追加</h1>
    <form @submit.prevent="submit" :class="$style.container">
      <div :class="$style.item_container">
        <label for="title">
          タイトル
          <span :class="$style.small">※必須, 最大150文字</span>
        </label>
        <input
          type="text"
          id="title"
          v-model="data.title"
          maxlength="150"
          :class="[error == null ? '' : 'error_input']"
        />
        <div>
          <span v-if="error != null" class="material-icons" :class="$style.materialIcon">warning</span>
          <span class="error_message icon_text">{{ error }}</span>
        </div>
      </div>
      <div :class="$style.button_container">
        <input class="button-primary" type="submit" value="保存" />
        <router-link :to="{ name: 'home' }" class="button">戻る</router-link>
      </div>
    </form>
  </article>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const BoardsRepository = RepositoryFactory.get("boards");

export default {
  name: "AddBoard",
  data: function() {
    return {
      data: [],
      error: null
    };
  },
  methods: {
    async submit() {
      if (
        typeof this.data.title === "undefined" ||
        this.data.title.length < 1
      ) {
        this.error = "タイトルを入力してください。";
        return;
      }
      let payload = {
        title: this.data.title
      };
      await BoardsRepository.addBoard(payload).then(() => {
        this.$router.push({ name: "home" });
      });
    }
  }
};
</script>

<style module>
.container {
  display: grid;
  grid-template-rows: max-content max-content;
  grid-template-columns: 1fr minmax(300px, 600px) 1fr;
  padding: 1rem 0;
}

.item_container {
  display: grid;
  grid-template-rows: max-content max-content max-content;
  grid-template-columns: 100%;
  row-gap: 0.5rem;
  padding: 0 0 1rem 0;
  width: 100%;
  grid-column: 2 / 3;
}

.button_container {
  display: grid;
  grid-template-rows: max-content;
  grid-template-columns: 1fr 1fr;
  column-gap: 0.5rem;
  padding: 1rem 0 0 0;
  width: 100%;
  grid-column: 2 / 3;
}

.small {
  font-size: small;
  font-weight: normal;
}

.materialIcon {
  font-size: 1.5em;
  color: #dc3545;
}
</style>
