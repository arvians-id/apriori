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
              <h3 class="mb-0">Datatable</h3>
              <p class="text-sm mb-0">
                This is an exmaple of datatable using the well known datatables.net plugin. This is a minimal setup in order to get started fast.
              </p>
            </div>
            <div class="table-responsive py-4">
              <table class="table table-flush" id="datatable">
                <thead class="thead-light">
                  <tr>
                    <th>No Transaksi</th>
                    <th>Nama Pelanggan</th>
                    <th>Nama Produk</th>
                    <th>Tanggal Dibuat</th>
                    <th>Terakhir Diubah</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in products" :key="item.id_transaction">
                    <td>{{ item.no_transaction }}</td>
                    <td>{{ item.customer_name }}</td>
                    <td>{{ item.product_name }}</td>
                    <td>{{ item.created_at }}</td>
                    <td>{{ item.updated_at }}</td>
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
    axios.get("http://localhost:3000/api/transactions").then((response) => {
      this.products = response.data.data;
      setTimeout(function(){
        $('#datatable').DataTable();
      }, 0);
    });
  },
  data: function () {
    return {
      products: [],
    };
  },
}
</script>