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
                <h3 class="mb-0">Create Manual Transaction</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                 <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Product Name</label> <small class="text-danger">*use ctrl for selecting the product</small>
                    <select multiple class="form-control" v-model="transaction.product_name" required>
                        <option v-for="(item) in products" :key="item.id_product">{{ item.name }}</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Customer Name</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" v-model="transaction.customer_name" required>
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
import Sidebar from "@/components/admin/Sidebar.vue"
import Topbar from "@/components/admin/Topbar.vue"
import Header from "@/components/admin/Header.vue"
import Footer from "@/components/admin/Footer.vue"
import axios from "axios";
import $ from "jquery";
import authHeader from "@/service/auth-header";

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  },
  mounted() {
    axios.get(`${process.env.VUE_APP_SERVICE_URL}/products`, { headers: authHeader() }).then((response) => {
      this.products = response.data.data;
      setTimeout(function(){
        $('#datatable').DataTable();
      }, 0);
    });
  },
  data: function () {
    return {
      products: [],
      transaction: {
        product_name: [],
        customer_name: ""
      }
    };
  },
  methods: {
    submit() {
      if (this.transaction.product_name.length > 0) {
        let productName = this.transaction.product_name
        this.transaction.product_name = productName.join(", ")
      }

      axios.post(`${process.env.VUE_APP_SERVICE_URL}/transactions`, this.transaction, { headers: authHeader() })
            .then(response => {
              if(response.data.code === 200) {
                alert(response.data.status)
                this.$router.push({
                  name: 'transaction'
                })
              }
            }).catch(error => {
                console.log(error.response.data.status)
            })
    }
  }
}
</script>