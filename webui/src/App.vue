<script>
import AppHeader from './components/AppHeader.vue';

export default {
  name: 'App',
  components: {AppHeader},
};
</script>


<template>
<div class="app-ctn">
  <AppHeader :key="this.$route.fullPath" />
  <RouterView class="mt-60" :key="this.$route.fullPath" />
</div>
</template>


<style>
* {
  font-family: 'Segoe UI', 'Open Sans', 'Helvetica Neue', sans-serif;
}
:root {
  --border-color: rgba(219,219,219,1);
  --grey-text: #aaaaaa;
  --blue-text: #0095F6;
}
.section {
  display: flex;
  justify-content: center;
}
.section.page {
  margin: 28px 0;
  margin-bottom: 0;
}
.section > .inner {
  width: 935px;
  position: relative;
}
button {
  cursor: pointer;
}
.clear-btn {
  text-decoration: none;
  outline: none;
  background: none;
  border: none;
  padding: 0;
  font-size: 0.875rem;
  color: var(--grey-text);
}
.drop-modal {
  box-shadow: 0 0 5px 3px rgb(0 0 0 / 10%);
  margin-top: 8px;
}
.triangle-pointer {
  content: '';
  width: 20px;
  height: 20px;
  background: #FFFFFF;
  box-shadow: 0 0 5px 3px rgb(0 0 0 / 10%);
  transform: rotate(45deg);
  position: absolute;
  top: -5px;
  right: 52px;
  z-index: -1;
}
.tp-cover {
  content: '';
  width: 20px;
  height: 20px;
  background: #FFFFFF;
  transform: rotate(45deg);
  position: absolute;
  top: -5px;
  right: 52px;
}
.bg-layer {
  position: fixed;
  inset: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
}

.fade {
  opacity: 0.3;
}
.mt-60 {
  padding-top: 60px;
}
.mt-0 {
  margin-top: 0px;
}
.reveals {
  animation: reveals .4s ease-in-out;
  -webkit-animation: reveals .4s ease-in-out;
}
@keyframes reveals {
  0% {opacity: 0;transform: scale(0.6);}
  60% {opacity: 1;transform: scale(1);}
  80% {opacity: 1;transform: scale(0.9);}
  100% {opacity: 1;transform: scale(1);}
}
</style>
