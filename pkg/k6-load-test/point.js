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
        point: {
            executor: 'per-vu-iterations',

            vus: USERS,

            // هر user یک درخواست point
            iterations: 1,

            maxDuration: '2m',
        },
    },
};


export default function () {

    const username = getUsername(__VU);
    const password = getPassword();


    // ==========================================
    // STEP 1: LOGIN
    // ==========================================

    const loginPayload = JSON.stringify({
        username: username,
        password: password,
    });


    const loginResponse = http.post(
        `${BASE_URL}/login`,
        loginPayload,
        {
            headers: {
                'Content-Type': 'application/json',
            },

            tags: {
                endpoint: 'login',
                purpose: 'authentication',
            },
        }
    );


    const loginOK = check(loginResponse, {
        'login for point successful': (r) =>
            r.status >= 200 && r.status < 300,
    });


    if (!loginOK) {

        console.error(
            `Point authentication failed | user=${username} | status=${loginResponse.status}`
        );

        return;
    }


    // ==========================================
    // STEP 2: POST /point
    // ==========================================

    const payload = JSON.stringify({
        point: 50,
    });


    const response = http.post(
        `${BASE_URL}/point`,
        payload,
        {
            headers: {
                'Content-Type': 'application/json',
            },

            tags: {
                endpoint: 'point',
            },
        }
    );


    check(response, {
        'point returned 2xx': (r) =>
            r.status >= 200 && r.status < 300,
    });


    if (response.status < 200 || response.status >= 300) {

        console.error(
            `Point failed | user=${username} | status=${response.status} | body=${response.body}`
        );
    }
}