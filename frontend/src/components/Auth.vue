<template>
  <div class="auth">
    <h1>Auth</h1>
    <div id="authForm">
      <p>
        <input type="text" v-model="usernameModel" placeholder="Username" />
      </p>
      <p>
        <input type="password" v-model="passwordModel" placeholder="Password" />
      </p>
      <!-- Password again: -->
      <button variant="primary" v-on:click="register()">Register</button>
      <button variant="primary" v-on:click="login()">Login</button>
    </div>
  </div>
</template>
<script>
import axios from "axios";

export default {
  name: "Auth",
  data: function () {
    return {
      usernameModel: "",
      passwordModel: "",
    };
  },
  methods: {
    register: function () {
      console.log(
        "Clicked register. Input fields: " +
          this.usernameModel +
          ", " +
          this.passwordModel
      );
      var data = {
        username: this.usernameModel,
        password: this.passwordModel,
      };
      axios
        .post("/register/", data, {
          headers: { "Content-Type": "application/json" },
        })
        .then((res) => {
          console.log(res.headers);
          this.$emit("setcookie");
        })
        .catch((err) => {
          console.log(err.response);
        });
    },
  },
};
</script>