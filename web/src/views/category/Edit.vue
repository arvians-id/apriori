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
                <h3 class="mb-0">Edit Category</h3>
              </div>
              <!-- Card body -->
              <div class="card-body" v-if="isLoading">
                <p class="mt-2 text-center" id="textLoad">Loading...</p>
              </div>
              <div class="card-body" v-else>
                <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Category Name</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" v-model="category.name" required>
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
      category: {
        name: "",
      },
      isLoading: true
    };
  },
  methods: {
    submit() {
      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/categories/${this.$route.params.id}`, this.category, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'category'
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
          })
    },
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/categories/${this.$route.params.id}`, { headers: authHeader() }).then(response => {
        this.category = {
          name: response.data.data.name,
        }
        this.isLoading = false
      }).catch(error => {
        document.getElementById("textLoad").innerHTML = "Failed to load data"
        this.isLoading = true
        console.log(error.response.data.status)
      })
    }
  }
}
</script>