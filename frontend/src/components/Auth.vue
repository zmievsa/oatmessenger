<template>
  <div class="auth">
    <h1>Auth</h1>
    <div id="authForm">
      <p>
        <input type="text" v-model="usernameModel" placeholder="Username" />
      </p>
      <p>
        <input type="text" v-model="emailModel" placeholder="Email" />
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
      emailModel: "",
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
      if (!this.emailModel.trim()) {
        window.alert("The email field is required.");
        return;
      }
      if (!this.usernameModel.trim()) {
        window.alert("The username field is required.");
        return;
      }
      if (!this.passwordModel.trim()) {
        window.alert("The password field is required.");
        return;
      }
      var data = {
        username: this.usernameModel,
        email: this.emailModel,
        password: this.passwordModel,
      };
      axios
        .post("/register/", data, {
          headers: { "Content-Type": "application/json" },
        })
        .then(() => {
          this.$emit("setcookie");
        })
        .catch((err) => {
          if (err.response.status === 403)
            window.alert("This user already exists.");
        });
    },
    login: function () {
      console.log(
        "Clicked login. Input fields: " +
          this.usernameModel +
          ", " +
          this.passwordModel
      );
      if (!this.usernameModel.trim()) {
        window.alert("The username field is required.");
        return;
      }
      if (!this.passwordModel.trim()) {
        window.alert("The password field is required.");
        return;
      }
      var data = {
        username: this.usernameModel,
        email: this.emailModel,
        password: this.passwordModel,
      };
      axios
        .post("/login/", data, {
          headers: { "Content-Type": "application/json" },
        })
        .then(() => {
          this.$emit("setcookie");
        })
        .catch(() => {
          window.alert(
            "The user/password combination you typed was not found."
          );
        });
    },
  },
};
</script>