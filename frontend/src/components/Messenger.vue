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
      <div id="box1">
        <h2>Dialogues</h2>
        <ul id="dialogueList">
          <li v-for="(dialogue, index) in dialogues" :key="`dialogue-${index}`">
            <button
              :userID="dialogue['user']['ID']"
              v-on:click="openChat(userID)"
            >
              {{ dialogue["user"]["Login"] }}
            </button>
          </li>
        </ul>
      </div>
      <div id="box2">
        <h2>User Search</h2>
        <input
          type="text"
          v-model="searchUsernameField"
          placeholder="Type username here"
        />
        <ul id="userSearchList">
          <li v-for="(user, index) in searchedUsers" :key="`user-${index}`">
            <button :userID="user['ID']" v-on:click="openChat(userID)">
              {{ user["Login"] }}
            </button>
          </li>
        </ul>
      </div>
      <div id="box3">
        <h2>Chat</h2>
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
    };
  },
  methods: {
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
    openChat() {
      console.log("openChat()");
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
#box1 {
  float: left;
  width: 33%;
  height: 280px;
}

#box2 {
  float: right;
  width: 33%;
  height: 280px;
}
#box3 {
  float: right;
  width: 33%;
  height: 280px;
}
#topbar {
  list-style-type: none;
  margin: 0;
  padding: 0;
}
</style>