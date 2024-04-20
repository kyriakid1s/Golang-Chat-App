<template>
  <div class="w-full overflow-y-hidden flex">
    <div
      class="flex flex-col w-[30%] h-screen p-5 text-lg border-r border-[rgba(0,0,0,1)]"
    >
      Chats:
      <div
        class="flex flex-col cursor-pointer"
        v-for="chat in chats"
        :key="chat.id"
      >
        <div
          class="flex items-center"
          @click="currentConv = chat.Members[0].username"
        >
          <p>{{ chat.Members[0].username }}</p>
          <span
            class="text-xs"
            :class="{
              'text-green-400': chat.Members[0].online == true,
              'text-red-400': !chat.Members[0].online,
            }"
            >●
          </span>
        </div>
      </div>
    </div>
    <Chat
      v-if="currentConv.length != 0"
      :members="currentConv"
      class="w-full overflow-y-hidden"
    />
    <div v-else>Choose a chat</div>
  </div>
</template>

<script>
import axios from "axios";
import { ref } from "vue";
import Chat from "./Chat.vue";
export default {
  components: { Chat },
  setup() {
    let chats = ref([]);
    let currentConv = ref("");

    async function fetchUser() {
      await axios
        .get("http://localhost:4000/v1/message/getchats", {
          withCredentials: true,
        })
        .then((res) => {
          chats.value = res.data.chats;
          console.log(chats);
        })
        .catch((err) => {
          console.log(err);
        });
    }

    return { chats, currentConv, fetchUser };
  },
  created() {
    this.fetchUser();
  },
};
</script>

<style lang="scss" scoped></style>
