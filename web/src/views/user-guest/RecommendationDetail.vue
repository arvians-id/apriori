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
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile">
            <!-- Card header -->
            <div class="card-header">
              <h3 class="mb-0">Detail Paket Rekomendasi</h3>
            </div>
            <div class="row align-items-center">
              <div class="col-12 col-lg-6 text-center">
                <img :src="getImage()" class="img-fluid my-5" width="500">
              </div>
              <div class="col-12 col-lg-6">
                <div class="card-pricing border-0 text-center mb-4">
                  <div class="card-body px-lg-7">
                    <div class="display-2">{{ apriori.apriori_discount }}%</div>
                    <span class="text-muted h2" style="text-decoration: line-through">Rp. {{ apriori.product_total_price }}</span>
                    /
                    <span class="text-muted">Rp. {{ apriori.price_discount }}</span>
                    <div class="text-center mt-3">
                      <h5 class="h3 text-uppercase">
                       Paket rekomendasi {{ apriori.apriori_item }} - {{ apriori.apriori_code }}
                      </h5>
                    </div>
                    <ul class="list-unstyled" v-if="apriori.apriori_item">
                      <li v-for="(value,i) in apriori.apriori_item.split(', ')" :key="i">
                        <div class="d-flex align-items-center justify-content-center">
                          <div>
                            <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                              <i class="ni ni-basket"></i>
                            </div>
                          </div>
                          <div>
                            <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                          </div>
                        </div>
                      </li>
                    </ul>
                    <p>{{ apriori.apriori_description }}</p>
                    <button type="button" @click="min(apriori)" class="btn btn-danger btn-sm mt-3">-</button>
                    <button class="btn disabled mt-3">{{ quantity }} item</button>
                    <button type="button" @click="add(apriori)" class="btn btn-primary btn-sm mt-3">+</button>
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
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      apriori: [],
      carts: [],
      totalCart: 0,
      quantity: 0,
    };
  },
  methods: {
    fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then((response) => {
        this.apriori = response.data.data;
      });

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      let productItem = this.carts.find(product => product.code == this.$route.params.id);
      this.quantity = productItem ? productItem.quantity : 0;
    },
    getImage() {
      return this.apriori.apriori_image
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