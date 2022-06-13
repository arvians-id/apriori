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
                <h3 class="mb-0">Buat Produk</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                 <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Kode Product</label>
                    <input type="text" class="form-control" v-model="product.code">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Nama Produk</label>
                    <input type="text" class="form-control" v-model="product.name">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Deskripsi</label>
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
      axios.post("http://localhost:3000/api/products", this.product)
          .then(response => {
            console.log(response)
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'product'
              })
            }
          }).catch(error => {
        alert(error.response.data.status)
      })
    }
  }
}
</script>