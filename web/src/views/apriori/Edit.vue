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
                <h3 class="mb-0">Edit Apriori</h3>
              </div>
              <!-- Card body -->
              <div class="card-body" v-if="isLoading">
                <p class="mt-2 text-center">Loading...</p>
              </div>
              <div class="card-body" v-else>
                <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Description</label>
                    <textarea v-model="apriori.description" class="form-control" rows="5"></textarea>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Image</label>
                    <input type="file" class="form-control" @change="uploadImage">
                  </div>
                  <div class="form-group">
                    <img :src="previewImage" width="150"/>
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
      apriori: {
        description: "",
        image: null
      },
      previewImage: "https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/no-image.png",
      isLoading: true
    };
  },
  methods: {
    submit() {
      const config = {
        headers: {
          'Content-Type': 'multipart/form-data',
          ...authHeader(),
        }
      }
      let formData = new FormData()
      formData.append("description", this.apriori.description)
      formData.append("image", this.apriori.image)

      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/update/${this.$route.params.id}`, formData, config)
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'apriori.code-detail',
                params: {
                  code: this.$route.params.code,
                  id: this.$route.params.id
                }
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
      })
    },
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then(response => {
        this.apriori = {
          description: response.data.data.apriori_description,
        }
        this.previewImage = response.data.data.apriori_image
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading = false
    },
    uploadImage(e) {
      let files = e.target.files[0]
      this.apriori.image = files
      this.previewImage = URL.createObjectURL(files)
    }
  }
}
</script>