<template>
  <article :class="$style.main_container">
    <h1 :class="$style.h1">編集</h1>
    <section :class="$style.base_container">
      <form @submit.prevent="submit" :class="$style.container">
        <h2 :class="$style.item_container">基本情報編集</h2>
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
          <input class="button-primary" type="submit" value="更新" />
        </div>
      </form>
    </section>

    <section :class="$style.members_container">
      <h2>メンバー除外</h2>
      <p>
        <span :class="$style.small">※メンバーを除外しても、そのメンバーの投稿は残ります</span>
      </p>
      <table>
        <thead>
          <tr>
            <th>名前</th>
            <th>ユーザーID</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="val in members" v-bind:key="val.uid">
            <td>{{ val.name }}</td>
            <td>{{ val.uid }}</td>
            <td>
              <button v-if="val.role_description != 'admin'" @click="remove(val.uid)">
                <span class="material-icons">person_remove</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <div :class="$style.back_button_container">
      <router-link :to="{ name: 'home' }" class="button">戻る</router-link>
    </div>
  </article>
</template>

<script>
import { RepositoryFactory } from "./../repositories/RepositoryFactory";
const BoardsRepository = RepositoryFactory.get("boards");

export default {
  name: "EditBoard",
  data: function () {
    return {
      data: [],
      error: null,
      members: [],
    };
  },
  props: {
    id: { type: String, required: true },
  },
  async created() {
    await BoardsRepository.getBoard(`${this.id}`).then((r) => {
      let d = new Object();
      d.title = r.data[0].title;
      d.code = r.data[0].code;
      this.data = d;
    });
    await BoardsRepository.getMembers(`${this.id}`).then((r) => {
      this.members = r.data;
    });
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
        title: this.data.title,
      };
      await BoardsRepository.updateBoard(this.id, payload).then(() => {
        this.$router.push({ name: "home" });
      });
    },
    async remove(uid) {
      let ret = confirm("除外します。よろしいですか？");
      if (ret) {
        await BoardsRepository.deleteBoardMember(this.id, uid);
        await BoardsRepository.getMembers(`${this.id}`).then((r) => {
          this.members = r.data;
        });
      }
    },
  },
};
</script>

<style module>
.main_container {
  display: grid;
  grid-template-rows: max-content max-content max-content max-content;
  grid-template-columns: 1fr minmax(300px, 600px) 1fr;
  padding: 1rem 0;
}

.h1 {
  grid-column: 1 / 3;
}

.base_container {
  grid-column: 2 / 3;
  padding: 1rem 0;
}

.container {
  display: grid;
}

.item_container {
  display: grid;
  padding: 1rem 0;
  width: 100%;
}

.code_container {
  display: grid;
  padding: 1rem 0;
  grid-column: 2 / 3;
  margin-bottom: 2.5rem;
}

.members_container {
  display: grid;
  padding: 1rem 0;
  grid-column: 2 / 3;
}

.button_container {
  display: grid;
  width: 100%;
}

.back_button_container {
  display: grid;
  grid-column: 2 / 3;
  width: 100%;
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
