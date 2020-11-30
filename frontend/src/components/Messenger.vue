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
            <button v-on:click="openChat(dialogue['user']['ID'])">
              {{ dialogue["user"]["Login"] }} ({{
                dialogue["user"]["FullName"]
              }})
            </button>
          </li>
        </ul>
      </div>
      <div class="column">
        <h2>Chat</h2>
        <nav v-if="currentDialogueIndex !== -1">
          <ul>
            <div
              v-for="(message, i) in dialogues[currentDialogueIndex][
                'messages'
              ]"
              :key="`message-${i}`"
            >
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
                  message["Datetime"].substring(12, 19)
                }}</span>
              </div>
            </div>
          </ul>
          <input
            type="text"
            name="message"
            placeholder="Type your message here"
            style="width: 70%; margin-left: 40px"
          />
          <span class="input-group-btn">
            <button
              type="button"
              class="btn btn-primary btn-flat"
              id="message-send"
              style="float: right"
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
              <button v-on:click="openChat(user['ID'])">
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
    };
  },
  methods: {
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
        this.user = res.data;
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
    openChat(val) {
      console.log("openChat(", val, ")");
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
          this.dialogues[dialogueIndex]["messages"] = res.data;
          // this.dialogues.$set(dialogueIndex, dialogue);
          this.currentDialogueIndex = dialogueIndex;
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