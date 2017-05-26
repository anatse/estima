import Vue from 'vue';
import Router from 'vue-router';

import cookie from '../utils/cookie';
import constGlobal from '../constGlobal';
import pageName from './pageName';

import LoginPage from '../pages/LoginPage';
import ProjectPage from '../pages/ProjectPage';

Vue.use(Router);

const router = new Router({
  routes: [
    {
      path: '/',
      redirect: '/project',
      name: pageName.MAIN_PAGE,
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: pageName.LOGIN_PAGE,
      component: LoginPage,
      meta: { requiresAuth: false },
    },
    {
      path: '/project',
      name: pageName.PROJECT_PAGE,
      component: ProjectPage,
      meta: { requiresAuth: true },
    },
  ],
});

/**
 * Проверка на доступ к страницам у авторизованного пользователя.
 */
router.beforeEach((to, from, next) => {
  const isAuth = cookie.get(constGlobal.cookieAuthName);
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // этот путь требует авторизации, проверяем залогинен ли
    // пользователь, и если нет, перенаправляем на страницу логина
    if (isAuth) {
      next();
    } else {
      next({ name: pageName.LOGIN_PAGE });
    }
  } else {
    next(); // всегда так или иначе нужно вызвать next()!
  }
});

export default router;
