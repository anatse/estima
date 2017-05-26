import userApi from '../../api/user';
import cookie from '../../utils/cookie';
import constGlobal from '../../constGlobal';

import * as types from '../mutation-types';
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
  authenticate({ commit }, { login, password }) {
    commit(types.LOGIN_REQUEST);

    return userApi.auth(
      login,
      password,
    ).then(
      user => commit(types.LOGIN_SUCCESS, user),
      error => commit(types.LOGIN_FAILURE, error),
    );
  },
  logout({ commit }) {
    cookie.del(constGlobal.cookieAuthName);
    commit(types.LOGOUT);
  },
};

/* eslint no-param-reassign: "off" */
// mutations
const mutations = {

  [types.LOGIN_REQUEST](state) {
    state.status = null;
  },

  [types.LOGIN_SUCCESS](state, user) {
    state.details = user;
    state.error = null;
    state.status = RESPONSE_STATUSES.SUCCESSFUL;
  },

  [types.LOGOUT](state) {
    state.details = getStateInit().details;
    state.status = null;
  },

  [types.LOGIN_FAILURE](state, error) {
    state.details = getStateInit().details;
    state.error = error;
    state.status = RESPONSE_STATUSES.FAILURE;
  },
};

export default {
  state: getStateInit(),
  getters,
  actions,
  mutations,
};
