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
              <h3 class="mb-0">Detail Recommendation</h3>
            </div>
            <div class="card-body">
              <div class="row align-items-center mx-auto" v-if="isLoading">
                <p class="p-3 mt-2 text-center">Loading...</p>
              </div>
              <div class="row d-flex justify-content-center" v-else>
                <div class="col-12 col-lg-4 text-center mb-2">
                  <img :src="getImage()" class="img-fluid mb-2" width="500">
                  <router-link :to="{ name: 'apriori.edit', params: { code: apriori.apriori_code, id: this.$route.params.id } }"  class="btn btn-primary btn-block mb-2">Edit</router-link>
                </div>
                <div class="col-12 col-lg-4">
                  <div class="text-left">
                    <h5 class="h2 text-uppercase p-0 m-0">Paket rekomendasi {{ apriori.apriori_item }}</h5>
                    <p class="p-0 m-0">code : {{ apriori.apriori_code }}</p>
                    <div class="h1 font-weight-bold">Rp. {{ apriori.price_discount }}</div>
                    <small class="p-1 d-inline rounded font-weight-bold" style="background-color: #ffeaef;; color: #ff5c84">{{ apriori.apriori_discount }}%</small>
                    <p class="d-inline ml-2" style="text-decoration: line-through">Rp. {{ apriori.product_total_price }}</p>
                    <hr class="m-2">
                    <ul class="nav nav-tabs" id="myTab" role="tablist">
                      <li class="nav-item">
                        <a class="nav-link active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true">Detail</a>
                      </li>
                    </ul>
                    <div class="tab-content" id="myTabContent">
                      <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                        <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                        <p>Original baru</p>
                        <p class="font-weight-bold mb-0 mt-2">Jumlah Barang</p>
                        <p>{{ apriori.apriori_item.split(", ").length }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Nama Barang</p>
                        <p>{{ UpperWord(apriori.apriori_item) }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                        <p>{{ apriori.apriori_description == undefined ? "Tidak ada deskripsi" : apriori.apriori_description }}</p>
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
      apriori: [],
      isLoading: true
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then((response) => {
        this.apriori = response.data.data;
      });

      this.isLoading = false;
    },
    getImage() {
      return this.apriori.apriori_image
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>