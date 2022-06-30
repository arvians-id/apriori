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
                <h3 class="mb-0">Edit Product</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Product Name</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" v-model="product.name" required>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Price</label>
                    <input type="number" class="form-control" v-model="product.price" required>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Description</label>
                    <input type="text" class="form-control" v-model="product.description">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Image</label>
                    <input type="file" class="form-control" @change="uploadImage">
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
  mounted() {
    this.fetchData()
  },
  data: function () {
    return {
      product: {
        name: "",
        price: 0,
        description: "",
        image: null
      }
    };
  },
  methods: {
    submit() {
      const config = {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': 'Bearer ' + localStorage.getItem('token'),
        }
      }
      let formData = new FormData()
      formData.append("name", this.product.name)
      formData.append("price", this.product.price)
      formData.append("description", this.product.description)
      formData.append("image", this.product.image)

      axios.patch(`http://localhost:3000/api/products/${this.$route.params.code}`, formData, config)
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
    },
    fetchData() {
      axios.get(`http://localhost:3000/api/products/${this.$route.params.code}`, { headers: authHeader() }).then(response => {
        this.product = {
          name: response.data.data.name,
          price: response.data.data.price,
          description: response.data.data.description,
        }
      });
    },
    uploadImage(e) {
      this.product.image = e.target.files[0]
    }
  }
}
</script>