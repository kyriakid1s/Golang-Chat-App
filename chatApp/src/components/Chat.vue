<template>
  <div class="flex flex-col overflow-y-hidden">
    <header
      class="w-full h-auto flex items-center p-5 text-lg border-b border-[rgba(0,0,0,0.2)]"
    >
      {{ members }}
    </header>
    <main class="overflow-y-auto" style="max-height: calc(100vh - 200px)">
      <div v-for="message in messages" :key="message.id">
        <p
          :class="{
            'ml-auto bg-blue-200': message.to == members,
            'bg-gray-200': message.to != members,
          }"
          class="p-2 m-1 rounded-xl w-fit"
        >
          {{ message.message }}
        </p>
      </div>
      <div ref="bottom"></div>
    </main>

    <div
      class="w-full h-auto m-auto p-4 border border-[rgba(0,0,0,0,2)] rounded-r-full rounded-l-full"
    >
      <input
        @keyup.enter="sendMessage"
        type="text"
        placeholder="Message..."
        v-model="newMessage"
        class="w-full bg-transparent outline-none"
      />
    </div>
  </div>
</template>

<script setup>
import axios from "axios";
import { ref, toRefs, onMounted, onUpdated, onUnmounted } from "vue";
import socket from "@/websocket";
const props = defineProps({
  members: String,
});
const { members } = toRefs(props);
let bottom = ref();

let messages = ref([]);
let newMessage = ref("");
async function getMessages(member) {
  console.log(member);
  messages.value = [];
  await axios
    .get(`http://localhost:4000/v1/message/get/${member}`, {
      withCredentials: true,
    })
    .then((res) => {
      console.log(res);
      messages.value = res.data.messages;
      console.log(messages);
    })
    .catch((err) => {
      console.log(err);
    });
}

async function sendMessage() {
  await axios
    .post(
      "http://localhost:4000/v1/message/send",
      { to: members.value, message: newMessage.value },
      { withCredentials: true }
    )
    .then(() => {
      messages.value.push({ to: members.value, message: newMessage.value });
      let data = { recipient: members.value, content: newMessage.value };
      socket.send(JSON.stringify(data));
      newMessage.value = "";
    })
    .catch((err) => {
      console.log(err);
    });
}
function scrollToLast() {
  bottom.value.scrollIntoView();
}

onMounted(async () => {
  await getMessages(members.value);
  socket.onopen = () => {};
  socket.onmessage = function (evt) {
    let data = JSON.parse(evt.data);
    messages.value.push({ to: data.recipient, message: data.content });
  };
});
onUpdated(() => {
  scrollToLast();
});
</script>
