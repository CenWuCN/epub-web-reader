<script setup lang="ts">

</script>

<template>
  <div id="win">
    <div id="leftblank" class="blank"></div>
    <div id="content">
      <RouterView />
    </div>
    <div id="rightblank" class="blank"></div>
  </div>
</template>

<style>
.blank{
  background-color: aqua;
  width: 20%;
}
</style>