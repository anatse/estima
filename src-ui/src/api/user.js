import constGlobal from '../constGlobal';

export default {
  auth(login, password) {
    return new Promise((resolve, reject) => {
      fetch(`${constGlobal.REST.CONTEXT}${constGlobal.REST.VERSION}/get-token?uname=${login}&upass=${password}`,
        {
          credentials: 'same-origin',
        })
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
  check(cookieAuth, details) {
    return new Promise((resolve, reject) => {
      if (cookieAuth && details.displayName === null) {
        fetch(`${constGlobal.REST.CONTEXT}${constGlobal.REST.VERSION}/users/current`,
          {
            credentials: 'same-origin',
          })
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
