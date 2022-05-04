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

    request.post("/api/auth/token", {
        authType : authType,
        authCode: code,
        setCookie: true,
      }).then((response) => {
        let user = getUser();
        ElNotification.success({
          title: 'Login Success',
          message: 'Hi~ ' + user.name,
          showClose: false,
        })
        router.push('/');
      }).catch((error) => {
        let msg = "OAuth failed";
        if (error.response) {
          msg += " " + error.response.data.msg
          console.log("oauth error =>", error.response.data)
        }
        ElMessage({
          message: msg,
          type: "error",
        });
        console.log("oauth failed =>", error);
        router.push('/login');
      });
    }
)

</script>