<template>
  <!-- Sidenav -->
  <Sidebar :totalNotification="totalNotification" />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar :totalCart="totalCart" :carts="carts" :totalNotification="totalNotification" :notifications="notifications" />
    <!-- Header -->
    <Header />
    <!-- Page content -->
    <div class="container-fluid mt--6">
      <div class="row">
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile">
            <!-- Card header -->
            <div class="card-body mt-5">
              <div v-if="isLoading">
                <div class="loading-skeleton row d-flex justify-content-center">
                  <div class="col-12 col-lg-4 mb-2">
                    <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/no-image.png" class="img-fluid mb-2">
                    <div class="border p-3 rounded text-center" style="color: #525f7f">
                      <p class="mb-0">Atur jumlah yang pembelian</p>
                      <div class="mt-2">
                        <button type="button" class="btn btn-danger btn-sm">-</button>
                        <button class="btn disabled">quantity</button>
                        <button type="button" class="btn btn-primary btn-sm">+</button>
                      </div>
                    </div>
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
                              <hr class="m-0 mb-3">
                              <div class="media">
                                <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/ryzy.jpg" width="59" class="mr-3" alt="...">
                                <div class="media-body">
                                  <p class="mt-0 mb-0 w-50 mb-1">Title</p>
                                  <p>Description</p>
                                </div>
                              </div>
                              <hr class="m-0 mb-3">
                              <p class="font-weight-bold mb-0 mt-2">Pengiriman</p>
                              <p class="mb-1"><i class="ni ni-pin-3"></i> Dikirim dari Tanggerang, Banten</p>
                              <p><i class="ni ni-delivery-fast"></i> Tersedia pengiriman dengan TIKI, JNE dan POS Indonesia</p>
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
                  <img :src="apriori.apriori_image" class="img-fluid mb-2" width="500">
                  <div class="border p-3 rounded text-center" style="color: #525f7f">
                    <p class="mb-0">Atur jumlah yang pembelian</p>
                    <div>
                      <button type="button" @click="min(apriori)" class="btn btn-danger btn-sm">-</button>
                      <button class="btn disabled">{{ quantity }} item</button>
                      <button type="button" @click="add(apriori)" class="btn btn-primary btn-sm">+</button>
                    </div>
                  </div>
                </div>
                <div class="col-12 col-lg-6">
                  <div class="text-left">
                    <h5 class="h2 text-uppercase p-0 m-0">Paket rekomendasi {{ apriori.apriori_item }}</h5>
                    <p class="p-0 m-0">code : {{ apriori.apriori_code }}</p>
                    <div class="h1 font-weight-bold">Rp. {{ apriori.price_discount }}</div>
                    <small class="p-1 d-inline rounded font-weight-bold" style="background-color: #ffeaef;; color: #ff5c84">{{ apriori.apriori_discount }}%</small>
                    <p class="d-inline ml-2" style="text-decoration: line-through">Rp. {{ apriori.product_total_price }}</p>
                    <hr class="mt-0 mb-1">
                    <div class="nav-wrapper">
                      <ul class="nav nav-pills nav-fill flex-column flex-md-row" id="tabs-icons-text" role="tablist">
                        <li class="nav-item">
                          <a class="nav-link mb-sm-3 mb-md-0 active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true"><i class="ni ni-tag mr-2"></i>Detail</a>
                        </li>
                      </ul>
                    </div>
                    <div class="card shadow">
                      <div class="card-body">
                        <div class="tab-content" id="myTabContent">
                          <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                            <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                            <p>Original baru</p>
                            <p class="font-weight-bold mb-0 mt-2">Jumlah Barang</p>
                            <p>{{ apriori.apriori_item.split(", ").length }}</p>
                            <p class="font-weight-bold mb-0 mt-2">Nama Barang</p>
                            <p>{{ UpperWord(apriori.apriori_item) }}</p>
                            <p class="font-weight-bold mb-0 mt-2">Berat Barang</p>
                            <p>{{ apriori.mass == undefined ? 0 : apriori.mass }} gram</p>
                            <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                            <p>{{ apriori.apriori_description == undefined ? "Tidak ada deskripsi" : apriori.apriori_description }}</p>
                            <hr class="m-0 mb-3">
                            <div class="media">
                              <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/ryzy.jpg" width="53" class="mr-3" alt="...">
                              <div class="media-body">
                                <h3 class="mt-0 mb-0">Toko Ryzy Olshop</h3>
                                <p>Produk Original Berkualitas dan Terpercaya..</p>
                              </div>
                            </div>
                            <hr class="m-0 mb-3">
                            <p class="font-weight-bold mb-0 mt-2">Pengiriman</p>
                            <p class="mb-1"><i class="ni ni-pin-3"></i> Dikirim dari Tanggerang, Banten</p>
                            <p><i class="ni ni-delivery-fast"></i> Tersedia pengiriman dengan TIKI, JNE dan POS Indonesia</p>
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
import Sidebar from "@/components/guest/Sidebar.vue"
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"
import Footer from "@/components/guest/Footer.vue"
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
    if(authHeader()["Authorization"] !== undefined) {
      this.fetchNotification()
    }
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      apriori: [],
      carts: [],
      totalCart: 0,
      quantity: 0,
      isLoading: true,
      totalNotification: 0,
      notifications: []
    };
  },
  methods: {
    async fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.apriori = response.data.data;
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      let productItem = this.carts.find(product => product.code == this.$route.params.id);
      this.quantity = productItem ? productItem.quantity : 0;

      this.isLoading = false;
    },
    async fetchNotification() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/notifications/user`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalNotification = response.data.data.filter(e => e.is_read === false).length
          this.notifications = response.data.data
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    },
    send(product) {
      return `whatsapp://send?phone=${process.env.VUE_APP_PHONE_NUMBER}&text=Hallo, saya ingin membeli produk code ${product.apriori_code} dengan paket rekomendasi ${product.apriori_item} seharga Rp${product.price_discount}. Apakah produk masih tersedia?`
    },
    add(item) {
      let productItem = this.carts.find(product => product.code == item.apriori_id);
      if (productItem) {
        productItem.quantity += 1
        productItem.totalPricePerItem = productItem.price * productItem.quantity
        this.quantity += 1
        this.totalCart += 1
      } else {
        this.carts.push({
          id_product: item.apriori_code,
          code: item.apriori_id,
          name: "Paket Rekomendasi " + this.UpperWord(item.apriori_item),
          price: item.price_discount,
          mass: item.mass == undefined ? 0 : item.mass,
          image: item.apriori_image,
          quantity: 1,
          totalPricePerItem: item.price_discount
        });
        this.quantity = 1
        this.totalCart += 1
      }

      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    min(item){
      let productItem = this.carts.find(product => product.code === item.apriori_id);
      if (productItem !== undefined) {
        if (productItem.quantity > 1) {
          productItem.quantity -= 1
          productItem.totalPricePerItem -= productItem.price
          this.quantity -= 1
          this.totalCart -= 1
        } else {
          this.quantity = 0
          this.totalCart -= 1
          this.carts.splice(this.carts.indexOf(productItem), 1);
        }
        localStorage.setItem('my-carts', JSON.stringify(this.carts));
      }
    }
  }
}
</script>