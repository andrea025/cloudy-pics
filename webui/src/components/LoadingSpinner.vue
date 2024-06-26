<script>
export default {
  name: 'LoadingSpinner',
  props: {
    loading: Boolean,
  },
};
</script>

<template>
	<div v-if="loading" class="loader-wrapper">
	  <div class="loader-class">
		<div></div>
		<div></div>
		<div></div>
		<div></div>
	  </div>
	  <p class="loader-text">Wait a sec...</p>
	</div>
</template>


<style scoped>
.loader-text {
  font-size: 2em;
  font-weight: bolder;
  margin-top: 30px;
  letter-spacing: 1px;
  color: rgb(82, 82, 82);
}
.loader-wrapper {
  text-align: center;
}
.loader-class {
  display: inline-block;
  position: relative;
  width: 200px;
  height: 200px;
}
.loader-class div {
  box-sizing: border-box;
  display: block;
  position: absolute;
  width: 200px;
  height: 200px;
  margin: 8px;
  border: 20px solid #e1306c;
  border-radius: 50%;
  animation: loader 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  border-color: #e1306c transparent transparent transparent;
}
.loader-class div:nth-child(1) {
  animation-delay: -0.45s;
}
.loader-class div:nth-child(2) {
  animation-delay: -0.3s;
}
.loader-class div:nth-child(3) {
  animation-delay: -0.15s;
}
@keyframes loader {
  0% {
	transform: rotate(0deg);
  }
  100% {
	transform: rotate(360deg);
  }
}
</style>
