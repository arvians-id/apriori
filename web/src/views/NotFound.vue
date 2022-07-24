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
        <div class="col-12">
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card header -->
              <div class="card-header">
                <h3 class="mb-0">404 Not Found</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <img src="../assets/notfound.png" class="img-fluid mx-auto d-block">
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"

export default {
  components: {
    Header,
    Topbar
  },
  mounted() {
    this.fetchData()
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      carts: [],
      totalCart: 0
    };
  },
  methods: {
    fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }
    }
  }
}
</script>