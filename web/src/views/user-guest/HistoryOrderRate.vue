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
                <h3 class="mb-0">Penilaian Produk</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <h3 class="font-weight-bold text-center">Penilaian produk untuk : {{ order.name }}</h3>
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
                    <select class="form-control" v-model="comment.tag" multiple>
                      <option value="1">Kualitas Barang</option>
                      <option value="2">Pelayanan Penjual</option>
                      <option value="3">Harga Barang</option>
                      <option value="4">Pengiriman</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Komentar</label>
                    <textarea type="text" class="form-control" v-model="comment.name" rows="5"></textarea>
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

<style scoped>
.card-img-top {
  width: 100%;
  object-fit: cover;
}
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
  },
  data: function () {
    return {
      comment: {
        description: "",
        rating: "",
        tag: "",
      },
      order: []
    };
  },
  methods: {
    fetchData(){
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order/${this.$route.params.id_order}/single`, { headers: authHeader() })
          .then(response => {
            this.order = response.data.data;
          }).catch(error => {
            console.log(error.response.data.status)
          })
    },
    submit() {
      // axios.post(`${process.env.VUE_APP_SERVICE_URL}/categories`, this.rateOrder, { headers: authHeader() })
      //     .then(response => {
      //       if(response.data.code === 200) {
      //         alert(response.data.status)
      //         this.$router.push({
      //           name: 'category'
      //         })
      //       }
      //     }).catch(error => {
      //   console.log(error.response.data.status)
      // })
    }
  }
}
</script>