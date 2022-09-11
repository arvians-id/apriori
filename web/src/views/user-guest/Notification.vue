<template>
  <!-- Sidenav -->
  <Sidebar :totalNotification="totalNotification" />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar :totalCart="totalCart" :carts="carts" :totalNotification="totalNotification" :notifications="notifications" />
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
              <div class="card-header d-flex justify-content-between">
                <!-- Title -->
                <h5 class="h3 mb-0">Pemberitahuan Saya</h5>
                <form @submit.prevent="markAll" method="POST" v-if="totalNotification > 0">
                  <button class="btn btn-danger btn-sm">Tandai Semua Telah Dibaca</button>
                </form>
              </div>
              <!-- Card body -->
              <div class="card-body p-0" v-if="isLoading">
                <div class="loading-skeleton m-3">
                  <div>
                    <p class="w-25 mb-0 mb-2">Title</p>
                    <p class="m-0 p-0">Description</p>
                    <hr class="mb-3 mt-3">
                  </div>
                  <div>
                    <p class="w-25 mb-0 mb-2">Title</p>
                    <p class="m-0 p-0">Description</p>
                    <hr class="mb-3 mt-3">
                  </div>
                  <div>
                    <p class="w-25 mb-0 mb-2">Title</p>
                    <p class="m-0 p-0">Description</p>
                    <hr class="mb-3 mt-3">
                  </div>
                  <div>
                    <p class="w-25 mb-0 mb-2">Title</p>
                    <p class="m-0 p-0">Description</p>
                    <hr class="mb-3 mt-3">
                  </div>
                  <div>
                    <p class="w-25 mb-0 mb-2">Title</p>
                    <p class="m-0 p-0">Description</p>
                  </div>
                </div>
              </div>
              <div class="card-body p-0" v-else>
                <!-- List group -->
                <ul class="list-group list-group-flush" data-toggle="checklist" v-if="notifications.length > 0">
                  <template v-for="(item,i) in notifications" :key="i">
                    <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4" v-if="!item.is_read" style="background-color: #f6f9fc">
                      <div class="checklist-item checklist-item-secondary">
                        <div class="checklist-info">
                          <h5 class="checklist-title mb-0">{{ item.title }}</h5>
                          <small class="d-block">{{ item.description }}</small>
                          <small>{{ item.created_at }}</small>
                        </div>
                        <form @submit.prevent="mark(item.id_notification)" method="POST" v-if="!item.is_read">
                          <button class="btn btn-sm btn-success">Tandai dibaca</button>
                        </form>
                      </div>
                    </li>
                    <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4" v-else>
                      <div class="checklist-item checklist-item-secondary">
                        <div class="checklist-info">
                          <h5 class="checklist-title mb-0">{{ item.title }}</h5>
                          <small class="d-block">{{ item.description }}</small>
                          <small>{{ item.created_at }}</small>
                        </div>
                        <form @submit.prevent="mark(item.id_notification)" method="POST" v-if="!item.is_read">
                          <button class="btn btn-sm btn-success">Tandai dibaca</button>
                        </form>
                      </div>
                    </li>
                  </template>
                </ul>
                <ul class="list-group list-group-flush list" v-else>
                  <li class="list-group-item">
                    <div class="alert alert-secondary">
                      <h5 class="alert-heading">Oops!</h5>
                      <p>Tidak ada pemeberitahuan baru.</p>
                    </div>
                  </li>
                </ul>
                <button @click="loadMore" v-if="notifications.length !== this.totalData" class="my-3 btn btn-secondary d-block mx-auto px-5">
                  Lihat lainnya <i class="ni ni-bold-down"></i>
                </button>
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
@import '../../assets/skeleton.css';
</style>

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
  data(){
    return {
      carts: [],
      totalCart: 0,
      notifications: [],
      limitData: 5,
      totalData: 0,
      isLoading: true,
      totalNotification: 0
    }
  },
  mounted() {
    if(authHeader()["Authorization"] === undefined) {
      this.$router.push({ name: 'auth.login' })
    }
    this.fetchData()
    this.fetchNotification()
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
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
    },
    async fetchNotification() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/notifications/user`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalNotification = response.data.data.filter(e => e.is_read === false).length
          this.totalData = response.data.data.length;
          this.notifications = response.data.data.slice(0, this.limitData);
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })

      this.isLoading = false
    },
    loadMore(){
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/notifications`,{ headers: authHeader() }).then((response) => {
        this.notifications = response.data.data.slice(0, this.limitData += 5);
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });
    },
    mark(id){
      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/notifications/mark/${id}`,null,{ headers: authHeader() }).then((response) => {
        if(response.data.code === 200) {
          console.log(response.data.status)
          this.fetchNotification()
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          alert(error.response.data.status)
        }
      })
    },
    markAll(){
      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/notifications/mark`,null,{ headers: authHeader() }).then((response) => {
        if(response.data.code === 200) {
          console.log(response.data.status)
          this.fetchNotification()
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          alert(error.response.data.status)
        }
      })
    },
    getColor(classNames, status) {
      let className = classNames
      if (status === "settlement") {
        className += "success"
      } else if (status === "pending") {
        className += "info"
      } else {
        className += "danger"
      }

      return className
    }
  }
}
</script>