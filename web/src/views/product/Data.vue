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
      <!-- Table -->
      <div class="row">
        <div class="col">
          <div class="card">
            <!-- Card header -->
            <div class="card-header">
              <h3 class="mb-0">Data Product</h3>
            </div>
            <div class="table-responsive py-3" v-if="isLoading">
              <p class="mt-2 text-center">Loading...</p>
            </div>
            <div class="table-responsive py-4" v-else>
              <table class="table table-flush" id="datatable">
                <thead class="thead-light">
                <tr>
                  <th>No</th>
                  <th>Product Code</th>
                  <th>Name</th>
                  <th>Price</th>
                  <th>Description</th>
                  <th class="text-center">Action</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item,i) in products" :key="item.id_product">
                  <td>{{ (i++) + 1 }}</td>
                  <td>{{ item.code }}</td>
                  <td>{{ item.name }}</td>
                  <td>Rp. {{ item.price }}</td>
                  <td>{{ item.description.length > 50 ? item.description.slice(0, 50) + "..." : item.description }}</td>
                  <td class="text-center">
                    <router-link :to="{ name: 'product.detail', params: { code: item.code } }" class="btn btn-secondary btn-sm">Detail</router-link>
                    <router-link :to="{ name: 'product.edit', params: { code: item.code } }" class="btn btn-primary btn-sm">Edit</router-link>
                    <form @submit.prevent="submit(item.code)" method="POST" class="d-inline">
                      <button class="btn btn-danger btn-sm">Delete</button>
                    </form>
                  </td>
                </tr>
                </tbody>
              </table>
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
import axios from "axios";
import $ from "jquery";
import Sidebar from "@/components/admin/Sidebar.vue"
import Topbar from "@/components/admin/Topbar.vue"
import Header from "@/components/admin/Header.vue"
import Footer from "@/components/admin/Footer.vue"
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
      products: [],
      isLoading: true
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products`, { headers: authHeader() }).then((response) => {
        this.products = response.data.data;
        setTimeout(function(){
          $('#datatable').DataTable();
        }, 0);
      });

      this.isLoading = false;
    },
    submit(no_product) {
      if(confirm("Are you sure to delete this product?")) {
        axios.delete(`${process.env.VUE_APP_SERVICE_URL}/products/` + no_product, { headers: authHeader() })
            .then(response => {
              if(response.data.code === 200) {
                alert(response.data.status)
                this.fetchData()
              }
            }).catch(error => {
          console.log(error.response.data.status)
        })
      }
    }
  }
}
</script>