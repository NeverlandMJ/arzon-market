<template>
  <div class="home">
    <form @submit.prevent="register">
      <input v-model="full_name" type="text" placeholder="FI"/>
      <input v-model="email" type="email" placeholder="email"/>
      <input v-model="password" type="password" placeholder="password"/>
      <input v-model="card_number" type="text" placeholder="cardNumber"/>
      <input v-model='balance' type="text" placeholder="balance"/>
      <button type="submit">Submit</button>
    </form>
   
  </div>
</template>

<script setup>
import { ref } from 'vue';

const full_name = ref('')
const email = ref('')
const password = ref('')
const card_number = ref('')
const balance = ref(null)

const register = ()=>{


}

</script>
