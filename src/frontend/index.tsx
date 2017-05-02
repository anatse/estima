// function login () {
//     var url = '/get-token?uname=' + encodeURIComponent(document.getEl
// ementsByName("user")[0].value)
// + "&upass=" + encodeURIComponent(document.getElementsByName("pwd")[0].value)
//     fetch(url, {
//         method: "get",
//         headers: {
//             'Accept': 'application/json',
//             'Cache': 'no-cache'
//         },
//         credentials: "same-origin"
//     }).then(function(response) {
//         response.text().then (function (text) {
//             console.log (text)
//         })
//     })
// }
//
// function products () {
//     var url = '/products'
//     fetch(url, {
//         method: "get",
//         headers: {
//             'Accept': 'application/json'
//         },
//         credentials: "same-origin"
//     }).then(function(response) {
//         response.text().then (function (text) {
//             console.log (text)
//         })
//     })
// }
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {AppRouters} from './components/AppRouters';

import './styles/fonts.css';
import './styles/app.css';

function contentLoaded() {
    ReactDOM.render(<AppRouters />, document.getElementById('app'));
}

document.addEventListener('DOMContentLoaded', contentLoaded);
