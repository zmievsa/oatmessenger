<template>
  <div class="messenger">
    <div v-if="showInterface === true" id="hide">
      <ul id="topbar">
        <h3>
          Welcome, {{ user["Login"] }} (Full name: {{ user["FullName"] }})
        </h3>
        <input
          type="text"
          v-model="newFullName"
          placeholder="Type new full name here"
        />
        <button v-on:click="setFullName()">Change full name</button>
        <button v-on:click="logout()">Log out</button>
      </ul>
      <div class="column">
        <h2>Dialogues</h2>
        <ul id="dialogueList">
          <li v-for="(dialogue, index) in dialogues" :key="`dialogue-${index}`">
            <button v-on:click="openChat('dialogue', dialogue['user']['ID'])">
              {{ dialogue["user"]["Login"] }} ({{
                dialogue["user"]["FullName"]
              }})
            </button>
          </li>
        </ul>
      </div>
      <div class="column" v-if="currentDialogueIndex !== -1">
        <h2>Chat</h2>
        <nav>
          <ul id="messageList">
            <div v-for="(message, i) in currentMessages" :key="`message-${i}`">
              <div
                v-bind:class="
                  message['UserIDFrom'] === user['ID']
                    ? 'container darker'
                    : 'container'
                "
              >
                <b class="right">
                  {{ getUser(message["UserIDFrom"])["Login"] }}
                </b>
                <p>{{ message["Text"] }}</p>
                <span class="time-right">{{
                  message["Datetime"].substring(12, 16)
                }}</span>
              </div>
            </div>
          </ul>
          <span>
            <input
              v-model="chatNewMessage"
              type="text"
              name="message"
              placeholder="Type your message here"
              style="width: 50%; margin-left: 20px"
            />

            <button
              type="button"
              style="width: 75px"
              v-on:click="
                sendMessage(
                  dialogues[currentDialogueIndex]['user']['ID'],
                  chatNewMessage
                )
              "
              id="message-send"
            >
              Send
            </button>
          </span>
        </nav>
      </div>
      <div class="column">
        <div style="float: right">
          <h2>User Search</h2>
          <input
            type="text"
            v-model="searchUsernameField"
            placeholder="Type username here"
          />
          <ul id="userSearchList">
            <li v-for="(user, index) in searchedUsers" :key="`user-${index}`">
              <button v-on:click="openChat('search', user['ID'])">
                {{ user["Login"] }} ({{ user["FullName"] }})
              </button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import axios from "axios";
// import Vue from "vue";
// import VueSocketIO from "vue-socket.io";

// import SocketIO from "socket.io-client";

