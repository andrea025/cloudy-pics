<script>
export default {
  name: 'ErrorMessage',
  props: ['msg'],
};
</script>


<template>
  <div class="alert alert-danger" role="alert">
    <p>{{ msg }}</p>
  </div>
</template>


<style>
p {
  text-align: center;
  font-weight: bold;
  font-size: larger;
}
</style>
