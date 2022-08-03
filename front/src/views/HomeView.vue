<template>
  <div class="home">
    <form @submit.prevent="register">
      <input v-model="full_name" type="text" placeholder="FI" />
      <input v-model="email" type="email" placeholder="email" />
      <input v-model="password" type="password" placeholder="password" />
      <button type="submit">Submit</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import http from "@/plugins/http";

const router = useRouter();
const full_name = ref("");
const email = ref("");
const password = ref("");

const register = () => {
  http
    .post("/users/post_register", {
      email: email.value,
      full_name: full_name.value,
      password: password,
    })
    .then((res) => {
      if (res.data.success) {
        router.push("profile");
        console.log(res.data);
        console.log("Uddaladik");
      }
    })
    .catch((error) => {
      console.log(error);
    });
};
</script>