export default {
  name: "Messenger",
  data: function () {
    return {
      user: {},
      showInterface: false,
      dialogues: [],
      searchUsernameField: "",
      searchedUsers: [],
      newFullName: "",
      currentDialogueIndex: -1,
      chatNewMessage: "",
      socket: null,
      cookie: "",
    };
  },
  computed: {
    currentMessages: function () {
      if (this.currentDialogueIndex === -1) return [];
      else return this.dialogues[this.currentDialogueIndex]["messages"];
    },
  },
  methods: {
    sendMessage(userID, message) {
      message = message.trim();
      if (!message) return;
      this.chatNewMessage = "";
      this.socket.send(
        JSON.stringify({
          userID_for: userID,
          text: message,
          attachments: "",
        })
      );
    },
    getUser(id) {
      console.log("getUser()");
      if (this.user["ID"] === id) return this.user;
      else return this.dialogues[this.currentDialogueIndex]["user"];
    },
    show() {
      console.log("Show messenger");
      this.showInterface = true;
    },
    buildInterface() {
      axios.get("/getUser/").then((res) => {
        this.user = res.data["user"];
        this.cookie = res.data["cookie"];
        console.log(res);
        var conn = new WebSocket(
          "ws://127.0.0.1:8090/ws/?cookie=" + this.cookie
        );
        conn.onclose = function (evt) {
          console.log("CONNECTION CLOSED: ", evt);
        };
        conn.onmessage = (evt) => {
          console.log("NEW SOCKET EVENT: ", evt);
          var data = JSON.parse(evt.data);
          var found = false;
          for (var i = 0; i < this.dialogues.length; i++) {
            var uid = this.dialogues[i]["user"]["ID"];
            if (uid === data["UserIDFrom"] || uid === data["UserIDFor"]) {
              this.dialogues[i]["messages"].push(data);
              this.$nextTick(() => {
                var obj = document.getElementById("messageList");
                obj.scrollTop = obj.scrollHeight;
              });
              found = true;
              break;
            }
          }
          if (!found) {
            axios
              .post("/getAnotherUser/", { userID: data["UserIDFrom"] })
              .then((res) => {
                this.dialogues.push({ user: res.data, messages: [data] });
              })
              .catch((err) => {
                console.log("Error: ", err.response);
              });
          }
        };
        this.socket = conn;
      });
      axios
        .get("/getAllDialogues/")
        .then((res) => {
          for (var i = 0; i < res.data.length; i++) {
            this.dialogues.push({ user: res.data[i], messages: [] });
          }
        })
        .catch((err) => {
          console.log("Error: ", err.response);
        });
    },
    openChat(type, val) {
      console.log("openChat(", val, ")");
      if (
        this.currentDialogueIndex !== -1 &&
        type === "dialogue" &&
        this.dialogues[this.currentDialogueIndex]["user"]["ID"] === val
      ) {
        this.currentDialogueIndex = -1;
        this.chatNewMessage = "";
        return;
      }
      axios
        .post(
          "/getMessages/",
          { user_with_id: val, offset: 0 },
          { headers: { "Content-Type": "application/json" } }
        )
        .then((res) => {
          var dialogueIndex = this.dialogues.findIndex((dialogue) => {
            return dialogue["user"]["ID"] === val;
          });
          if (dialogueIndex === -1) {
            var searchedUsersIndex = this.searchedUsers.findIndex((user) => {
              return user["ID"] === val;
            });
            this.dialogues.push({
              user: this.searchedUsers[searchedUsersIndex],
              messages: res.data,
            });
            dialogueIndex = this.dialogues.length - 1;
          } else this.dialogues[dialogueIndex]["messages"] = res.data;
          this.currentDialogueIndex = dialogueIndex;
          this.chatNewMessage = "";
        })
        .catch((err) => {
          console.log(err);
        });
    },
    logout() {
      console.log("Log out");
      axios.post("/logout/").then(() => {
        this.user = {};
        this.dialogues = [];
        this.searchedUsers = [];
        this.searchUsernameField = "";
        this.showInterface = false;
        this.$emit("setcookie");
        this.socket.close();
      });
    },
    setFullName() {
      axios
        .post(
          "/setFullName/",
          { fullname: this.newFullName },
          { headers: { "Content-Type": "application/json" } }
        )
        .then(() => {
          this.user["FullName"] = this.newFullName;
          this.newFullName = "";
        });
    },
  },
  destroyed: function () {
    this.socket.close();
  },
  watch: {
    showInterface: function (val, oldval) {
      if (val === true && oldval === false) {
        this.buildInterface();
      }
    },
    searchUsernameField: function (val, oldval) {
      val = val.trim();
      if (!val) {
        this.searchedUsers = [];
        return;
      }
      if (val === oldval.trim()) {
        return;
      }
      console.log("findUsers(", val, ")");
      axios
        .post(
          "/findUsers/",
          { name: val },
          { headers: { "Content-Type": "application/json" } }
        )
        .then((res) => {
          console.log("then()");
          console.log("Data received: ", res.data);
          this.searchedUsers = res.data;
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
};
</script>
<style>
nav ul {
  height: 300px;
  width: 100%;
}
nav ul {
  overflow: hidden;
  overflow-y: scroll;
}
.column {
  float: left;
  width: 33%;
}
.rightcolumn {
  float: right;
  width: 33%;
  border: 1px solid black;
}
#topbar {
  list-style-type: none;
  margin-left: 0;
  padding-left: 0;
  margin-bottom: 1%;
}
/* Chat containers */
.container {
  border: 2px solid #dedede;
  background-color: #f1f1f1;
  border-radius: 5px;
  padding: 10px;
  margin: 10px 0;
}

/* Darker chat container */
.darker {
  border-color: #ccc;
  background-color: #ddd;
}

/* Clear floats */
.container::after {
  content: "";
  clear: both;
  display: table;
}

/* Style time text */
.time-right {
  float: right;
  color: #aaa;
}

/* Style time text */
.time-left {
  float: left;
  color: #999;
}
</style>