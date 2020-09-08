<template>
  <article>
    <h1>新規追加</h1>
    <form @submit.prevent="submit" :class="$style.container">
      <div :class="$style.item_container">
        <label for="y">やったこと</label>
        <textarea id="y" v-model="data.y" rows="4" cols="40"></textarea>
      </div>
      <div :class="$style.item_container">
        <label for="w">わかったこと</label>
        <textarea id="w" v-model="data.w" rows="4" cols="40"></textarea>
      </div>
      <div :class="$style.item_container">
        <label for="t">つぎやること</label>
        <textarea id="t" v-model="data.t" rows="4" cols="40"></textarea>
      </div>
      <div :class="$style.button_container">
        <input class="button-primary" type="submit" value="保存" />
        <router-link :to="{ name: 'board' }" class="button">戻る</router-link>
      </div>
    </form>
  </article>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const BoardsRepository = RepositoryFactory.get("boards");

export default {
  name: "AddItem",
  data: function() {
    return {
      data: [],
    };
  },
  props: {
    id: { type: String, required: true },
  },
  async created() {
    await BoardsRepository.getBoard(`${this.id}`);
  },
  methods: {
    async submit() {
      let payload = {
        y: this.data.y,
        w: this.data.w,
        t: this.data.t,
      };
      await BoardsRepository.addItem(`${this.id}`, payload).then(() => {
        this.$router.push({ name: "board" });
      });
    },
  },
};
</script>

<style module>
.container {
  display: grid;
  grid-template-rows: max-content max-content max-content max-content;
  grid-template-columns: 1fr minmax(300px, 600px) 1fr;
  padding: 1rem 0;
}

.item_container {
  display: grid;
  grid-template-rows: max-content max-content;
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
</style>
