/**
 * Получить куки по названию.
 *
 * @param name Название куки.
 */
function get(name) {
  const matches = document.cookie.match(new RegExp(
    `(?:^|; )${name.replace(/([.$?*|{}()[]\/+^])/g, '\\$1')}=([^;]*)`,
  ));
  return matches ? decodeURIComponent(matches[1]) : undefined;
}

function set(name, value, options) {
  const optionsNew = options || {};
  let newValue = value;

  let expires = optionsNew.expires;

  if (typeof expires === 'number' && expires) {
    const d = new Date();
    d.setTime(d.getTime() + (expires * 1000));
    expires = d;
    optionsNew.expires = d;
  }
  if (expires && expires.toUTCString) {
    optionsNew.expires = expires.toUTCString();
  }

  newValue = encodeURIComponent(newValue);

  let updatedCookie = `${name}=${newValue}`;

  Object.keys(optionsNew).forEach((propName) => {
    updatedCookie += `; ${propName}`;
    const propValue = optionsNew[propName];
    if (propValue !== true) {
      updatedCookie += `=${propValue}`;
    }
  });

  document.cookie = updatedCookie;
}

function del(name) {
  set(name, '', {
    expires: -1,
  });
}

export default {
  get,
  del,
  set,
};
