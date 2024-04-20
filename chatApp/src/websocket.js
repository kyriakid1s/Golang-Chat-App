import axios from "axios";
let user = "";
await axios
  .get("http://localhost:4000/v1/ws/getuser", {
    withCredentials: true,
  })
  .then((res) => {
    user = res.data.user;
  })
  .catch((err) => {
    console.log(err);
  });
const url = `ws://localhost:4000/ws?user=${user}`;
const socket = new WebSocket(url);

export default socket;
