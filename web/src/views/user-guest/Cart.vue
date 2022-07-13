<template>
  <!-- Sidenav -->
  <Sidebar />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar :totalCart="totalCart" :carts="carts" />
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
                <!-- Title -->
                <h5 class="h3 mb-0">Keranjang Belanja</h5>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <!-- List group -->
                <ul class="list-group list-group-flush list my--3" v-if="carts.length > 0">
                  <li class="list-group-item px-0" v-for="(item, i) in carts" :key="i">
                    <div class="row align-items-center">
                      <div class="col-auto">
                        <!-- Avatar -->
                        <router-link
                            :to="{ name: 'guest.product.recommendation', params: { code: item.id_product, id: item.code } }"
                            v-if="item.name.includes(`Paket Rekomendasi`)"
                            class="avatar">
                          <img alt="Image placeholder" :src="getImage(item.image)">
                        </router-link>
                        <router-link
                            :to="{ name: 'guest.product.detail', params: { code: item.code } }"
                            class="avatar"
                            v-else>
                          <img alt="Image placeholder" :src="getImage(item.image)">
                        </router-link>
                      </div>
                      <div class="col ml--4">
                        <h4 class="mb-0">
                          <router-link
                              :to="{ name: 'guest.product.recommendation', params: { code: item.id_product, id: item.code } }"
                              v-if="item.name.includes(`Paket Rekomendasi`)">{{ item.name }}
                          </router-link>
                          <router-link
                              :to="{ name: 'guest.product.detail', params: { code: item.code } }"
                              v-else>{{ item.name }}
                          </router-link>
                        </h4>
                        <small>Rp {{ numberWithCommas(item.price) }}</small>
                      </div>
                      <div class="col-12 d-block d-lg-none">
                        <p>Rp {{ numberWithCommas(item.totalPricePerItem) }} - {{ item.quantity }} item</p>
                        <button type="button" @click="min(item)" class="btn btn-danger btn-sm">Kurangi</button>
                        <button type="button" @click="add(item)" class="btn btn-primary btn-sm">Tambah</button>
                      </div>
                      <div class="col-auto d-none d-lg-block">
                        <p class="text-center">Rp {{ numberWithCommas(item.totalPricePerItem) }} - {{ item.quantity }} item</p>
                        <button type="button" @click="min(item)" class="btn btn-danger btn-sm">Kurangi</button>
                        <button type="button" @click="add(item)" class="btn btn-primary btn-sm">Tambah</button>
                      </div>
                    </div>
                  </li>
                  <li class="list-group-item px-0">
                    <div class="row align-items-center">
                      <div class="col-auto">
                      </div>
                      <div class="col">
                      </div>
                      <div class="col-auto">
                        <table class="table table-responsive table-borderless">
                          <tr>
                            <td>Sub Total</td>
                            <td>:</td>
                            <td class="text-right">Rp {{ numberWithCommas(totalPrice - 5000) }}</td>
                          </tr>
                          <tr>
                            <td>Pajak Aplikasi</td>
                            <td>:</td>
                            <td class="text-right">Rp {{ numberWithCommas(5000) }}</td>
                          </tr>
                          <tr>
                            <td>Total harga</td>
                            <td>:</td>
                            <td class="text-right font-weight-bold"><h3>Rp {{ numberWithCommas(totalPrice) }}</h3></td>
                          </tr>
                          <tr>
                            <td></td>
                            <td></td>
                            <td>
                              <a href="javascript:void(0);" @click="clearCart" class="btn btn-danger btn-sm">Bersihkan keranjang</a>
                              <button @click="send" class="btn btn-primary btn-sm">Pesan Sekarang</button>
                            </td>
                          </tr>
                        </table>
                      </div>
                    </div>
                  </li>
                </ul>
                <ul class="list-group list-group-flush list" v-else>
                  <li class="list-group-item">
                    <div class="alert alert-secondary">
                      <h5 class="alert-heading">Oops!</h5>
                      <p>Keranjang belanjaan masih kosong nih. <router-link :to="{ name: 'guest.product' }" >Beli produk disini!</router-link></p>
                    </div>
                  </li>
                </ul>
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
.avatar{
  background-color: transparent;
}
.table-borderless > tbody > tr > td,
.table-borderless > tbody > tr > th,
.table-borderless > tfoot > tr > td,
.table-borderless > tfoot > tr > th,
.table-borderless > thead > tr > td,
.table-borderless > thead > tr > th {
  border: none;
}
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
  data(){
    return {
      carts: [],
      totalPrice: 5000,
      totalCart: 0
    }
  },
  mounted() {
    this.fetchData()
    if(authHeader()["Authorization"] !== undefined) {
      this.loadScript()
    }
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  methods: {
    fetchData() {
      localStorage.getItem("my-carts")
        ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
        : (this.carts = []);

      this.carts.map(item => {
        this.totalPrice += item.price * item.quantity;
      })

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }
    },
    loadScript(){
      let midtransJs = "https://api.sandbox.midtrans.com/v2/assets/js/midtrans.min.js"
      let tagMidtransJs = document.createElement("script");
      tagMidtransJs.setAttribute("src", midtransJs);
      document.head.appendChild(tagMidtransJs);

      let formData = new FormData()
      formData.append("gross_amount", this.totalPrice)
      formData.append("user_id", "7")
      formData.append("items", JSON.stringify(this.carts))

      let snapJs = "https://app.sandbox.midtrans.com/snap/snap.js"
      let tagSnapJs = document.createElement("script");
      tagSnapJs.setAttribute("src", snapJs);
      tagSnapJs.setAttribute("data-client-key", "SB-Mid-client-1WI-DDXBFya0sHp_");
      document.head.appendChild(tagSnapJs);
    },
    getImage(image) {
      return image;
    },
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    },
    add(item) {
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem) {
        productItem.quantity += 1
        productItem.totalPricePerItem = productItem.price * productItem.quantity
        this.totalPrice += productItem.price
        this.totalCart += 1
      } else {
        this.carts.push({
          id_product: item.id_product,
          code: item.code,
          name: item.name,
          price: item.price,
          image: item.image,
          quantity: 1,
          totalPricePerItem: item.price
        });
        this.totalCart += 1
      }

      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    min(item){
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem.quantity > 1) {
        productItem.quantity -= 1
        productItem.totalPricePerItem -= productItem.price
        this.totalCart -= 1
      } else {
        this.carts.splice(this.carts.indexOf(productItem), 1);
        this.totalCart -= 1
      }

      this.totalPrice -= productItem.price
      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    send() {
      if(authHeader()["Authorization"] === undefined) {
        this.$router.push({ name: 'auth.login' })
      } else {
        let formData = new FormData()
        formData.append("gross_amount", this.totalPrice)
        formData.append("user_id", localStorage.getItem("user"))
        formData.append("customer_name", localStorage.getItem("name"))
        formData.append("items", JSON.stringify(this.carts))

        axios.post(`${process.env.VUE_APP_SERVICE_URL}/payments/pay`, formData).then(response => {
          window.snap.pay(response.data.data.token, {
            onSuccess: function(result) {
              console.log(result)
            },
            onPending: function(result) {
              console.log(result)
            },
            onError: function(result) {
              console.log(result)
            }
          })
        }).catch(error => {
          console.log(error)
        })
      }

    },
    clearCart() {
      if(confirm("Apakah anda yakin ingin menghapus semua keranjang belanjaan?")){
        this.carts = []
        localStorage.setItem('my-carts', JSON.stringify([]));
        this.totalCart = 0
        this.totalPrice = 0
      }
    }
  }
}
</script>