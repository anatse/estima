import userApi from '../../api/user';
import cookie from '../../utils/cookie';
import constGlobal from '../../constGlobal';

import * as TMutations from '../mutation-types';
import * as TActions from '../action-types';
import { RESPONSE_STATUSES } from '../const';

// initial state
function getStateInit() {
  return {
    status: null,
    error: null,
    details: {
      displayName: null,
      email: null,
      name: null,
    },
  };
}

// getters
const getters = {
  isUserAuth: state => state.details.email !== null,
  getUser: state => state.details,
  error: state => state.error,
};


// actions
const actions = {
  [TActions.USER_LOGIN]({ commit }, { uname, upass }) {
    commit(TMutations.USER_LOGIN_REQUEST);

    return userApi.login(
      uname,
      upass,
    ).then(
      user => commit(TMutations.USER_LOGIN_SUCCESS, user),
      error => commit(TMutations.USER_LOGIN_FAILURE, error),
    );
  },
  [TActions.USER_LOGOUT]({ commit }) {
    cookie.del(constGlobal.cookieAuthName);
    commit(TMutations.USER_LOGOUT);
  },
  [TActions.USER_CURRENT]({ commit, state }) {
    return userApi.current(cookie.get(constGlobal.cookieAuthName), state.details).then(
      user => commit(TMutations.USER_LOGIN_SUCCESS, user),
      error => commit(TMutations.USER_LOGIN_FAILURE, error),
    );
  },
};

/* eslint no-param-reassign: "off" */
// mutations
const mutations = {

  [TMutations.USER_LOGIN_REQUEST](state) {
    state.status = null;
  },

  [TMutations.USER_LOGIN_SUCCESS](state, user) {
    state.details = user;
    state.error = null;
    state.status = RESPONSE_STATUSES.SUCCESSFUL;
  },

  [TMutations.USER_LOGIN_FAILURE](state, error) {
    state.details = getStateInit().details;
    state.error = error;
    state.status = RESPONSE_STATUSES.FAILURE;
  },

  [TMutations.USER_LOGOUT](state) {
    state.details = getStateInit().details;
    state.status = null;
  },
};

export default {
  state: getStateInit(),
  getters,
  actions,
  mutations,
};
