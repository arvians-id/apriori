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
                     <label class="form-control-label">Price</label>
                     <input type="number" class="form-control" v-model="product.price" required>
                   </div>
                  <div class="form-group">
                    <label class="form-control-label">Description</label>
                    <textarea class="form-control" v-model="product.description" rows="5"></textarea>
                  </div>
                   <div class="form-group">
                     <label class="form-control-label">Image</label>
                     <input type="file" class="form-control" @change="uploadImage">
                   </div>
                   <div class="form-group">
                     <img :src="previewImage" width="150"/>
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
        price: 0,
        description: "",
        image: null
      },
      previewImage: "https://my-apriori.s3.ap-southeast-1.amazonaws.com/no-image.png"
    }
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
      formData.append("code", this.product.code)
      formData.append("name", this.product.name)
      formData.append("price", this.product.price)
      formData.append("description", this.product.description)
      formData.append("image", this.product.image)

      axios.post(`${process.env.VUE_APP_SERVICE_URL}/products`, formData, config)
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
    uploadImage(e) {
      let files = e.target.files[0]
      this.product.image = files
      this.previewImage = URL.createObjectURL(files)
    }
  }
}
</script>