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
                <h3 class="mb-0">Create Product</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                 <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Product Code</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" v-model="product.code" required>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Product Name</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" v-model="product.name" required>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Description</label>
                    <input type="text" class="form-control" v-model="product.description">
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
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
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
      product: {
        code: "",
        name: "",
        description: "",
      }
    }
  },
  methods: {
    submit() {
      axios.post("http://localhost:3000/api/products", this.product, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'product'
              })
            }
          }).catch(error => {
        console.log(error.response.data.status)
      })
    }
  }
}
</script>