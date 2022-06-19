<template>
  <nav id="navbar-main" class="navbar navbar-horizontal navbar-transparent navbar-main navbar-expand-lg navbar-light">
    <div class="container">
      <router-link class="navbar-brand" :to="{ name: 'login' }">
        <img src="/frontend/img/brand/white.png">
      </router-link>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbar-collapse" aria-controls="navbar-collapse" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="navbar-collapse navbar-custom-collapse collapse" id="navbar-collapse">
        <div class="navbar-collapse-header">
          <div class="row">
            <div class="col-6 collapse-brand">
              <router-link :to="{ name: 'login' }">
                <img src="/frontend/img/brand/blue.png">
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
            <router-link :to="{ name: 'login' }" class="nav-link">
              <span class="nav-link-inner--text">Login</span>
            </router-link>
          </li>
          <li class="nav-item">
            <router-link :to="{ name: 'register' }" class="nav-link">
              <span class="nav-link-inner--text">Register</span>
            </router-link>
          </li>
          <li class="nav-item">
            <router-link :to="{ name: 'forgot-password' }" class="nav-link">
              <span class="nav-link-inner--text">Reset Password</span>
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>