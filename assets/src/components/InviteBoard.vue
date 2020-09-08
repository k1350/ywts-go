<template>
  <article>
    <h1>ボードへの参加</h1>
    <div :class="$style.container">
      <p :class="$style.item_container">{{ title }}</p>
      <div :class="$style.button_container">
        <button @click="join" class="button-primary">参加する</button>
        <router-link :to="{ name: 'home' }" class="button">参加しない</router-link>
      </div>
    </div>
  </article>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const BoardsRepository = RepositoryFactory.get("boards");

export default {
  name: "InviteBoard",
  data: function() {
    return {
      title: null,
      id: null
    };
  },
  props: {
    code: { type: String, required: true }
  },
  async created() {
    const bid = await BoardsRepository.getBoardId(`${this.code}`);
    this.id = bid.data[0].id.toString();

    const m = await BoardsRepository.getMembers(`${this.id}`);
    this.title = m.data[0].title;
    const uids = await m.data.map(d => {
      return d.uid;
    });
    if (uids.indexOf(this.$store.state.auth.uid) > -1) {
      this.$router.push({ name: "board", params: { id: this.id } });
    }
  },
  methods: {
    async join() {
      let payload = {
        code: this.code
      };
      await BoardsRepository.joinBoard(`${this.id}`, payload).then(() => {
        this.$router.push({ name: "board", params: { id: this.id } });
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
  font-weight: bold;
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
