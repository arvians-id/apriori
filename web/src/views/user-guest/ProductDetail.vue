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
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile">
            <!-- Card header -->
            <div class="card-body mt-5">
              <div v-if="isLoading">
                <div class="loading-skeleton row d-flex justify-content-center">
                  <div class="col-12 col-lg-3 mb-2">
                    <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/no-image.png" class="img-fluid mb-2">
                    <div class="border p-3 rounded text-center" style="color: #525f7f">
                      <p class="mb-0">Atur jumlah yang pembelian</p>
                      <div class="mt-2">
                        <button type="button" class="btn btn-danger btn-sm">-</button>
                        <button class="btn disabled">quantity</button>
                        <button type="button" class="btn btn-primary btn-sm">+</button>
                      </div>
                    </div>
                    <div class="card mt-4">
                      <div class="card-body pt-0">
                        <div class="d-flex align-items-center justify-content-center p-0 m-0 mt-3">
                          <p>Star</p>
                          <p class="text-center d-inline font-weight-bold ml-2" style="font-size: 68px">4.8</p>
                          <p class="ml-1">/5.0</p>
                        </div>
                        <div class="text-center">
                          <p class="font-weight-bold mb-1">94% pembeli merasa puas</p>
                          <p>899 rating • 289 ulasan</p>
                        </div>
                        <p class="py-5">Star</p>
                      </div>
                    </div>
                    <div class="card mt-4">
                      <div class="card-body">
                        <h3 class="card-title text-uppercase">Filter Ulasan</h3>
                        <p class="mb-2 font-weight-bold">Rating</p>
                        <p class="py-5">Star</p>
                        <hr class="my-3">
                        <p class="mb-2 font-weight-bold">Topik Ulasan</p>
                        <p class="py-5">Star</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-lg-6">
                    <div class="text-left">
                      <h5 class="h2 p-0 m-0 mb-1">This is title of product</h5>
                      <p class="p-0 m-0 mb-1 w-50">This is title of product</p>
                      <p class="w-50">Pricing</p>
                      <hr class="mt-0 mb-1">
                      <div class="nav-wrapper">
                        <p class="pt-4 mb-0">this is button</p>
                      </div>
                      <div class="card">
                        <div class="card-body">
                          <div class="tab-content" id="skeleton-myTabContent">
                            <div class="tab-pane fade show active" id="skeleton-detail" role="tabpanel" aria-labelledby="detail-tab">
                              <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                              <p>Original baru</p>
                              <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                              <p class="font-weight-bold">Text</p>
                              <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                              <p>gram</p>
                              <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                              <p>Description</p>
                              <hr class="m-0 mb-3">
                              <div class="media">
                                <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/ryzy.jpg" width="59" class="mr-3" alt="...">
                                <div class="media-body">
                                  <p class="mt-0 mb-0 w-50 mb-1">Title</p>
                                  <p>Description</p>
                                </div>
                              </div>
                              <hr class="m-0 mb-3">
                              <p class="font-weight-bold mb-0 mt-2">Pengiriman</p>
                              <p class="mb-1"><i class="ni ni-pin-3"></i> Dikirim dari Tanggerang, Banten</p>
                              <p><i class="ni ni-delivery-fast"></i> Tersedia pengiriman dengan TIKI, JNE dan POS Indonesia</p>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="card">
                        <div class="card-body">
                          <h3 class="card-title text-uppercase">Penilaian Produk</h3>
                          <div class="media" v-for="item in 5" :key="item">
                            <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/ryzy.jpg" width="59" class="mr-3" alt="...">
                            <div class="media-body">
                              <p class="mt-0 mb-0 w-50 mb-1">Title</p>
                              <p class="py-5">Description</p>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-lg-2">
                    <h3 class="mb-0 mt-3">Produk terkait</h3>
                    <hr class="mb-3 p-0">
                    <div class="card card-pricing border shadow-none" v-for="item in 5" :key="item">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body">
                        <p class="card-title m-0 mb-1">this is title of product</p>
                        <p class="card-text p-0 m-0 w-50">pricing</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else>
                <div class="row d-flex justify-content-center">
                  <div class="col-12 col-lg-3 mb-2">
                    <img :src="product.image" class="img-fluid mb-2" width="500">
                    <div class="border p-3 rounded text-center" style="color: #525f7f">
                      <p class="mb-0">Atur jumlah yang pembelian</p>
                      <div>
                        <button type="button" @click="min(product)" class="btn btn-danger btn-sm">-</button>
                        <button class="btn disabled">{{ quantity }} item</button>
                        <button type="button" @click="add(product)" class="btn btn-primary btn-sm">+</button>
                      </div>
                    </div>
                    <div class="card mt-4">
                      <div class="card-body pt-0" v-if="mainRatings > 0">
                        <div class="d-flex align-items-center justify-content-center p-0 m-0">
                          <i class="fas fa-star text-warning fa-lg"></i>
                          <p class="text-center d-inline font-weight-bold ml-2" style="font-size: 68px">{{ mainRatings.toFixed(1) }}</p>
                          <p class="ml-1" style="padding-top: 38px">/5.0</p>
                        </div>
                        <div class="text-center">
                          <p class="font-weight-bold mb-1">{{ satisfactionRating.toFixed() }}% pembeli merasa puas</p>
                          <p>{{ totalRatings }} rating • {{ ratings.reduce((a, b) => a + b.result_comment, 0) }} ulasan</p>
                        </div>
                        <div class="row d-flex align-items-center" v-for="item in 5" :key="item">
                          <div class="col-4 text-right">
                            <i class="fas fa-star text-warning"></i>
                            {{ 5 - item + 1 }}
                          </div>
                          <div class="col-6">
                            <div class="progress mb-0">
                              <div
                                  class="progress-bar bg-orange"
                                  role="progressbar"
                                  aria-valuenow=""
                                  aria-valuemin="0"
                                  aria-valuemax="100"
                                  :style="`width: ${ratings.filter(x => x.rating === 5 - item + 1).length > 0 ? (ratings.filter(x => x.rating === 5 - item + 1)[0].result_rating / (5 - item + 1)) / totalRatings * 100 : 0}%;`"></div>
                            </div>
                          </div>
                          <div class="col-2 px-0">
                            {{ ratings.filter(x => x.rating === 5 - item + 1).length > 0 ? ratings.filter(x => x.rating === 5 - item + 1)[0].result_rating / (5 - item + 1) : 0 }}
                          </div>
                        </div>
                      </div>
                      <div class="card-body pt-0" v-else>
                        <p class="mt-4 font-weight-bold text-center mb-0">Produk ini belum memiliki rating</p>
                      </div>
                    </div>
                    <div class="card mt-4">
                      <div class="card-body">
                        <h3 class="card-title text-uppercase">Filter Ulasan</h3>
                        <p class="mb-2 font-weight-bold">Rating</p>
                        <div class="form-check mb-2" v-for="item in 5" :key="item">
                          <input class="form-check-input" type="radio" :value="5 - item + 1" v-model="filter.rating" :id="'radio-' +item" @change="filterComment">
                          <label class="form-check-label" :for="'radio-' +item">
                            <i class="fas fa-star text-warning"></i> {{ 5 - item + 1 }}
                          </label>
                        </div>
                        <hr class="my-3">
                        <p class="mb-2 font-weight-bold">Topik Ulasan</p>
                          <div class="form-check mb-2" v-for="(item, i) in ['Kualitas Barang', 'Pelayanan Penjual', 'Harga Barang', 'Pengiriman']" :key="i">
                            <input class="form-check-input" type="checkbox" v-model="filter.tags" :value="item" :id="'checkbox-' + i" @change="filterComment">
                            <label class="form-check-label" :for="'checkbox-' + i">
                              {{ item }}
                            </label>
                          </div>
                        <button class="btn btn-warning btn-sm" @click="fetchComments">Reset Filter</button>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-lg-6">
                    <div class="text-left">
                      <h5 class="h2 text-uppercase p-0 m-0">{{ product.name }}</h5>
                      <p class="p-0 m-0">code : {{ product.code }}</p>
                      <div class="h1 font-weight-bold">Rp. {{ product.price }}</div>
                      <hr class="mt-0 mb-1">
                      <div class="nav-wrapper">
                        <ul class="nav nav-pills nav-fill flex-column flex-md-row" id="tabs-icons-text" role="tablist">
                          <li class="nav-item">
                            <a class="nav-link mb-sm-3 mb-md-0 active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true"><i class="ni ni-tag mr-2"></i>Detail</a>
                          </li>
                          <li class="nav-item">
                            <a class="nav-link mb-sm-3 mb-md-0" id="lainnya-tab" data-toggle="tab" href="#lainnya" role="tab" aria-controls="lainnya" aria-selected="false"><i class="ni ni-ungroup mr-2"></i>Lainnya</a>
                          </li>
                        </ul>
                      </div>
                      <div class="card shadow">
                        <div class="card-body">
                          <div class="tab-content" id="myTabContent">
                            <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                              <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                              <p>Original baru</p>
                              <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                              <p>
                                <router-link :to="{ name: 'guest.product', query: { category: category } }" v-for="(category, i) in product.category.split(', ')" :key="i" class="text-primary font-weight-bold">
                                  {{ category }}{{ product.category.split(', ').length - 1 != i ? ', ' : '' }}
                                </router-link>
                              </p>
                              <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                              <p>{{ product.mass }} gram</p>
                              <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                              <p>{{ product.description == "" ? "Tidak ada deskripsi" : product.description }}</p>
                              <hr class="m-0 mb-3">
                              <div class="media">
                                <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/ryzy.jpg" width="53" class="mr-3" alt="...">
                                <div class="media-body">
                                  <h3 class="mt-0 mb-0">Toko Ryzy Olshop</h3>
                                  <p>Produk Original Berkualitas dan Terpercaya..</p>
                                </div>
                              </div>
                              <hr class="m-0 mb-3">
                              <p class="font-weight-bold mb-0 mt-2">Pengiriman</p>
                              <p class="mb-1"><i class="ni ni-pin-3"></i> Dikirim dari Tanggerang, Banten</p>
                              <p><i class="ni ni-delivery-fast"></i> Tersedia pengiriman dengan TIKI, JNE dan POS Indonesia</p>
                            </div>
                            <div class="tab-pane fade" id="lainnya" role="tabpanel" aria-labelledby="lainnya-tab">
                              <p class="font-weight-bold mb-0 mt-2">Tanggal Dibuat</p>
                              <p>{{ product.created_at }}</p>
                              <p class="font-weight-bold mb-0 mt-2">Terakhir Diubah</p>
                              <p>{{ product.updated_at }}</p>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="card">
                      <div class="card-body">
                        <h3 class="card-title text-uppercase">Penilaian Produk</h3>
                        <ul class="list-unstyled" v-if="comments.length > 0">
                          <li class="media my-4" style="border-bottom: 1px solid #e9ecef" v-for="item in comments" :key="item.id_comment">
                            <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/user.png" width="53" class="mr-3" alt="...">
                            <div class="media-body">
                              <h4 class="mt-0 mb-0">{{ item.user_order.payment.user.name }}</h4>
                              <small>{{ item.created_at }}</small>
                              <div class="mb-2">
                                <i class="fas fa-star text-warning" v-for="itemRating in item.rating" :key="itemRating"></i>
                              </div>
                              <template v-if="item.tag !== ''">
                                <button class="btn btn-dark btn-sm disabled mb-2" v-for="(item, i) in item.tag.split(',')" :key="i">{{ item }}</button>
                              </template>
                              <p>{{ item.description }}</p>
                            </div>
                          </li>
                        </ul>
                        <div v-else>
                          <div class="alert alert-secondary mt-3">
                            <h5 class="alert-heading">Oops!</h5>
                            <p>Tidak ada penilaian.</p>
                          </div>
                        </div>
                        <button @click="loadMore()" v-if="comments.length !== this.totalData" class="btn btn-secondary d-block mx-auto px-5">
                          Lihat lainnya <i class="ni ni-bold-down"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-lg-2">
                    <h3 class="mb-0 mt-3 text-uppercase">Produk terkait</h3>
                    <hr class="mb-3 p-0">
                    <div v-if="productSimilarCategory.length > 0">
                      <div class="card card-pricing border shadow-none" v-for="item in productSimilarCategory" :key="item.id_product">
                        <div class="embed-responsive embed-responsive-16by9">
                          <img class="card-img-top embed-responsive-item" :src="item.image" alt="Preview Image">
                        </div>
                        <div class="card-body">
                          <router-link :to="{ name: 'guest.product.detail', params: { code: item.code } }" class="card-title m-0">{{ item.name }}</router-link>
                          <p class="card-text p-0 m-0">Rp. {{ item.price }}</p>
                        </div>
                      </div>
                    </div>
                    <div v-else>
                      <div class="alert alert-secondary mt-3">
                        <h5 class="alert-heading">Oops!</h5>
                        <p>Tidak ada produk yang serupa.</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card">
            <!-- Card header -->
            <div class="card-header">
              <!-- Title -->
              <h5 class="h3 mb-0">Rekomendasi</h5>
            </div>
            <div class="card-body" v-if="isLoading2">
              <div class="loading-skeleton">
                <div class="row">
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="card-body" v-else>
              <div class="row" v-if="recommendation.length > 0">
                <div class="col-12 col-md-6 col-lg-4" v-for="item in recommendation" :key="item.apriori_id">
                  <div class="card card-pricing border-0 mb-4">
                    <div class="embed-responsive embed-responsive-16by9">
                      <img class="card-img-top embed-responsive-item" :src="item.apriori_image" alt="Preview Image">
                    </div>
                    <div class="card-body pb-3">
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="text-dark d-none d-lg-block">
                        Paket Barang {{ item.apriori_item.length > 20 ? UpperWord(item.apriori_item.slice(0, 20)) + "..." : UpperWord(item.apriori_item) }}
                      </router-link>
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="text-dark d-block d-lg-none">
                        Paket Barang {{ UpperWord(item.apriori_item) }}
                      </router-link>
                      <ul class="list-unstyled">
                        <li v-for="(value,i) in item.apriori_item.split(', ')" :key="i">
                          <div class="d-flex align-items-center">
                            <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                              <i class="ni ni-basket"></i>
                            </div>
                            <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                          </div>
                        </li>
                      </ul>
                      <div class="card-footer p-0 pt-2 text-center m-0">
                        <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="font-weight-bold">
                          {{ item.apriori_discount }}%
                        </router-link>
                        <small class="text-muted mb-0">diskon</small>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="row" v-else>
                <div class="col-12">
                  <div class="alert alert-secondary">
                    <h5 class="alert-heading">Oops!</h5>
                    <p>Tidak ada rekomendasi yang tersedia.</p>
                  </div>
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

