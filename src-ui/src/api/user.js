import api from '../api';

export default {
  auth(login, password) {
    return new Promise((resolve, reject) => {
      fetch(`${api.CONTEXT}${api.VERSION}/get-token?uname=${login}&upass=${password}`,
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
};
