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
              <h3 class="mb-0">Data Apriori</h3>
            </div>
            <div class="table-responsive py-4">
              <table class="table table-flush" id="datatable">
                <thead class="thead-light">
                <tr>
                  <th>No</th>
                  <th>Code</th>
                  <th>Rentang Tanggal</th>
                  <th>Status</th>
                  <th>Terakhir Dibuat</th>
                  <th class="text-center">Action</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item,i) in apriories" :key="item.id_apriori">
                  <td>{{ (i++) + 1 }}</td>
                  <td>{{ item.code }}</td>
                  <td>{{ item.range_date }}</td>
                  <td>{{ item.is_active }}</td>
                  <td>{{ item.created_at }}</td>
                  <td class="text-center">
                    <router-link to="" class="btn btn-secondary btn-sm">Detail</router-link>
                    <router-link to="" class="btn btn-primary btn-sm">Ubah</router-link>
                    <form class="d-inline" onsubmit="return confirm(`Apakah anda yakin ingin menghapus data ini?`)">
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
    axios.get("http://localhost:3000/api/apriori").then((response) => {
      this.apriories = response.data.data;
      setTimeout(function(){
        $('#datatable').DataTable();
      }, 0);
    });
  },
  data: function () {
    return {
      apriories: [],
    };
  },
}
</script>