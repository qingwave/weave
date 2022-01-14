<template>
  <div class="h-full login-bg">
    <div class="flex h-full justify-center items-center">
      <el-card class="h-max min-w-1/4 w-18rem text-center items-center">
        <div class="inline-flex mt-4 mb-8 items-center">
          <img src="@/assets/weave.png" class="h-12 mr-2" />
          <h1 class="font-bold text-2xl">Weave</h1>
        </div>
        <el-tabs v-model="activeTab" :key="reloadTab">
          <el-tab-pane label="Login" name="login" stretch>
            <el-form ref="loginFormRef" :model="loginUser" :rules="rules">
              <el-form-item prop="name">
                <el-input v-model="loginUser.name" placeholder="Please input">
                  <template #prefix>
                    <User />
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item prop="password">
                <el-input
                  v-model="loginUser.password"
                  type="password"
                  placeholder="Please input password"
                  show-password
                >
                  <template #prefix>
                    <Lock />
                  </template>
                </el-input>
              </el-form-item>
            </el-form>

            <el-button class="w-full" type="primary" @click="login()">Login</el-button>
          </el-tab-pane>
          <el-tab-pane label="Register" name="register" stretch>
            <el-form
              ref="registerFormRef"
              :model="registerUser"
              :rules="rules"
              label-position="left"
              label-width="auto"
            >
              <el-form-item label="Username" prop="name">
                <el-input placeholder="user name" v-model="registerUser.name"></el-input>
              </el-form-item>
              <el-form-item label="Email" prop="email">
                <el-input placeholder="email" v-model="registerUser.email"></el-input>
              </el-form-item>
              <el-form-item label="Password" prop="password">
                <el-input placeholder="password" minlength="6" v-model="registerUser.password"></el-input>
              </el-form-item>
            </el-form>
            <el-button class="w-full" type="primary" @click="register()">Register</el-button>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.login-bg {
  background-image: url('@/assets/login-bg.svg');
  background-repeat: no-repeat;
}
</style>

<script>
import { ElMessage, ElNotification } from "element-plus";
import { User, Lock } from '@icon-park/vue-next';
import { ref, unref, reactive, toRefs } from 'vue';
import request from '@/axios'
import { useRouter } from 'vue-router';

export default {
  name: "Login",
  setup() {
    const router = useRouter();
    const loginFormRef = ref();
    const registerFormRef = ref();
    const activeTab = ref('login');
    const reloadTab = ref(0);
    const state = reactive({
      loginUser: {
        name: "",
        password: "",
      },
      registerUser: {
        name: "",
        email: "",
        password: "",
      },
      rules: {
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
      }
    })

    const login = () => {
      const form = unref(loginFormRef)
      if (!form) {
        return
      }

      form.validate((valid) => {
        if (valid) {
           request.post("/api/auth/token", {
            name: state.loginUser.name,
            password: state.loginUser.password,
            setCookie: true,
          }).then((response) => {
            console.log(response.data.data.msg);
            router.push('/');
            ElNotification.success({
              title: 'Login Success',
              message: 'Hi~ ' + state.loginUser.name,
              showClose: false,
            })
          }).catch((error) => {
            ElMessage({
              message: "Auth failed",
              type: "error",
            });
            console.log("auth failed =>", error);
          });
        } else {
          ElMessage({
            message: "Input invalid",
            type: "error",
          });
        }
      });
    };

    const register = () => {
      const form = unref(registerFormRef)
      if (!form) {
        return
      }

      form.validate((valid) => {
        if (valid) {
          request.post("/api/auth/user", {
            name: state.registerUser.name,
            password: state.registerUser.password,
            email: state.registerUser.email,
          }).then((response) => {
            console.log(response.data.data);
            ElMessage({
              message: 'Register success',
              type: 'success',
            })
            state.loginUser.name = state.registerUser.name;
            state.loginUser.password = state.registerUser.password;
            activeTab.value = 'login';
            reloadTab.value++;
          }).catch((error) => {
            ElMessage({
              message: "Register failed",
              type: "error",
            });
            if (error.response) {
              console.log("register failed =>", error.response.data);
            } else {
              console.log("register failed =>", error)
            }
            
          });
        } else {
          ElMessage({
            message: "Input invalid",
            type: "error",
          });
        }
      });
    }

    return {
      ...toRefs(state),
      activeTab,
      reloadTab,
      loginFormRef,
      registerFormRef,
      login,
      register
    }
  },
  components: {
    User,
    Lock
  }
};
</script>
