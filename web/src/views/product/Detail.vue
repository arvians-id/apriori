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
            <div class="card-body">
              <div class="row align-items-center mx-auto" v-if="isLoading">
                <p class="p-3 mt-2 text-center">Loading...</p>
              </div>
              <div class="row d-flex justify-content-center" v-else>
                <div class="col-12 col-lg-4 text-center mb-2">
                  <img :src="getImage()" class="img-fluid mb-2" width="500">
                  <router-link :to="{ name: 'product.edit', params: { code: this.$route.params.code } }"  class="btn btn-primary btn-block mb-2">Edit</router-link>
                  <form @submit.prevent="submit(this.$route.params.code)" method="POST" class="d-inline">
                    <button class="btn btn-danger btn-block">Delete</button>
                  </form>
                </div>
                <div class="col-12 col-lg-4">
                  <div class="text-left">
                    <h5 class="h2 text-uppercase p-0 m-0">{{ product.name }}</h5>
                    <p class="p-0 m-0">code : {{ product.code }}</p>
                    <div class="h1 font-weight-bold">Rp. {{ product.price }}</div>
                    <hr class="m-2">
                    <ul class="nav nav-tabs" id="myTab" role="tablist">
                      <li class="nav-item">
                        <a class="nav-link active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true">Detail</a>
                      </li>
                      <li class="nav-item">
                        <a class="nav-link" id="other-detail-tab" data-toggle="tab" href="#other-detail" role="tab" aria-controls="other-detail" aria-selected="false">Lainnya</a>
                      </li>
                    </ul>
                    <div class="tab-content" id="myTabContent">
                      <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                        <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                        <p>Original baru</p>
                        <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                        <p>{{ product.category }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                        <p>{{ product.mass }} gram</p>
                        <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                        <p>{{ product.description }}</p>
                      </div>
                      <div class="tab-pane fade" id="other-detail" role="tabpanel" aria-labelledby="other-detail-tab">
                        <p class="font-weight-bold mb-0 mt-2">Tanggal Dibuat</p>
                        <p>{{ product.created_at }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Terakhir Diubah</p>
                        <p>{{ product.updated_at }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Status Produk</p>
                        <p>{{ product.is_empty == 0 ? "Produk aktif" : "Produk tidak aktif" }}</p>
                      </div>
                    </div>
                  </div>
                </div>
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
      product: [],
      isLoading: true,
      isLoading2: true,
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      });

      this.isLoading = false
    },
    getImage() {
      return this.product.image
    },
    submit(no_product) {
      if(confirm("Are you sure to delete this product?")) {
        axios.delete(`${process.env.VUE_APP_SERVICE_URL}/products/` + no_product, { headers: authHeader() })
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
}
</script>