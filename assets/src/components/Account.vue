<template>
  <article>
    <h1>アカウント管理</h1>
    <form @submit.prevent="submit" :class="$style.container">
      <div :class="$style.item_container">
        <label for="name">
          名前
          <span :class="$style.small">※必須, 最大150文字</span>
        </label>
        <input
          type="text"
          id="name"
          v-model="data.name"
          maxlength="150"
          :class="[error.name == null ? '' : 'error_input']"
        />
        <div>
          <span
            v-if="error.name != null"
            class="material-icons"
            :class="$style.materialIcon"
            >warning</span
          >
          <span class="error_message icon_text">{{ error.name }}</span>
        </div>
      </div>
      <div :class="$style.item_container">
        <label for="email">
          メールアドレス
          <span :class="$style.small">※必須</span>
        </label>
        <input
          type="email"
          id="email"
          v-model="data.email"
          maxlength="256"
          :class="[error.email == null ? '' : 'error_input']"
        />
        <div>
          <span
            v-if="error.email != null"
            class="material-icons"
            :class="$style.materialIcon"
            >warning</span
          >
          <span class="error_message icon_text">{{ error.email }}</span>
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
const AuthRepository = RepositoryFactory.get("auth");

export default {
  name: "Account",
  data: function () {
    return {
      data: [],
      error: { name: null, email: null },
    };
  },
  async created() {
    this.data.name = this.$store.state.auth.name;
    this.data.email = this.$store.state.auth.email;
  },
  methods: {
    async submit() {
      if (typeof this.data.name === "undefined" || this.data.name.length < 1) {
        this.error.name = "氏名を入力してください。";
      } else {
        this.error.name = null;
      }
      if (
        typeof this.data.email === "undefined" ||
        this.data.email.length < 1
      ) {
        this.error.email = "メールアドレスを入力してください。";
      } else {
        this.error.email = null;
      }
      if (this.error.name != null || this.error.email != null) {
        return;
      }
      let payload = {
        uid: this.$store.state.auth.uid,
        name: this.data.name,
        email: this.data.email,
      };
      await AuthRepository.updateUser(payload).then(() => {
        this.$store.commit("auth/setName", this.data.name);
        this.$store.commit("auth/setEmail", this.data.email);
        this.$router.push({ name: "home" });
      });
    },
  },
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
