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
        <div class="col-12 col-md-6">
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card header -->
              <div class="card-header">
                <!-- Title -->
                <h5 class="h3 mb-0">Konfirmasi Pemesanan</h5>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <!-- List group -->
                <form @submit.prevent="submit" method="POST">
                  <div class="row">
                    <div class="col-6">
                      <div class="form-group">
                        <label class="form-control-label">Nama Depan</label> <small class="text-danger">*</small>
                        <input type="text" class="form-control" name="first_name" v-model="checkout.first_name" required readonly>
                      </div>
                    </div>
                    <div class="col-6">
                      <div class="form-group">
                        <label class="form-control-label">Nama Belakang</label>
                        <input type="text" class="form-control" name="last_name" v-model="checkout.last_name" readonly>
                      </div>
                    </div>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Provinsi</label> <small class="text-danger">*</small>
                    <select class="form-control" v-model="province_id" @change="getCity" required>
                      <option value="" disabled selected>Select</option>
                      <option v-for="province in provinces" :value="province.province_id" :key="province.province_id">{{ province.province }}</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Kota/Kabupaten</label> <small class="text-danger">*</small>
                    <select class="form-control" v-model="city_id" required>
                      <option value="" disabled selected>Select</option>
                      <option v-for="city in cities" :value="city.city_id" :key="city.city_id">{{ city.type + " " + city.city_name + ", " + city.postal_code }}</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Alamat Lengkap</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" name="full_address" v-model="checkout.address" required>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Jasa Ekspedisi</label> <small class="text-danger">*</small>
                    <select class="form-control" v-model="checkout.courier" @change="getCost" required>
                      <option value="" disabled selected>Select</option>
                      <option value="jne">JNE</option>
                      <option value="tiki">TIKI</option>
                      <option value="pos">POS Indonesia</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Layanan Pengiriman</label> <small class="text-danger">*</small>
                    <select class="form-control" v-model="totalCost" @change="getTotalPrice" required>
                      <option value="" disabled selected>Select</option>
                      <template v-if="costs.results !== undefined">
                        <option v-for="(cost,i) in costs.results[0].costs" :value="getCourierService(cost.service,cost.description,cost.cost[0].etd,cost.cost[0].value)" :key="i">{{ cost.service + " ("+ cost.description +") | " + cost.cost[0].etd + " (day) | Rp." + cost.cost[0].value }}</option>
                      </template>
                    </select>
                  </div>
                  <button class="btn btn-primary" :disabled="carts.length < 1">Pesan</button>
                </form>
              </div>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6">
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card body -->
              <div class="card-body">
                <!-- List group -->
                <ul class="list-group list-group-flush list my--3 p-3" v-if="carts.length > 0">
                  <li class="list-group-item px-0" v-for="(item, i) in carts" :key="i">
                    <div class="row align-items-center">
                      <div class="col-auto">
                        <!-- Avatar -->
                        <router-link
                            :to="{ name: 'guest.product.recommendation', params: { code: item.id_product, id: item.code } }"
                            v-if="item.name.includes(`Paket Rekomendasi`)"
                            class="avatar">
                          <img alt="Image placeholder" :src="item.image">
                        </router-link>
                        <router-link
                            :to="{ name: 'guest.product.detail', params: { code: item.code } }"
                            class="avatar"
                            v-else>
                          <img alt="Image placeholder" :src="item.image">
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
                      </div>
                      <div class="col-auto d-none d-lg-block">
                        <p class="text-center">Rp {{ numberWithCommas(item.totalPricePerItem) }} - {{ item.quantity }} item</p>
                      </div>
                    </div>
                  </li>
                  <li class="list-group-item px-0">
                    <div class="d-flex justify-content-between">
                      <div>
                        <div class="p-0 my-2">Subtotal</div>
                        <div class="p-0 my-2">Pajak</div>
                        <div class="p-0 my-2">Ongkos Kirim</div>
                      </div>
                      <div class="text-right">
                        <div class="p-0 my-2">Rp {{ numberWithCommas(totalPrice - 5000 - totalCost) }}</div>
                        <div class="p-0 my-2">Rp {{ numberWithCommas(5000) }}</div>
                        <div class="p-0 my-2">{{ numberWithCommas(totalCost) == 0 ? "Kurir belum ditentukan." : "Rp " + numberWithCommas(totalCost) }}</div>
                      </div>
                    </div>
                  </li>
                  <li class="list-group-item px-0">
                    <div class="d-flex justify-content-between">
                      <div>Total Pembayaran</div>
                      <div class="font-weight-bold">Rp {{ numberWithCommas(totalPrice) }}</div>
                    </div>
                    <a href="javascript:void(0);" @click="clearCart" class="btn btn-danger btn-sm mt-2">Bersihkan pemesanan</a>
                  </li>
                </ul>
                <ul class="list-group list-group-flush list my--3 mx--3" v-else>
                  <li class="list-group-item">
                    <div class="alert alert-secondary">
                      <h5 class="alert-heading">Oops!</h5>
                      <p>Pemesanan belanjaan masih kosong nih. <router-link :to="{ name: 'guest.product' }" >Beli produk disini!</router-link></p>
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
import getRoles from "@/service/get-roles";

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
      totalCart: 0,
      checkout: {
        first_name: "",
        last_name: "",
        courier: "",
        address: "",
        courier_service: ""
      },
      province_id: "",
      city_id: "",
      provinces: [],
      cities: [],
      costs: [],
      totalCost: 0,
      totalNotification: 0,
      notifications: [],
    }
  },
  async mounted() {
    this.fetchData()
    if(authHeader()["Authorization"] !== undefined) {
      this.loadScript()
      await this.fetchNotification()
      let getRole = await getRoles();
      let names = getRole.name.split(" ")
      if (names.length < 2) {
        this.checkout.first_name = names[0]
      } else {
        this.checkout.first_name = names[0]
        this.checkout.last_name = names[1]
      }
    }
    this.getProvinces()
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  methods: {
    getProvinces(){
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/raja-ongkir/province`, { headers: authHeader() })
          .then(response => {
            this.provinces = response.data.data.rajaongkir.results
          }).catch(error => {
            console.log(error)
          });
    },
    getCity(){
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/raja-ongkir/city?province=${this.province_id}`, { headers: authHeader() })
          .then(response => {
            this.cities = response.data.data.rajaongkir.results
          }).catch(error => {
        console.log(error)
          });
    },
    getCost(){
      let headers = {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
          ...authHeader()
        }
      }

      let getWeight = this.carts.map(item => {
        return item.mass * item.quantity
      }).reduce((a, b) => a + b, 0)

      let formData = new FormData()
      formData.append("origin", 457)
      formData.append("destination", this.city_id)
      formData.append("weight", getWeight)
      formData.append("courier", this.checkout.courier)

      axios.post(`${process.env.VUE_APP_SERVICE_URL}/raja-ongkir/cost`, formData, headers)
          .then(response => {
            this.costs = response.data.data.rajaongkir
          }).catch(error => {
            console.log(error)
          });
    },
    getTotalPrice(){
      this.totalPrice = 5000
      this.carts.map(item => {
        this.totalPrice += item.price * item.quantity;
      })
      this.totalPrice += this.totalCost
    },
    getCourierService(service, note, etd, cost){
      this.checkout.courier_service = service + " ("+ note +") | " + etd + " (day) | Rp." + cost
      return cost
    },
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
      })
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
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    },
    async submit() {
      if(authHeader()["Authorization"] === undefined) {
        this.$router.push({ name: 'auth.login' })
      } else {
        getRoles().then(response => {
          let formData = new FormData()

          let names = this.checkout.last_name === "" ? this.checkout.first_name : this.checkout.first_name + " " + this.checkout.last_name
          formData.append("gross_amount", this.totalPrice)
          formData.append("user_id", response.id_user)
          formData.append("customer_name", names)
          formData.append("items", JSON.stringify(this.carts))
          formData.append("address", this.checkout.address)
          formData.append("courier", this.checkout.courier)
          formData.append("courier_service", this.checkout.courier_service)
          formData.append("shipping_cost", this.totalCost)

          axios.post(`${process.env.VUE_APP_SERVICE_URL}/payments/pay`, formData, { headers: authHeader() }).then(response => {
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
            if (error.response.status === 400 || error.response.status === 404) {
              console.log(error.response.data.status)
            }
          })
        })
      }
    },
    clearCart() {
      if(confirm("Apakah anda yakin ingin menghapus semua pesanan anda?")){
        this.carts = []
        localStorage.setItem('my-carts', JSON.stringify([]));
        this.totalCart = 0
        this.totalPrice = 0
      }
    }
  }
}
</script>