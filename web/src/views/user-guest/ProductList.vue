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
              <div class="card-body">
                <div class="row">
                  <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in products" :key="item.id_product">
                    <div class="card">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" :src="getImage(item.image)" alt="Preview Image">
                      </div>
                      <ul class="list-group list-group-flush">
                        <li class="list-group-item text-danger font-weight-bold">Rp{{ numberWithCommas(item.price) }}</li>
                      </ul>
                      <div class="card-body">
                        <router-link :to="{ name: 'guest.product.detail', params: { code: item.code } }" class="card-title mb-3 text-dark">{{ item.name }}</router-link>
                        <p class="card-text mb-4">{{ item.description.length > 35 ? item.description.slice(0, 35) + "..." : item.description }}</p>
                        <div class="d-flex justify-content-between">
                          <button type="button" @click="add(item)" class="btn btn-primary btn-sm">Tambah</button>
                          <button type="button" @click="min(item)" class="btn btn-danger btn-sm">Kurangi</button>
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
  .card-img-top {
    width: 100%;
    object-fit: cover;
  }
</style>

<script>
import Sidebar from "@/components/guest/Sidebar.vue"
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"
import Footer from "@/components/guest/Footer.vue"
import axios from "axios";

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
      products: [],
      carts: [],
    };
  },
  methods: {
    fetchData() {
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/products`).then((response) => {
        this.products = response.data.data;
      });
    },
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    },
    getImage(image) {
      return image;
    },
    add(item) {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem) {
        productItem.quantity += 1
        productItem.totalPricePerItem = productItem.price * productItem.quantity
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
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem.quantity > 1) {
        productItem.quantity -= 1
        productItem.totalPricePerItem -= productItem.price
      } else {
        this.carts.splice(this.carts.indexOf(productItem), 1);
      }
      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    }
  }
}
</script>