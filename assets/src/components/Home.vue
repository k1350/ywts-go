<template>
  <article>
    <h1>参加ボード一覧</h1>
    <section :class="$style.container" ref="cards">
      <div class="add-card" :class="$style.add" @click="addBoard">
        <span class="material-icons" :class="[$style.addBox]">add_box</span>
      </div>
      <div
        class="card"
        v-for="val in timeData"
        v-bind:key="val.id"
        v-show="data != null && data.length > 0"
      >
        <div :class="$style.cardContent">
          <p :class="$style.title">
            <router-link :to="{ name: 'board', params: { id: val.id.toString() } }">{{ val.title }}</router-link>
          </p>
          <p :class="$style.title">作成日</p>
          <p>{{ val.created }}</p>
        </div>
        <div>
          <div :class="$style.authorContainer" v-if="val.is_admin == 1">
            <router-link
              :to="{
                name: 'editBoard',
                params: { id: val.id.toString() },
              }"
              class="button"
            >
              <span class="material-icons">edit</span>
            </router-link>
            <button @click="deleteBoard(val.id)">
              <span class="material-icons">delete_forever</span>
            </button>
          </div>
          <div :class="$style.authorContainer" v-else>
            <button @click="exitBoard(val.id)">
              <span class="material-icons">exit_to_app</span>
            </button>
          </div>
        </div>
        <div v-if="val.is_admin == 1">参加URL</div>
        <div :class="$style.joinContainer" v-if="val.is_admin == 1">
          <input type="url" :id="'join' + val.id" :value="val.code" readonly />
          <button @click="copyToClipboard('join' + val.id)">
            <span class="material-icons">content_copy</span>
          </button>
        </div>
      </div>
    </section>
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
  name: "Home",
  data: function () {
    return {
      data: [],
    };
  },
  async created() {
    await BoardsRepository.get()
      .then((response) => {
        this.data = response.data;
      })
      .then(() => {
        changeGridRow(this.$refs.cards.children);
      });
  },
  computed: {
    timeData: function () {
      return this.data.map((item) => ({
        id: item.id,
        title: item.title,
        created: moment(moment.utc(item.created))
          .local()
          .format("YYYY-MM-DD HH:mm:ss"),
        updated: moment(moment.utc(item.updated))
          .local()
          .format("YYYY-MM-DD HH:mm:ss"),
        is_admin: item.is_admin,
        code: window.location.origin + "/invite-board/" + item.code,
      }));
    },
  },
  methods: {
    addBoard() {
      this.$router.push({ name: "addBoard" });
    },
    async deleteBoard(id) {
      let ret = confirm("削除します。よろしいですか？");
      if (ret) {
        await BoardsRepository.deleteBoard(`${id}`);
        await BoardsRepository.get()
          .then((response) => {
            this.data = response.data;
          })
          .then(() => {
            changeGridRow(this.$refs.cards.children);
          });
      }
    },
    copyToClipboard(id) {
      var copyTarget = document.getElementById(id);
      copyTarget.select();
      document.execCommand("Copy");
      alert("参加URLをクリップボードにコピーしました : " + copyTarget.value);
    },
    async exitBoard(id) {
      let ret = confirm("ボードから退出します。よろしいですか？");
      if (ret) {
        await BoardsRepository.deleteBoardMember(`${id}`, this.$store.state.auth.uid);
        await BoardsRepository.get()
          .then((response) => {
            this.data = response.data;
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

.title {
  font-weight: bold;
}

.add {
  text-align: center;
}

.addBox {
  font-size: 8rem;
}

.authorContainer {
  display: grid;
  grid-template-rows: max-content;
  grid-template-columns: 45% 45%;
  column-gap: 10%;
}

.joinContainer {
  display: grid;
  grid-template-rows: max-content;
  grid-template-columns: 60% 35%;
  column-gap: 5%;
}
</style>
