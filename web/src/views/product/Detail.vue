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
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile">
            <!-- Card header -->
            <div class="card-header">
              <h3 class="mb-0">Detail Product</h3>
            </div>
            <div class="row align-items-center">
              <div class="col-12 col-lg-6 text-center">
                <img :src="getImage()" class="img-fluid img-thumbnail my-5" width="500">
              </div>
              <div class="col-12 col-lg-6">
                <div class="card-header text-center border-0 pt-8 pt-md-4 pb-0 pb-md-4">
                  <div class="d-flex justify-content-between">
                    <router-link :to="{ name: 'product.edit', params: { code: this.$route.params.code } }"  class="btn btn-sm btn-default float-right">Edit</router-link>
                    <form @submit.prevent="submit(this.$route.params.code)" method="POST" class="d-inline">
                      <button class="btn btn-danger btn-sm">Delete</button>
                    </form>
                  </div>
                </div>
                <div class="card-body pt-0">
                  <div class="row">
                    <div class="col">
                      <div class="card-profile-stats d-flex justify-content-center">
                        <div>
                          <span class="heading">Created At</span>
                          <span class="description">{{ product.created_at }}</span>
                        </div>
                        <div>
                          <span class="heading">Last Modified</span>
                          <span class="description">{{ product.updated_at }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="text-center">
                    <h5 class="h3 text-uppercase">
                      {{ product.name }} - {{ product.code }}
                    </h5>
                    <div class="h5 font-weight-300">
                      <i class="ni location_pin mr-2"></i>Rp. {{ product.price }}
                    </div>
                    <div>
                      <i class="ni education_hat mr-2"></i> {{ product.description }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- Progress track -->
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
      product: [],
    };
  },
  methods: {
    fetchData() {
      axios.get(`http://localhost:3000/api/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      });
    },
    getImage() {
      return this.product.image
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    },
    submit(no_product) {
      axios.delete("http://localhost:3000/api/products/" + no_product, { headers: authHeader() })
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