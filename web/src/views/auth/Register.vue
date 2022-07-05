<template>
  <!-- Navbar -->
  <Navbar />
  <div class="main-content">
    <!-- Header -->
    <Header />
    <!-- Page content -->
    <div class="container mt--8 pb-5">
      <div class="row justify-content-center">
        <div class="col-lg-5 col-md-7">
          <div class="card bg-secondary border-0 mb-0">
            <div class="card-body px-lg-5 py-lg-5">
              <div class="text-center text-muted mb-4">
                <small>Sign up your account</small>
              </div>
              <form @submit.prevent="submit" method="POST" role="form">
                <div class="form-group mb-3">
                  <div class="input-group input-group-merge input-group-alternative">
                    <div class="input-group-prepend">
                      <span class="input-group-text"><i class="ni ni-single-02"></i></span>
                    </div>
                    <input class="form-control" placeholder="Full Name" type="text" v-model="user.name" required>
                  </div>
                </div>
                <div class="form-group mb-3">
                  <div class="input-group input-group-merge input-group-alternative">
                    <div class="input-group-prepend">
                      <span class="input-group-text"><i class="ni ni-email-83"></i></span>
                    </div>
                    <input class="form-control" placeholder="Email" type="email" v-model="user.email" required>
                  </div>
                </div>
                <div class="form-group">
                  <div class="input-group input-group-merge input-group-alternative">
                    <div class="input-group-prepend">
                      <span class="input-group-text"><i class="ni ni-lock-circle-open"></i></span>
                    </div>
                    <input class="form-control" placeholder="Password" type="password" v-model="user.password" required>
                  </div>
                </div>
                <div class="text-center">
                  <button type="submit" class="btn btn-primary my-4">Register</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <!-- Footer -->
  <Footer />
</template>

<script>
import Footer from "@/components/auth/Footer.vue"
import Navbar from "@/components/auth/Navbar.vue"
import Header from "@/components/auth/Header.vue"
import axios from "axios";

export default {
  components: {
    Footer,
    Navbar,
    Header
  },
  data(){
    return {
      user: {
        name: "",
        email: "",
        password: "",
      }
    }
  },
  methods: {
    submit() {
      axios.post(`${process.env.VUE_APP_SERVICE_URL}/auth/register`, this.user)
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'auth.login'
              })
            }
          }).catch(error => {
            alert(error.response.data.status)
      })
    }
  }
}
</script>