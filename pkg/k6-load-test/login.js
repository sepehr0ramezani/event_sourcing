import http from 'k6/http';
import { check } from 'k6';

import {
    BASE_URL,
    USERS,
    getUsername,
    getPassword,
} from './config.js';


export const options = {
    scenarios: {
        login: {
            executor: 'per-vu-iterations',

            vus: USERS,

            // هر کاربر دقیقاً یک Login
            iterations: 1,

            maxDuration: '2m',
        },
    },
};


export default function () {

    const username = getUsername(__VU);
    const password = getPassword();


    const payload = JSON.stringify({
        username: username,
        password: password,
    });


    const response = http.post(
        `${BASE_URL}/login`,
        payload,
        {
            headers: {
                'Content-Type': 'application/json',
            },

            tags: {
                endpoint: 'login',
            },
        }
    );


    check(response, {
        'login returned 2xx': (r) =>
            r.status >= 200 && r.status < 300,
    });


    if (response.status < 200 || response.status >= 300) {
        console.error(
            `Login failed | user=${username} | status=${response.status} | body=${response.body}`
        );
    }
}