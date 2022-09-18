<template>
  <!-- Sidenav -->
  <Sidebar :totalNotification="totalNotification" />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar :totalCart="totalCart" :carts="carts" :totalNotification="totalNotification" :notifications="notifications" />
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
                  {{ user.name }}
                </h5>
                <div class="h5 font-weight-300">
                  <i class="ni location_pin mr-2"></i>{{ user.email }}
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
                  <input type="text" class="form-control" name="name" v-model="user.name" required>
                </div>
                <div class="form-group">
                  <label class="form-control-label">Email</label> <small class="text-danger">*</small>
                  <input type="email" class="form-control" name="email" v-model="user.email" required>
                </div>
                <div class="form-group">
                  <label class="form-control-label">Address</label> <small class="text-danger">*</small>
                  <input type="text" class="form-control" name="address" v-model="user.address" required>
                </div>
                <div class="form-group">
                  <label class="form-control-label">Phone Number</label> <small class="text-danger">*</small>
                  <input type="text" class="form-control" name="phone" v-model="user.phone" required>
                </div>
                <div class="form-group">
                  <label class="form-control-label">Password</label> <small class="text-danger">*</small>
                  <input type="password" class="form-control" name="password" v-model="user.password">
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
import Sidebar from "@/components/guest/Sidebar.vue"
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"
import Footer from "@/components/guest/Footer.vue"
import axios from "axios";
import authHeader from "@/service/auth-header";

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  },
  data(){
    return {
      carts: [],
      totalCart: 0,
      user: {
        name: "",
        email: "",
        address: "",
        phone: "",
        password: "",
      },
      totalNotification: 0,
      notifications: []
    }
  },
  mounted() {
    this.fetchData()
    if(authHeader()["Authorization"] !== undefined) {
      this.fetchNotification()
    }
  },
  methods: {
    submit() {
      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/profile/update`, this.user, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'member.profile'
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
          });
    },
    async fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/profile`, { headers: authHeader() })
          .then(response => {
            this.user = {
              name: response.data.data.name,
              email: response.data.data.email,
              address: response.data.data.address,
              phone: response.data.data.phone,
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              console.log(error.response.data.status)
            }
          });
    },
    async fetchNotification() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/notifications/user`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalNotification = response.data.data.filter(e => e.is_read === false).length
          this.notifications = response.data.data
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });
    },
  }
}
</script>