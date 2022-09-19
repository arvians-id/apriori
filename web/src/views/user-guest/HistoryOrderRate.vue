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
              <div class="card-header">
                <h3 class="mb-0">Penilaian Produk</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <div class="loading-skeleton" v-if="isLoading">
                  <h3 class="font-weight-bold text-center" v-if="isExists == false">Penilaian produk untuk : {{ order.name }}</h3>
                  <h3 class="font-weight-bold text-center" v-else>Anda sudah memberi nilai pada produk ini</h3>
                  <form @submit.prevent="submit" method="POST">
                    <div class="form-group rating">
                      <p style="padding-bottom: 130px; padding-left: 500px">This is for rating</p>
                    </div>
                    <div class="form-group">
                      <label class="form-control-label">this is label for input</label>
                      <p style="padding-bottom: 70px; padding-left: 500px">This is for rating</p>
                    </div>
                    <div class="form-group">
                      <label class="form-control-label">this is label for input</label>
                      <p style="padding-bottom: 110px; padding-left: 500px">This is for rating</p>
                    </div>
                    <button class="btn btn-primary" type="submit" :disabled="isExists">Submit form</button>
                  </form>
                </div>
                <div v-else>
                  <h3 class="font-weight-bold text-center" v-if="isExists == false">Penilaian produk untuk : {{ order.name }}</h3>
                  <h3 class="font-weight-bold text-center" v-else>Anda sudah memberi nilai pada produk ini</h3>
                  <form @submit.prevent="submit" method="POST">
                    <div class="form-group rating">
                      <input type="radio" name="rating" v-model="comment.rating" value="5" id="5" required><label for="5">☆</label>
                      <input type="radio" name="rating" v-model="comment.rating" value="4" id="4" required><label for="4">☆</label>
                      <input type="radio" name="rating" v-model="comment.rating" value="3" id="3" required><label for="3">☆</label>
                      <input type="radio" name="rating" v-model="comment.rating" value="2" id="2" required><label for="2">☆</label>
                      <input type="radio" name="rating" v-model="comment.rating" value="1" id="1" required><label for="1">☆</label>
                    </div>
                    <div class="form-group">
                      <label class="form-control-label">Tag</label>
                      <select class="form-control" v-model="comment.tag" name="tag" multiple>
                        <option value="Kualitas Barang">Kualitas Barang</option>
                        <option value="Pelayanan Penjual">Pelayanan Penjual</option>
                        <option value="Harga Barang">Harga Barang</option>
                        <option value="Pengiriman">Pengiriman</option>
                      </select>
                    </div>
                    <div class="form-group">
                      <label class="form-control-label">Komentar</label>
                      <textarea type="text" class="form-control" name="description" v-model="comment.description" rows="5"></textarea>
                    </div>
                    <button class="btn btn-primary" type="submit" :disabled="isExists">Submit form</button>
                  </form>
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
@import '../../assets/skeleton.css';

.rating {
  display: flex;
  flex-direction: row-reverse;
  justify-content: center;
}

.rating > input{ display:none;}

.rating > label {
  position: relative;
  width: 1em;
  font-size: 6vw;
  color: #FFD600;
  cursor: pointer;
}
.rating > label::before{
  content: "\2605";
  position: absolute;
  opacity: 0;
}
.rating > label:hover:before,
.rating > label:hover ~ label:before {
  opacity: 1 !important;
}

.rating > input:checked ~ label:before{
  opacity:1;
}

.rating:hover > input:checked ~ label:before{ opacity: 0.4; }

body{ background: #222225; color: white;}
h1, p{
  text-align: center;

}

h1{
  margin-top:150px;
}
p{ font-size: 1.2rem;}
@media only screen and (max-width: 600px) {
  h1{font-size: 14px;}
  p{font-size: 12px;}
}

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
  mounted() {
    this.fetchData()
    this.fetchOrder()
    if(authHeader()["Authorization"] !== undefined) {
      this.fetchNotification()
    }
  },
  data: function () {
    return {
      comment: {
        user_order_id: parseInt(this.$route.params.id_order),
        product_code: "",
        description: "",
        rating: 0,
        tag: [],
      },
      order: [],
      isExists: false,
      isLoading: true,
      carts: [],
      totalCart: 0,
      totalNotification: 0,
      notifications: []
    };
  },
  methods: {
    async fetchData(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order/single/${this.$route.params.id_order}`, { headers: authHeader() })
          .then(response => {
            this.order = response.data.data;
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              console.log(error.response.data.status)
            }
          })

      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }
    },
    async fetchOrder(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/comments/user-order/${this.$route.params.id_order}`, { headers: authHeader() })
          .then(response => {
            if(response.data.data != null){
              this.comment = {
                user_order_id: response.data.data.user_order_id,
                product_code: response.data.data.product_code,
                description: response.data.data.description,
                rating: response.data.data.rating,
                tag: response.data.data.tag.split(', '),
              };
              this.isExists = true;
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              console.log(error.response.data.status)
            }
          })

      this.isLoading = false
    },
    async fetchNotification() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/notifications/user`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalNotification = response.data.data.filter(e => e.is_read === false).length
          this.notifications = response.data.data
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })
    },
    submit() {
      this.comment.product_code = this.order.code
      if (this.comment.tag.length > 0) {
        let tag = this.comment.tag
        this.comment.tag = tag.join(", ")
      } else {
        this.comment.tag = ""
      }
      this.comment.rating = parseInt(this.comment.rating)

      axios.post(`${process.env.VUE_APP_SERVICE_URL}/comments`, this.comment, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'member.history'
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
          })
    }
  }
}
</script>