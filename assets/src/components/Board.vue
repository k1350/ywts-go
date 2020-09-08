<template>
  <article class="width: 100%;">
    <h1>{{ board.title }} のYWT一覧</h1>
    <section :class="$style.container" ref="cards">
      <div class="add-card" :class="$style.add" @click="addItem">
        <span class="material-icons" :class="[$style.addBox]">add_box</span>
      </div>
      <div
        class="card"
        v-for="val in nl2brData"
        v-bind:key="val.id"
        v-show="data != null && data.length > 0"
      >
        <div :class="$style.cardContent">
          <p :class="$style.title">やったこと</p>
          <p v-html="$sanitize(val.y)"></p>
          <p :class="$style.title">わかったこと</p>
          <p v-html="$sanitize(val.w)"></p>
          <p :class="$style.title">つぎやること</p>
          <p v-html="$sanitize(val.t)"></p>
        </div>
        <div>
          <p>{{ val.author }}</p>
          <div :class="$style.authorContainer" v-if="val.author_uid == $store.state.auth.uid">
            <router-link
              :to="{
                name: 'edit',
                params: { id: id.toString(), iid: val.id.toString() },
              }"
              class="button"
            >
              <span class="material-icons">edit</span>
            </router-link>
            <button @click="deleteItem(val.id)">
              <span class="material-icons">delete_forever</span>
            </button>
          </div>
        </div>
      </div>
    </section>
    <div :class="$style.back">
      <router-link :to="{ name: 'home' }" class="button">戻る</router-link>
    </div>
  </article>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const BoardsRepository = RepositoryFactory.get("boards");
const moment = require("moment");

function changeGridRow(c) {
  for (var i = 1; i < c.length; i++) {
    var boxInnerH = c[i].scrollHeight;
    c[i].style.gridRow = "span " + Math.ceil(boxInnerH / 100.0);
  }
}

export default {
  name: "Board",
  data: function () {
    return {
      data: [],
      board: "",
    };
  },
  computed: {
    nl2brData: function () {
      return this.data.map((item) => ({
        id: item.id,
        board_id: item.board_id,
        y: item.y.replace(/\n/g, "<br/>"),
        t: item.t.replace(/\n/g, "<br/>"),
        w: item.w.replace(/\n/g, "<br/>"),
        author_uid: item.author_uid,
        author: item.author,
        created: moment(moment.utc(item.created))
          .local()
          .format("YYYY-MM-DD HH:mm:ss"),
        updated: moment(moment.utc(item.updated))
          .local()
          .format("YYYY-MM-DD HH:mm:ss"),
      }));
    },
  },
  props: {
    id: { type: String, required: true },
  },
  async created() {
    await BoardsRepository.getBoard(`${this.id}`).then((r) => {
      this.board = r.data[0];
    });
    await BoardsRepository.getItems(`${this.id}`)
      .then((r) => {
        this.data = r.data;
      })
      .then(() => {
        changeGridRow(this.$refs.cards.children);
      });
  },
  methods: {
    addItem() {
      this.$router.push({ name: "add", params: { id: this.id.toString() } });
    },
    async deleteItem(iid) {
      let ret = confirm("削除します。よろしいですか？");
      if (ret) {
        await BoardsRepository.deleteItem(`${this.id}`, `${iid}`);
        await BoardsRepository.getItems(`${this.id}`)
          .then((r) => {
            this.data = r.data;
          })
          .then(() => {
            changeGridRow(this.$refs.cards.children);
          });
      }
    },
  },
};
</script>

<style module>
.container {
  display: grid;
  grid-template-rows: 100px;
  grid-template-columns: repeat(auto-fit, 300px);
  grid-auto-flow: dense;
  column-gap: 2rem;
  row-gap: 2rem;
  justify-content: center;
}

.back {
  margin-top: 1rem;
  text-align: center;
}

.title {
  font-weight: bold;
}

.add {
  text-align: center;
}

.authorContainer {
  display: grid;
  grid-template-rows: max-content;
  grid-template-columns: 45% 45%;
  column-gap: 10%;
}

.addBox {
  font-size: 8rem;
}
</style>
