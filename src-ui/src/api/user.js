import constGlobal from '../constGlobal';

const methodGlobal = {
  GET: 'GET',
  POST: 'POST',
};

const urlAddress = {
  login: `${constGlobal.REST.CONTEXT}${constGlobal.REST.VERSION}/login`,
  curent: `${constGlobal.REST.CONTEXT}${constGlobal.REST.VERSION}/users/current`,
};

/**
 * Подготавливаем конфигурацию для запроса.
 *
 * @param method
 * @param body
 * @param options
 * @return {{method: string, body: {}, credentials: string}}
 */
function configureOptions(method = methodGlobal.GET, body = {}, options = {}) {
  switch (method) {
    case methodGlobal.GET:
      break;
    case methodGlobal.POST:
      body = JSON.stringify(body);
      options = {
        headers: {
          Accept: 'application/json, text/plain, */*',
          'Content-Type': 'application/json',
        },
        ...options,
      };
      break;
    default:
  }

  return {
    method,
    body,
    credentials: 'same-origin',
    ...options,
  };
}

export default {
  /**
   * Авторизация в системе.
   *
   * @param uname
   * @param upass
   * @return {Promise}
   */
  login(uname, upass) {
    return new Promise((resolve, reject) => {
      fetch(urlAddress.login, configureOptions(methodGlobal.POST, { uname, upass }))
        .then(response => response.json())
        .then((response) => {
          if (response.success) {
            resolve(response.body);
          } else {
            reject(response.error);
          }
        })
        .catch(reject);
    });
  },

  /**
   * Получаем данные о текущем пользователе.
   *
   * @param cookieAuth {string} Данные в cookie.
   * @param details {Object} Данные о пользователе.
   * @return {Promise}
   */
  current(cookieAuth, details) {
    return new Promise((resolve, reject) => {
      if (cookieAuth && details.displayName === null) {
        fetch(urlAddress.curent, configureOptions())
          .then(response => response.json())
          .then((response) => {
            if (response.success) {
              resolve(response.body);
            } else {
              reject(response.error);
            }
          })
          .catch(reject);
      } else {
        resolve(details);
      }
    });
  },
};
