<template>
  <div class="h-full login-bg bg-slate-50">
    <div class="flex h-full justify-center items-center">
      <div class="h-max min-w-[16rem] w-1/4 max-w-[24rem] text-center items-center">
        <div class="inline-flex mt-4 mb-8 items-center">
          <img src="@/assets/weave.png" class="h-12 mr-2" />
          <h1 class="font-bold text-4xl font-mono">Weave</h1>
        </div>

        <div v-if="showLogin">
            <el-form ref="loginFormRef" :model="loginUser" size="large" :rules="rules" show-message>
              <el-form-item prop="name">
                <el-input v-model="loginUser.name" placeholder="admin">
                  <template #prefix>
                    <User />
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item prop="password">
                <el-input v-model="loginUser.password" type="password" placeholder="123456" show-password>
                  <template #prefix>
                    <Lock />
                  </template>
                </el-input>
              </el-form-item>
            </el-form>

            <el-button class="w-full" type="primary" size="large" @click="login(loginFormRef)">SIGN IN</el-button>
            <div class="w-full flex mt-[0.25rem]">
              <div class="w-full text-right">
                <el-button link @click="showLogin=false">SIGN UP</el-button>
              </div>
            </div>
            
            <div class="my-[0.5rem]">
              <el-button link @click="oauthLogin('github')">
                <Github theme="outline" size="30" fill="#333" />
              </el-button>
              <el-button link @click="oauthLogin('wechat')">
                <Wechat theme="filled" size="30" fill="#7ed321" />
              </el-button>
            </div>
        </div>

          <div v-if="showLogin == false">
            <el-form ref="registerFormRef" :model="registerUser" label-position="top" :rules="rules"
              label-width="auto" size="large">
              <el-form-item label="Username" prop="name">
                <el-input placeholder="user name" v-model="registerUser.name" size="large"></el-input>
              </el-form-item>
              <el-form-item label="Email" prop="email">
                <el-input placeholder="email" v-model="registerUser.email"></el-input>
              </el-form-item>
              <el-form-item label="Password" prop="password">
                <el-input placeholder="password" minlength="6" v-model="registerUser.password"></el-input>
              </el-form-item>
            </el-form>
            
            <el-button class="w-full" type="primary" size="large" @click="register(registerFormRef)">SIGN UP</el-button>
            <div class="mt-[0.25rem] text-right">
              <el-button link @click="showLogin=true">SIGN IN</el-button>
            </div>
          </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-bg {
  background-image: url('@/assets/login-bg.svg');
  background-repeat: no-repeat;
  background-size: 100% auto;
  background-position: 0 100%;
}
</style>

<script setup>
import { ElMessage, ElNotification } from "element-plus"
import { User, Lock, Github, Wechat } from '@icon-park/vue-next'
import { ref, reactive } from 'vue'
import request from '@/axios'
import { useRouter } from 'vue-router'
import { authInfo } from '@/config.js'

const router = useRouter();

const loginFormRef = ref();
const registerFormRef = ref();
const redirectUri = window.location.origin + '/oauth'

const showLogin = ref(true);

const loginUser = reactive({
  name: "admin",
  password: "123456",
});
const registerUser = reactive({
  name: "",
  email: "",
  password: "",
});
const rules = reactive({
  name: [
    { required: true, message: 'Please input user name', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please input password', trigger: 'blur' },
    { min: 6, message: 'Length should be great than 6', trigger: 'blur' }
  ],
  email: [
    { required: true, message: 'Please input email', trigger: 'blur' },
    { type: 'email', message: 'Please input correct email address', trigger: ['blur', 'change'] },
  ]
});

const login = async (form) => {
  if (!form) {
    return
  }

  let name = loginUser.name;

  let success = function() {
    ElNotification.success({
          title: 'Login Success',
          message: 'Hi~ ' + name,
          showClose: true,
          duration: 1500,
        })
    router.push('/');
  }

  await form.validate((valid, fields) => {
    if (valid) {
      request.post("/api/v1/auth/token", {
        name: loginUser.name,
        password: loginUser.password,
        setCookie: true,
      }).then((response) => {
        success()
      })
    } else {
      console.log("input invaild", fields)
      ElMessage({
        message: "Input invalid" + fields,
        type: "error",
      });
    }
  });
};

const oauthLogin = (authType) => {
  if (!authInfo[authType]) {
    return
  }

  let uri = "";
  const endpoint = authInfo[authType].endpoint;
  const scope = authInfo[authType].scope;
  const clientId = authInfo[authType].clientId;
  const state = btoa(`${window.location.search}&app=weave&oauth=${authType}`)

  if (authType === "google") {
    uri = `${endpoint}?client_id=${clientId}&redirect_uri=${redirectUri}&scope=${scope}&response_type=code&state=${state}`;
  } else if (authType === "github") {
    uri = `${endpoint}?client_id=${clientId}&redirect_uri=${redirectUri}&scope=${scope}&response_type=code&state=${state}`;
  } else if (authType === "wechat") {
    if (navigator.userAgent.includes("MicroMessenger")) {
      uri = `${authInfo[authType].mpEndpoint}?appid=${authInfo[authType].clientId2}&redirect_uri=${redirectUri}&state=${state}&scope=${authInfo[authType].mpScope}&response_type=code#wechat_redirect`;
    } else {
      uri = `${endpoint}?appid=${clientId}&redirect_uri=${redirectUri}&scope=${scope}&response_type=code&state=${state}#wechat_redirect`;
    }
  } else {
    console.log(`auth type ${authType} not supported`)
    return
  }
  window.location.href = uri;
};

const register = async (form) => {
  if (!form) {
    return
  }

  await form.validate((valid, fields) => {
    if (valid) {
      request.post("/api/v1/auth/user", {
        name: registerUser.name,
        password: registerUser.password,
        email: registerUser.email,
      }).then((response) => {
        ElMessage({
          message: 'Register success',
          type: 'success',
        })
        loginUser.name = registerUser.name;
        loginUser.password = registerUser.password;
        activeTab.value = 'login';
      })
    } else {
      console.log("Input invalid =>", fields)
      ElMessage({
        message: "Input invalid",
        type: "error",
      });
    }
  });
};
</script>
