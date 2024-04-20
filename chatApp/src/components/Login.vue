<template>
  <div class="w-screen h-screen flex items-center justify-center">
    <form
      action="#"
      class="w-auto flex flex-col items-center justify-center py-5 px-10 border-2 rounded-xl"
      @submit.prevent="onSubmit"
    >
      <label for="username">Username: </label>
      <span class="border-2 p-">
        <input
          class="bg-transparent"
          name="username"
          type="text"
          v-model="username"
        />
      </span>

      <label for="password"> Password:</label>
      <span class="border-2 p-">
        <input
          class="bg-transparent"
          name="password"
          type="password"
          v-model="password"
        />
      </span>
      <button type="submit">Log in</button>
    </form>
  </div>
</template>

<script>
import { useRouter } from "vue-router";
import axios from "axios";
import { ref } from "vue";

export default {
  setup() {
    let router = useRouter();
    let username = ref("");
    let password = ref("");

    async function onSubmit() {
      await axios
        .post(
          "http://localhost:4000/v1/user/login",
          {
            username: this.username,
            password: this.password,
          },
          {
            withCredentials: true,
          }
        )
        .then((res) => {
          console.log(res);
          router.push("/chats");
        })
        .catch((err) => {
          console.log(err);
        });
    }

    return { onSubmit, username, password };
  },
};
</script>
