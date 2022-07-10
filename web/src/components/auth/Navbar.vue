<template>
  <nav id="navbar-main" class="navbar navbar-horizontal navbar-transparent navbar-main navbar-expand-lg navbar-light">
    <div class="container">
      <router-link class="navbar-brand" :to="{ name: 'guest.index' }" style="font-family: 'Fugaz One', cursive;">
        RYZYOLSHOP
      </router-link>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbar-collapse" aria-controls="navbar-collapse" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="navbar-collapse navbar-custom-collapse collapse" id="navbar-collapse">
        <div class="navbar-collapse-header">
          <div class="row">
            <div class="col-6 collapse-brand">
              <router-link :to="{ name: 'auth.login' }" style="font-family: 'Fugaz One', cursive;">
                RYZYOLSHOP
              </router-link>
            </div>
            <div class="col-6 collapse-close">
              <button type="button" class="navbar-toggler" data-toggle="collapse" data-target="#navbar-collapse" aria-controls="navbar-collapse" aria-expanded="false" aria-label="Toggle navigation">
                <span></span>
                <span></span>
              </button>
            </div>
          </div>
        </div>
        <ul class="navbar-nav ml-auto">
          <li class="nav-item">
            <router-link :to="{ name: 'auth.login' }" class="nav-link">
              <span class="nav-link-inner--text">Login</span>
            </router-link>
          </li>
          <li class="nav-item">
            <router-link :to="{ name: 'auth.register' }" class="nav-link">
              <span class="nav-link-inner--text">Register</span>
            </router-link>
          </li>
          <li class="nav-item">
            <router-link :to="{ name: 'auth.forgot-password' }" class="nav-link">
              <span class="nav-link-inner--text">Reset Password</span>
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Fugaz+One&display=swap');

.navbar-brand {
  font-size: 1.50rem;
}
</style>