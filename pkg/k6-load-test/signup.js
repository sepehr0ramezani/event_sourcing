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
        signup: {
            executor: 'per-vu-iterations',

            // دقیقاً 100 کاربر
            vus: USERS,

            // هر کاربر فقط یک بار signup می‌کند
            iterations: 1,

            maxDuration: '2m',
        },
    },
};


export default function () {

    // __VU از 1 تا 100 است
    const username = getUsername(__VU + 400);
    const password = getPassword();

    const payload = JSON.stringify({
        username: username,
        password: password,

        // point باید int باشد
        point: 100,
    });


    const response = http.post(
        `${BASE_URL}/signup`,
        payload,
        {
            headers: {
                'Content-Type': 'application/json',
            },

            tags: {
                endpoint: 'signup',
            },
        }
    );


    check(response, {
        'signup returned 2xx': (r) =>
            r.status >= 200 && r.status < 300,
    });


    if (response.status < 200 || response.status >= 300) {
        console.error(
            `Signup failed | user=${username} | status=${response.status} | body=${response.body}`
        );
    }
}