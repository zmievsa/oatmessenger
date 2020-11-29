<template>
  <div id="app">
    <div v-if="cookieExists == false" id="hide">
      <Auth msg="Auth" @setcookie="cookieExists = true" />
    </div>
    <div v-if="cookieExists == true" id="hide">
      <Messenger msg="Messenger" />
    </div>
  </div>
</template>

<script>
import Messenger from "./components/Messenger.vue";
import axios from "axios";
import Auth from "./components/Auth.vue";
// import Vuex from "vuex";
// Vue.use(Vuex);
// import "es6-promise/auto";
export default {
  name: "App",
  components: {
    Auth,
    Messenger,
  },
  data: function () {
    return {
      cookieExists: false,
    };
  },
  methods: {
    checkCookieExistence: function () {
      console.log("Checking cookie existence.");
      axios
        .get("/checkCookieExistence/", { withCredentials: true })
        .then(() => {
          console.log("Cookie found!");
          this.cookieExists = true;
        })
        .catch(() => {
          console.log("No cookie found!");
          this.cookieExists = false;
        });
    },
    showMessenger: function () {},
  },
  watch: {
    cookieExists: function (val, oldval) {
      if (val === true && oldval === false) {
        this.showMessenger(); // TODO Check if this is unnecessary
      }
    },
  },
  created() {
    this.checkCookieExistence();
  },
};
</script>