.card-img-top {
  width: 100%;
  object-fit: cover;
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
  watch: {
    '$route': function (to) {
      if (to.name === 'guest.product.detail') {
        this.allFetch()
      }
    }
  },
  mounted() {
    this.allFetch()
  },
  data: function () {
    return {
      product: [],
      recommendation: [],
      productSimilarCategory: [],
      carts: [],
      totalCart: 0,
      quantity: 0,
      isLoading: true,
      isLoading2: true,
      categories: [],
      allComments: [],
      comments: [],
      limitData: 4,
      totalData: 0,
      mainRatings: 0,
      ratings: [],
      totalRatings: 0,
      satisfactionRating: 0,
      filter: {
        tags: [],
        rating: ''
      },
      totalNotification: 0,
      notifications: []
    };
  },
  methods: {
    allFetch(){
      this.fetchCategories()
      this.fetchData()
      this.fetchDataRecommendation()
      this.fetchSimilarCategory()
      this.fetchComments()
      this.fetchRatings()
      if(authHeader()["Authorization"] !== undefined) {
        this.fetchNotification()
      }
      document.getElementsByTagName("body")[0].classList.remove("bg-default");
    },
    async filterComment(){
      let tag = this.filter.tags.join(",")
      let rating = this.filter.rating
      let filter = ""
      if(tag !== "") {
        filter = "?tags=" + tag
        if (rating !== "") {
          filter += "&rating="+rating
        }
      } else if (rating !== "") {
        filter = "?rating=" + rating
        if (tag !== "") {
          filter += "&tags="+tag
        }
      }

      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/comments/product/${this.$route.params.code}${filter}`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length
          this.allComments = response.data.data
          this.comments = response.data.data.slice(0, this.limitData)
        }else {
          this.totalData = 0
          this.allComments = []
          this.comments = []
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })
    },
    async fetchRatings(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/comments/rating/${this.$route.params.code}`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.ratings = response.data.data
          this.totalRatings = response.data.data.reduce((a, b) => a + b.result_rating / b.rating, 0)
          this.mainRatings = response.data.data.reduce((a, b) => a + b.result_rating , 0) / this.totalRatings
          this.satisfactionRating = this.ratings.filter(x => x.rating >= 4).reduce((a, b) => a + (b.result_rating / b.rating) , 0) / this.totalRatings * 100
        } else {
          this.ratings = []
          this.totalRatings = 0
          this.mainRatings = 0
          this.satisfactionRating = 0
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })
    },
    async fetchComments(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/comments/product/${this.$route.params.code}`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          console.log(response.data.data)
          this.totalData = response.data.data.length
          this.allComments = response.data.data
          this.comments = response.data.data.slice(0, this.limitData)
        }else {
          this.totalData = 0
          this.allComments = []
          this.comments = []
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })

      this.filter = {
        tags: [],
        rating: ''
      }
    },
    async fetchCategories(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/categories`, { headers: authHeader() })
        .then(response => {
          this.categories = response.data.data
        }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              console.log(error.response.data.status)
            }
        })
    },
    async fetchSimilarCategory(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}/category`, { headers: authHeader() })
          .then(response => {
            if(response.data.data != null) {
              this.productSimilarCategory = response.data.data.slice(0, 3)
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              console.log(error.response.data.status)
            }
          })
    },
    async fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      let productItem = this.carts.find(product => product.code === this.$route.params.code);
      this.quantity = productItem ? productItem.quantity : 0;

      this.isLoading = false;
    },
    async fetchDataRecommendation() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}/recommendation`, { headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.recommendation = response.data.data;
        } else {
          this.recommendation = []
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading2 = false
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
      });
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    },
    add(item) {
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem) {
        productItem.quantity += 1
        productItem.totalPricePerItem = productItem.price * productItem.quantity
        this.quantity += 1
        this.totalCart += 1
      } else {
        this.carts.push({
          id_product: item.id_product,
          code: item.code,
          name: item.name,
          price: item.price,
          mass: item.mass,
          image: item.image,
          quantity: 1,
          totalPricePerItem: item.price
        });
        this.quantity = 1
        this.totalCart += 1
      }

      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    min(item){
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem !== undefined) {
        if (productItem.quantity > 1) {
          productItem.quantity -= 1
          productItem.totalPricePerItem -= productItem.price
          this.quantity -= 1
          this.totalCart -= 1
        } else {
          this.quantity = 0
          this.totalCart -= 1
          this.carts.splice(this.carts.indexOf(productItem), 1);
        }
        localStorage.setItem('my-carts', JSON.stringify(this.carts));
      }
    },
    loadMore(){
      this.comments = this.allComments.slice(0, this.limitData += 4);
    },
  }
}
</script>