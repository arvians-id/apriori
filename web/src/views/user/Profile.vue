<template>
  <!-- Sidenav -->
  <Sidebar />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar />
    <!-- Header -->
    <Header />
    <!-- Page content -->
    <div class="container-fluid mt--6">
      <div class="row">
        <div class="col-xl-4 order-xl-2">
          <div class="card card-profile">
            <div class="card-body pt-0">
              <div class="text-center">
                <h5 class="h3 mt-4">
                  My Name
                </h5>
                <div class="h5 font-weight-300">
                  <i class="ni location_pin mr-2"></i>emailsaya@gmail.com
                </div>
              </div>
            </div>
          </div>
          <!-- Progress track -->
        </div>
        <div class="col-xl-8 order-xl-1">
          <div class="card">
            <div class="card-header">
              <div class="row align-items-center">
                <div class="col-8">
                  <h3 class="mb-0">My Profile </h3>
                </div>
              </div>
            </div>
            <div class="card-body">
              <form @submit.prevent="submit" method="POST">
                <div class="form-group">
                  <label class="form-control-label">Full Name</label> <small class="text-danger">*</small>
                  <input type="text" class="form-control" v-model="user.name">
                </div>
                <div class="form-group">
                  <label class="form-control-label">Email</label> <small class="text-danger">*</small>
                  <input type="email" class="form-control" v-model="user.email">
                </div>
                <div class="form-group">
                  <label class="form-control-label">Password</label> <small class="text-danger">*</small>
                  <input type="password" class="form-control" v-model="user.password">
                </div>
                <button class="btn btn-primary" type="submit">Save</button>
              </form>
            </div>
          </div>
        </div>
      </div>
      <!-- Footer -->
      <Footer />
    </div>
  </div>
</template>

<script>
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import axios from "axios";

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
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
      axios.post("http://localhost:3000/api/userss", this.user)
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'user'
              })
            }
          }).catch(error => {
        alert(error.response.data.status)
      })
    }
  }
}
</script>