<template>
  <div class="messenger">
    <div v-if="showInterface === true" id="hide">
      <ul id="topbar">
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
    };
  },
  methods: {
    show() {
      console.log("Show messenger");
      this.showInterface = true;
    },
    buildInterface() {
      axios
        .get("/getAllDialogues/")
        .then((res) => {
          this.dialogues = res.data.dialogues;
        })
        .catch((err) => {
          window.alert("Error: ", err.response);
        });
    },
    openChat() {
      console.log("openChat()");
    },
    logout() {
      console.log("Log out");
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