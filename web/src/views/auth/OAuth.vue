<template>
</template>

<script setup>
import { onMounted } from 'vue';
import request from '@/axios'
import { getUser } from '@/utils';
import { useRouter } from 'vue-router';

const router = useRouter(); 

onMounted(
    () => {
        const params = new URLSearchParams(window.location.search);
        const state = params.get("state");
        const code = params.get("code");
        const queryString = atob(state);
        const newParams = new URLSearchParams(queryString);
        const authType = newParams.get("oauth");

    request.post("/api/v1/auth/token", {
        authType : authType,
        authCode: code,
        setCookie: true,
      }).then((response) => {
        let user = getUser();
        ElNotification.success({
          title: 'Login Success',
          message: 'Hi~ ' + user.name,
          showClose: false,
          duration: 1500,
        })
        router.push('/');
      })
    }
)

</script>