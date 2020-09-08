<template>
  <article>
    <h1 v-if="data.author_uid == $store.state.auth.uid">編集</h1>
    <h1 v-else>閲覧</h1>
    <form @submit.prevent="submit" :class="$style.container">
      <div :class="$style.item_container">
        <label for="y">やったこと</label>
        <textarea
          v-if="data.author_uid == $store.state.auth.uid"
          id="y"
          v-model="data.y"
          rows="4"
          cols="40"
        ></textarea>
        <p v-else v-html="$sanitize(nl2brData.y)"></p>
      </div>
      <div :class="$style.item_container">
        <label for="w">わかったこと</label>
        <textarea
          v-if="data.author_uid == $store.state.auth.uid"
          id="w"
          v-model="data.w"
          rows="4"
          cols="40"
        ></textarea>
        <p v-else v-html="$sanitize(nl2brData.w)"></p>
      </div>
      <div :class="$style.item_container">
        <label for="t">つぎやること</label>
        <textarea
          v-if="data.author_uid == $store.state.auth.uid"
          id="t"
          v-model="data.t"
          rows="4"
          cols="40"
        ></textarea>
        <p v-else v-html="$sanitize(nl2brData.t)"></p>
      </div>
      <div :class="$style.button_container">
        <input
          v-if="data.author_uid == $store.state.auth.uid"
          class="button-primary"
          type="submit"
          value="保存"
        />
        <router-link :to="{ name: 'board' }" class="button">戻る</router-link>
      </div>
    </form>
  </article>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const BoardsRepository = RepositoryFactory.get("boards");

export default {
  name: "EditItem",
  data: function() {
    return {
      data: []
    };
  },
  computed: {
    nl2brData: function() {
      if (this.data.length < 1) {
        return [];
      }
      return {
        id: this.data.id,
        board_id: this.data.board_id,
        y: this.data.y.replace(/\n/g, "<br/>"),
        t: this.data.t.replace(/\n/g, "<br/>"),
        w: this.data.w.replace(/\n/g, "<br/>"),
        author_uid: this.data.author_uid,
        created: this.data.created,
        updated: this.data.updated
      };
    }
  },
  props: {
    id: { type: String, required: true },
    iid: { type: String, required: true }
  },
  async created() {
    await BoardsRepository.getItem(`${this.id}`, `${this.iid}`).then(
      response => {
        this.data = response.data[0];
      }
    );
  },
  methods: {
    async submit() {
      let payload = {
        id: this.data.id.toString(),
        y: this.data.y,
        w: this.data.w,
        t: this.data.t
      };
      await BoardsRepository.updateItem(
        `${this.id}`,
        `${this.iid}`,
        payload
      ).then(() => {
        this.$router.push({ name: "board" });
      });
    }
  }
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
