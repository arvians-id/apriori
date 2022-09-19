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
                     <label class="form-control-label">Price</label> <small class="text-danger">*</small>
                     <input type="number" class="form-control" v-model="product.price" required>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Category Name</label> <small class="text-danger">*use ctrl for selecting the category</small>
                     <select class="form-control" v-model="product.category" multiple required>
                        <option v-for="(category, i) in categories" :value="category.name" :key="i">{{ category.name }}</option>
                     </select>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Mass (gram)</label> <small class="text-danger">*</small>
                     <input type="number" class="form-control" v-model="product.mass" required>
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
  data(){
    return {
      product: {
        code: "",
        name: "",
        price: 0,
        category: [],
        mass: 0,
        description: "",
        image: null
      },
      categories: [],
      previewImage: "https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/no-image.png"
    }
  },
  mounted() {
    this.fetchCategories();
  },
  methods: {
    async fetchCategories(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/categories`, { headers: authHeader() })
      .then(response => {
        this.categories = response.data.data;
      })
      .catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });
    },
    submit() {
      const config = {
        headers: {
          'Content-Type': 'multipart/form-data',
          ...authHeader(),
        }
      }
      let formData = new FormData()
      formData.append("code", this.product.code)
      formData.append("name", this.product.name)
      formData.append("price", this.product.price)
      formData.append("description", this.product.description)
      formData.append("image", this.product.image)
      formData.append("mass", this.product.mass)
      if (this.product.category.length > 0) {
        let categoryName = this.product.category
        this.product.category = categoryName.join(", ")
        formData.append("category", this.product.category)
      }

      axios.post(`${process.env.VUE_APP_SERVICE_URL}/products`, formData, config)
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'product'
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
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