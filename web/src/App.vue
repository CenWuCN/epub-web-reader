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
#win {
  min-width: 100%;
  /* min-height: 100%; */
  display: grid;
  grid-template-columns: 2fr 6fr 2fr;
  grid-template-rows: 100vh;
}
.blank{
  background-color: aqua;
}
#content {
  background-color: lightblue;
}
</style>