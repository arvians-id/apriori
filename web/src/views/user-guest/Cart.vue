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
                        <a href="#" class="avatar">
                          <img alt="Image placeholder" :src="getImage(item.image)">
                        </a>
                      </div>
                      <div class="col ml--2">
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
                      <div class="col-auto">
                        <p>Rp {{ numberWithCommas(item.totalPricePerItem) }} - {{ item.quantity }} item</p>
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
                      <div class="col-auto text-center">
                        <p>Total harga : Rp {{ numberWithCommas(totalPrice) }}</p>
                        <a :href="send(carts, totalPrice)" class="btn btn-primary btn-sm">Pesan Sekarang</a>
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
</style>
<script>
import Sidebar from "@/components/guest/Sidebar.vue"
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"
import Footer from "@/components/guest/Footer.vue"

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
      totalPrice: 0,
    }
  },
  mounted() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      localStorage.getItem("my-carts")
        ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
        : (this.carts = []);

      this.carts.map(item => {
        this.totalPrice += item.price * item.quantity;
      })
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
      }

      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    min(item){
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem.quantity > 1) {
        productItem.quantity -= 1
        productItem.totalPricePerItem -= productItem.price
      } else {
        this.carts.splice(this.carts.indexOf(productItem), 1);
      }

      this.totalPrice -= productItem.price
      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    send(item, totalPrice) {
      let text = `Hallo saya ingin memesan produk-produk berikut : %0a${item.map(item => `${item.name} - ${item.quantity} x`).join('%0a')} %0aTotal harga : Rp ${this.numberWithCommas(totalPrice)}`;
      return `whatsapp://send/?phone=${process.env.VUE_APP_PHONE_NUMBER}&text=${text}`
    }
  }
}
</script>