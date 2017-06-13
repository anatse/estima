<template>
  <div class="login">
    <form class="login_form">
      <div class="login_row login_row__align_center">
        <h1 class="login_header">Estimator</h1>
      </div>
      <div
        v-if="error"
        class="login_row login_row__align_center login_row__error"
      >
        {{error}}
      </div>
      <div
        v-if="isUserAuth"
        class="login_row login_row__align_center login_row__auth"
      >
        Вы уже авторизованы в системе ({{user.displayName}}).
      </div>
      <div class="login_row login_row__align_center">
        <input
          id="Authentication__input__login"
          class="login_input"
          type="text"
          placeholder="Логин"
          v-model="uname"
        />
      </div>
      <div class="login_row login_row__align_center">
        <input
          id="Authentication__input__password"
          class="login_input"
          type="password"
          placeholder="Пароль"
          v-model="upass"
        />
      </div>
      <div class="login_row login_row__align_center">
        <button
          id="Authentication__button__enter"
          type="button"
          class="login_action"
          :class="{ 'login_action__disabled': loginDisabled }"
          :disabled="loginDisabled"
          @click="authenticate({uname, upass})"
        >
        Вход
        </button>
      </div>
      <div
        v-if="isUserAuth"
        class="login_row login_row__align_center"
      >
        <button
          id="Authentication__button__logout"
          type="button"
          class="login_action login_action__logout"
          @click="logout"
        >
          Выход
        </button>
      </div>
    </form>
  </div>
</template>

<script>
  import constGlobal from '../constGlobal';
  import * as TActions from '../store/action-types';

  export default {
    name: 'LoginPage',
    data() {
      return {
        uname: '',
        upass: '',
      };
    },
    methods: {
      // Авторизация в системе.
      authenticate({ uname, upass }) {
        this.$store.dispatch(TActions.USER_LOGIN, { uname, upass }).then(() => {
          this.$router.push({ name: constGlobal.PAGE_NAME.MAIN_PAGE });
        });
      },
      // Выход из системы.
      logout() {
        this.$store.dispatch(TActions.USER_LOGOUT);
      },
    },
    computed: {
      loginDisabled() {
        return !(this.uname.length !== 0 && this.upass.length !== 0);
      },
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
  .login {
    display: table;
    position: absolute;
    margin: 0;
    padding: 0;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }

  .login_form {
    display: table-cell;
    margin: 0;
    padding: 0;
    vertical-align: middle;
    text-align: center;
  }

  .login_header {
    text-transform: uppercase;
    margin: 0;
    padding: 0;
  }

  .login_row {
    padding: 0 0 8px 0;
  }

  .login_row__align_center {
    text-align: center;
  }

  .login_input {
    width: 120px;
    outline: none;
    padding: 8px;
    margin: 0;
    border: 1px solid darkgray;
  }

  .login_input:focus {
    border: 1px solid orange;
  }

  .login_action {
    cursor: pointer;
    background-color: lightgray;
    width: 140px;
    outline: none;
    padding: 8px;
    margin: 0;
    border: 0;
  }

  .login_action:hover {
    background-color: orange;
  }

  .login_action__logout {
    background-color: red;
    color: white;
  }

  .login_action__disabled {
    cursor: default;
  }

  .login_action__disabled:hover {
    background-color: lightgray;
  }

  .login_row__auth {
    color: green;
  }

  .login_row__error {
    color: red;
  }
</style>
