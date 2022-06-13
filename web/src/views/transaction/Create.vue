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
                <h3 class="mb-0">Buat Transaksi Manual</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                 <form>
                  <div class="form-group">
                    <label class="form-control-label">Nama Produk</label> <small class="text-danger">*use ctrl for selecting the product</small>
                    <select multiple class="form-control" v-model="product_name">
                        <option v-for="(item) in products" :key="item.id_product">{{ item.name }}</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Nama Pelanggan</label>
                    <input type="text" class="form-control" v-model="customer_name">
                  </div>
                  <button class="btn btn-primary" type="submit">Submit form</button>
                </form>
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
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import axios from "axios";
import $ from "jquery";

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  },
  mounted() {
    axios.get("http://localhost:3000/api/products").then((response) => {
      this.products = response.data.data;
      setTimeout(function(){
        $('#datatable').DataTable();
      }, 0);
    });
  },
  data: function () {
    return {
      products: []
    };
  },
}
</script>