<script>
import HomeIcon from './HomeIcon.vue';
import ProfileIcon from './ProfileIcon.vue';
import SearchIcon from './SearchIcon.vue';

export default {
  name: 'AppHeader',
  components: {HomeIcon, SearchIcon, ProfileIcon},
  data: function() {
    return {
      loading: false,
      errorMsg: null,
      id: localStorage.getItem('id'),
      usernameSearch: '',
      users: [],
    };
  },
  methods: {
    logout() {
      localStorage.clear();
      this.$router.replace('/login');
    },
    searchUser() {
      this.$router.push({path: '/users', query: {username: this.usernameSearch}});
      this.usernameSearch = '';
    },
  },
};
</script>


<template>
  <div class="header-ctn section" v-if="this.id">
    <div class="inner">
      <router-link class="logo" to="/home">
        <img src="@/assets/bitmaps/logo.png" alt="">
      </router-link>
      <div class="search-bar">
        <form @submit.prevent="this.searchUser()">
          <input type="text" placeholder="Search" v-model="usernameSearch"  />
          <SearchIcon class="search-icon" />
        </form>
      </div>
      <div class="quick-actions">
        <router-link to="/home">
          <HomeIcon :active="$route.name === 'home'" />
        </router-link>
        <router-link :to="'/users/' + this.id">
          <ProfileIcon :active="this.$route.path == '/users/' + this.id" />
        </router-link>
        <button class="header-btn" @click="logout">
            Logout
        </button>
      </div>
    </div>
  </div>
  <div class="header-ctn section" v-if="!this.id">
    <div class="inner-alone">
      <router-link class="logo-alone" to="/login">
        <img src="@/assets/bitmaps/logo.png">
      </router-link>
    </div>
  </div>
</template>


<style scoped>
.header-ctn {
  border-bottom: 1px solid var(--border-color);
  position: fixed;
  width: 100%;
  background: #FFF;
  z-index: 1;
  top: 0;
}
.header-ctn > .inner {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.header-ctn > .inner-alone {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.header-ctn > .inner > * {
  flex-basis: 30%;
}
.header-ctn .logo {
  height: 24px;
}
.header-ctn .logo-alone {
  height: 24px;
}
.search-bar {
  position: relative;
}
.search-bar input {
  background: rgba(239,239,239,1);
  border: none;
  outline: none;
  height: 36px;
  width: 270px;
  max-width: 230px;
  padding: 0 8px;
  padding-left: 46px;
  border-radius: 8px;
  color: rgba(142,142,142,1);
  font-size: 1rem;
}
.search-bar input::placeholder {
  color: rgba(142,142,142,1);
  font-weight: 300;
}
.search-bar .search-icon {
  position: absolute;
  left: 16px;
  top: 10px;
}
.quick-actions {
  display: flex;
  align-items: center;
}
.quick-actions > * {
  margin-left: 22px;
  position: relative;
  text-align: left;
}
.header-btn {
  text-decoration: none;
  height: 30px;
  width: 85px;
  display: grid;
  place-items: center;
  border: 1px solid var(--border-color);
  font-size: 0.875rem;
  color: #000000;
  border-radius: 4px;
  font-weight: 600;
}
</style>
