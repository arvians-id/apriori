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
              <div v-if="isLoading">
                <div class="loading-skeleton row d-flex justify-content-center">
                  <div class="col-12 col-lg-4 mb-2">
                    <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/no-image.png" class="img-fluid mb-2">
                    <p class="pt-4 mb-0">this is button</p>
                    <p class="pt-4 mb-0 mt-2">this is button</p>
                  </div>
                  <div class="col-12 col-lg-6">
                    <div class="text-left">
                      <h5 class="h2 p-0 m-0 mb-1">This is title of product</h5>
                      <p class="p-0 m-0 mb-1 w-50">This is title of product</p>
                      <p class="w-50">Pricing</p>
                      <hr class="mt-0 mb-1">
                      <div class="nav-wrapper">
                        <p class="pt-4 mb-0">this is button</p>
                      </div>
                      <div class="card shadow">
                        <div class="card-body">
                          <div class="tab-content" id="skeleton-myTabContent">
                            <div class="tab-pane fade show active" id="skeleton-detail" role="tabpanel" aria-labelledby="detail-tab">
                              <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                              <p>Original baru</p>
                              <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                              <p class="font-weight-bold">Text</p>
                              <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                              <p>gram</p>
                              <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                              <p>Description</p>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="row d-flex justify-content-center" v-else>
                <div class="col-12 col-lg-4 text-center mb-2">
                  <img :src="product.image" class="img-fluid mb-2" width="500">
                  <router-link :to="{ name: 'product.edit', params: { code: this.$route.params.code } }"  class="btn btn-primary btn-block mb-2">Edit</router-link>
                  <form @submit.prevent="submit(this.$route.params.code)" method="POST" class="d-inline">
                    <button class="btn btn-danger btn-block">Delete</button>
                  </form>
                </div>
                <div class="col-12 col-lg-6">
                  <div class="text-left">
                    <h5 class="h2 text-uppercase p-0 m-0">{{ product.name }}</h5>
                    <p class="p-0 m-0">code : {{ product.code }}</p>
                    <div class="h1 font-weight-bold">Rp. {{ product.price }}</div>
                    <hr class="mt-0 mb-1">
                    <div class="nav-wrapper">
                      <ul class="nav nav-pills nav-fill flex-column flex-md-row" id="tabs-icons-text" role="tablist">
                        <li class="nav-item">
                          <a class="nav-link mb-sm-3 mb-md-0 active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true"><i class="ni ni-tag mr-2"></i>Detail</a>
                        </li>
                        <li class="nav-item">
                          <a class="nav-link mb-sm-3 mb-md-0" id="lainnya-tab" data-toggle="tab" href="#lainnya" role="tab" aria-controls="lainnya" aria-selected="false"><i class="ni ni-ungroup mr-2"></i>Lainnya</a>
                        </li>
                      </ul>
                    </div>
                    <div class="card shadow">
                      <div class="card-body">
                        <div class="tab-content" id="myTabContent">
                          <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                            <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                            <p>Original baru</p>
                            <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                            <p>
                              <router-link :to="{ name: 'guest.product', query: { category: category } }" v-for="(category, i) in product.category.split(', ')" :key="i" class="text-primary font-weight-bold">
                                {{ category }}{{ product.category.split(', ').length - 1 != i ? ', ' : '' }}
                              </router-link>
                            </p>
                            <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                            <p>{{ product.mass }} gram</p>
                            <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                            <p>{{ product.description == "" ? "Tidak ada deskripsi" : product.description }}</p>
                          </div>
                          <div class="tab-pane fade" id="lainnya" role="tabpanel" aria-labelledby="lainnya-tab">
                            <p class="font-weight-bold mb-0 mt-2">Tanggal Dibuat</p>
                            <p>{{ product.created_at }}</p>
                            <p class="font-weight-bold mb-0 mt-2">Terakhir Diubah</p>
                            <p>{{ product.updated_at }}</p>
                          </div>
                        </div>
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

<style scoped>
@import '../../assets/skeleton.css';
</style>

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
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading = false
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
              if (error.response.status === 400 || error.response.status === 404) {
                alert(error.response.data.status)
              }
            })
      }
    }
  }
}
</script>