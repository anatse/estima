<template>
  <div v-if="isUserAuth" class="navbar">
    <div class="navbar_item navbar_item__left">
      Система ведения бэклога
    </div>
    <div class="navbar_item">
      {{user.displayName}}
    </div>
    <div class="navbar_item">
      <button
        type="button"
        class="navbar_action"
        @click="logout"
      >
        Выход
      </button>
    </div>
  </div>
</template>

<script>
  import * as TActions from '../store/action-types';
  import constGlobal from '../constGlobal';

  export default {
    name: 'Navbar',
    methods: {
      // Выход из системы.
      logout() {
        this.$store.dispatch(TActions.USER_LOGOUT).then(() => {
          this.$router.push({ name: constGlobal.PAGE_NAME.LOGIN_PAGE });
        });
      },
    },
    computed: {
      isUserAuth() {
        return this.$store.getters.isUserAuth;
      },
      error() {
        return this.$store.getters.error;
      },
      user() {
        return this.$store.getters.getUser;
      },
    },
  };
</script>

<style>
.navbar {
    display: block;
    padding: 8px;
    text-align: right;
}

.navbar_item {
    display: inline-block;
    padding: 0 8px;
    border-right: 2px solid gray
}

.navbar_item__left {
  float: left;
}

.navbar_item:last-child {
    border-right: none;
    padding: 0 4px 0 8px;
}

.navbar_item:first-child {
    border-right: none;
    padding: 0 8px 0 4px;
}

.navbar_action {
    background: red;
    color: white;
    border: none;
}

.navbar_action:hover {
    background: orange;
    cursor: pointer;
}
</style>
