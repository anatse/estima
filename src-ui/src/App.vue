<template>
  <div id="app" class="app">
    <router-view></router-view>
  </div>
</template>

<script>
import constGlobal from './constGlobal';
import * as TActions from './store/action-types';

export default {
  name: 'App',
  created() {
    this.$store.dispatch(TActions.USER_CURRENT).catch(() => {
      this.$router.push({ name: constGlobal.PAGE_NAME.LOGIN_PAGE });
    });
  },
};
</script>

<style type="text/css">
  /* cyrillic */
  @font-face {
    font-family: 'Cuprum';
    font-style: normal;
    font-weight: 400;
    src: local('Cuprum Regular'), local('Cuprum-Regular'), url(./assets/fonts/cuprum_cyrillic.woff2) format('woff2');
    unicode-range: U+0400-045F, U+0490-0491, U+04B0-04B1, U+2116;
  }

  /* latin */
  @font-face {
    font-family: 'Cuprum';
    font-style: normal;
    font-weight: 400;
    src: local('Cuprum Regular'), local('Cuprum-Regular'), url(./assets/fonts/cuprum_latin.woff2) format('woff2');
    unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2212, U+2215;
  }

  .body {
    padding: 0;
    margin: 0;
  }

  .app {
    font-family: 'Cuprum', sans-serif;
    margin: 0;
    padding: 0;
  }
</style>
