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
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile" id="print">
            <!-- Card header -->
            <div class="card-header justify-content-between d-flex">
              <h3 class="mb-0">Rincian Pesanan</h3>
              <h3 class="mb-0">{{ payment.transaction_status.toUpperCase() }}</h3>
            </div>
            <!-- Card body -->
            <div class="card-body">
              <h3>Nota Pesanan</h3>
              <div class="bg-secondary p-3 rounded">
                <p>Nama Pembeli : {{ user.name }}</p>
                <p>Alamat Pembeli : {{ user.address }}</p>
                <p>No. Handphone Pembeli : {{ user.phone }}</p>
                <p>Nama Toko Penjual : Toko Ryzy Olshop</p>
              </div>
              <div class="row mt-3">
                <div class="col-6 col-lg-3">
                  <p>No. Pesanan <br> {{ payment.order_id }} </p>
                </div>
                <div class="col-6 col-lg-3">
                  <p>Waktu Pembayaran <br> {{ payment.transaction_time }} </p>
                </div>
                <div class="col-6 col-lg-3">
                  <p>Metode Pembayaran <br> {{ payment.bank_type }} </p>
                </div>
                <div class="col-6 col-lg-3">
                  <p>Nomor Virtual Account <br> {{ payment.va_number == "" ? (payment.biller_code == "" ? "" : `Biller Code : ${payment.biller_code} | Bill Key : ${payment.bill_key}`) : payment.va_number }} </p>
                </div>
              </div>
              <h3 class="mt-5">Rincian Pesanan</h3>
              <div class="table-responsive">
                <table class="table table-bordered">
                  <thead class="thead-light">
                  <tr class="text-center">
                    <th>Kuantitas</th>
                    <th>Produk</th>
                    <th>Harga Produk</th>
                    <th>Subtotal</th>
                  </tr>
                  </thead>
                  <tbody>
                  <tr v-for="(item, i) in orders" :key="i">
                    <td class="text-center">{{ item.quantity }}</td>
                    <td>{{ item.name }}</td>
                    <td class="text-right">Rp {{ numberWithCommas(item.price) }}</td>
                    <td class="text-right">Rp {{ numberWithCommas(item.total_price_item) }}</td>
                  </tr>
                  <tr class="font-weight-bold">
                    <td class="text-right" colspan="3">Subtotal</td>
                    <td class="text-right">Rp {{ numberWithCommas(totalPrice - 5000) }}</td>
                  </tr>
                  <tr class="font-weight-bold">
                    <td class="text-right" colspan="3">Pajak Aplikasi</td>
                    <td class="text-right">Rp {{ numberWithCommas(5000) }}</td>
                  </tr>
                  <tr class="font-weight-bold">
                    <td class="text-right" colspan="3">Total Pembayaran</td>
                    <td class="text-right">Rp {{ numberWithCommas(totalPrice) }}</td>
                  </tr>
                  </tbody>
                </table>
              </div>
              <div class="justify-content-between d-flex" id="information">
                <div>
                </div>
                <div class="mt-5">
                  <a href="javascript:void(0);" @click="print">
                    <i class="fa fa-print fa-2x text-dark"></i>
                  </a>
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
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      carts: [],
      totalCart: 0,
      orders: [],
      totalPrice: 5000,
      user: {
        name: "",
        address: "",
        phone: ""
      },
      payment: {
        id_payload: "",
        user_id: "",
        order_id: "",
        transaction_time: "",
        transaction_status: "",
        transaction_id: "",
        status_code: "",
        signature_key: "",
        settlement_time: "",
        payment_type: "",
        merchant_id: "",
        gross_amount: "",
        fraud_status: "",
        bank_type: "",
        va_number: "",
        biller_code: "",
        bill_key: ""
      }
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

      axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order/${this.$route.params.order_id}`, { headers: authHeader() }).then(response => {
        this.orders = response.data.data
        if(response.data.data != null) {
          response.data.data.map(item => {
            this.totalPrice += item.total_price_item;
          })
        }

        axios.get(`${process.env.VUE_APP_SERVICE_URL}/payments/${this.$route.params.order_id}`, { headers: authHeader() }).then(response => {
          axios.get(`${process.env.VUE_APP_SERVICE_URL}/users/${response.data.data.user_id}`, { headers: authHeader() }).then(response => {
            this.user = {
              name: response.data.data.name,
              address: response.data.data.address,
              phone: response.data.data.phone
            }
          })
        }).catch(error => {
          console.log(error)
        })
      })

      axios.get(`${process.env.VUE_APP_SERVICE_URL}/payments/${this.$route.params.order_id}`, { headers: authHeader() }).then(response => {
        this.payment = {
          id_payload: response.data.data.id_payload,
          user_id: response.data.data.user_id,
          order_id: response.data.data.order_id,
          transaction_time: response.data.data.transaction_time,
          transaction_status: response.data.data.transaction_status,
          transaction_id: response.data.data.transaction_id,
          status_code: response.data.data.status_code,
          signature_key: response.data.data.signature_key,
          settlement_time: response.data.data.settlement_time,
          payment_type: response.data.data.payment_type,
          merchant_id: response.data.data.merchant_id,
          gross_amount: response.data.data.gross_amount,
          fraud_status: response.data.data.fraud_status,
          bank_type: response.data.data.bank_type,
          va_number: response.data.data.va_number,
          biller_code: response.data.data.biller_code,
          bill_key: response.data.data.bill_key
        }
      })
    },
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    },
    print(){
      const except = document.getElementById("information");
      except.style.display = "none";
      const prtHtml = document.getElementById('print').innerHTML;

      let stylesHtml = '';
      for (const node of [...document.querySelectorAll('link[rel="stylesheet"], style')]) {
        stylesHtml += node.outerHTML;
      }

      const WinPrint = window.open('', '', 'left=0,top=0,width=800,height=900,toolbar=0,scrollbars=0,status=0');

      WinPrint.document.write(`<!DOCTYPE html>
        <html>
          <head>
            ${stylesHtml}
          </head>
          <body>
            ${prtHtml}
          </body>
        </html>`);

      WinPrint.document.close();
      WinPrint.focus();
      WinPrint.print();
      WinPrint.close();
    }
  }
}
</script>