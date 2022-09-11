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
        <div class="col-12">
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card header -->
              <div class="card-header">
                <h3 class="mb-0">Edit User</h3>
              </div>
              <!-- Card body -->
              <div class="card-body" v-if="isLoading">
                <p class="mt-2 text-center">Loading...</p>
              </div>
              <div class="card-body" v-else>
                <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Full Name</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" name="name" v-model="user.name" required>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Role</label> <small class="text-danger">*</small>
                    <input type="number" min="1" max="2" class="form-control" v-model="user.role" required>
                    <span class="text-danger text-sm">*1 Admin and *2 Member</span>
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
                    <label class="form-control-label">Password</label>
                    <input type="password" class="form-control" name="password" v-model="user.password">
                  </div>
                  <button class="btn btn-primary" type="submit">Submit form</button>
                </form>
              </div>
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
import Sidebar from "@/components/admin/Sidebar.vue"
import Topbar from "@/components/admin/Topbar.vue"
import Header from "@/components/admin/Header.vue"
import Footer from "@/components/admin/Footer.vue"
import axios from "axios";
import authHeader from "@/service/auth-header";

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  },
  mounted() {
    this.fetchData()
  },
  data: function () {
    return {
      user: {
        role: 2,
        name: "",
        email: "",
        address: "",
        phone: "",
        password: "",
      },
      isLoading: true
    };
  },
  methods: {
    submit() {
      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/users/${this.$route.params.id}`, this.user, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'user'
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
          })
    },
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/users/${this.$route.params.id}`, { headers: authHeader() }).then(response => {
        this.user = {
          role: response.data.data.role,
          name: response.data.data.name,
          email: response.data.data.email,
          address: response.data.data.address,
          phone: response.data.data.phone,
        };
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading = false;
    }
  }
}
</script>