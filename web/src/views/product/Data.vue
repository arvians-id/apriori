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
              <h3 class="mb-0">Data Produk</h3>
            </div>
            <div class="table-responsive py-4">
              <table class="table table-flush" id="datatable">
                <thead class="thead-light">
                <tr>
                  <th>No</th>
                  <th>Code Produk</th>
                  <th>Nama</th>
                  <th>Deskripsi</th>
                  <th>Tanggal Dibuat</th>
                  <th>Terakhir Diubah</th>
                  <th class="text-center">Action</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item,i) in products" :key="item.id_product">
                  <td>{{ (i++) + 1 }}</td>
                  <td>{{ item.code }}</td>
                  <td>{{ item.name }}</td>
                  <td>{{ item.description }}</td>
                  <td>{{ item.created_at }}</td>
                  <td>{{ item.updated_at }}</td>
                  <td class="text-center">
                    <router-link to="" class="btn btn-primary btn-sm">Ubah</router-link>
                    <form @submit.prevent="submit(item.code)" method="POST" class="d-inline">
                      <button class="btn btn-danger btn-sm">Hapus</button>
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
import "jquery/dist/jquery.min.js";
import "datatables.net-dt/js/dataTables.dataTables";
import "datatables.net-dt/css/jquery.dataTables.min.css";
import axios from "axios";
import $ from "jquery";
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"

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
      products: []
    };
  },
  methods: {
    fetchData() {
      axios.get("http://localhost:3000/api/products").then((response) => {
        this.products = response.data.data;
        setTimeout(function(){
          $('#datatable').DataTable();
        }, 0);
      });
    },
    submit(no_product) {
      axios.delete("http://localhost:3000/api/products/" + no_product)
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.fetchData()
            }
          }).catch(error => {
        alert(error.response.data.status)
      })
    }
  }
}
</script>