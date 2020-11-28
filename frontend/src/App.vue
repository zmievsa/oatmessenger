<template>
  <div id="app">
    <Auth msg="Auth" />
    <Messenger msg="Messenger" />
  </div>
</template>

<script>
import Auth from "./components/Auth.vue";
import Messenger from "./components/Messenger.vue";

export default {
  name: "App",
  components: {
    Auth,
    Messenger,
  },
};
</script>
