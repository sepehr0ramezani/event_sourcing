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
        leaderboard: {
            executor: 'per-vu-iterations',

            vus: USERS,

            // هر user یک بار leaderboard
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
    //
    // این Login برای احراز هویت است.
    // JWT از response نمی‌گیریم.
    //
    // k6 cookie دریافتی را خودش در Cookie Jar
    // همین VU ذخیره می‌کند.
    //

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
        'login for leaderboard successful': (r) =>
            r.status >= 200 && r.status < 300,
    });


    if (!loginOK) {

        console.error(
            `Leaderboard authentication failed | user=${username} | status=${loginResponse.status}`
        );

        return;
    }


    // ==========================================
    // STEP 2: LEADERBOARD
    // ==========================================
    //
    // هیچ Authorization دستی نمی‌فرستیم.
    //
    // Cookie که Login ساخته، خودکار ارسال می‌شود.
    //

    const response = http.get(
        `${BASE_URL}/leaderboard`,
        {
            tags: {
                endpoint: 'leaderboard',
            },
        }
    );


    check(response, {
        'leaderboard returned 2xx': (r) =>
            r.status >= 200 && r.status < 300,
    });


    if (response.status < 200 || response.status >= 300) {

        console.error(
            `Leaderboard failed | user=${username} | status=${response.status} | body=${response.body}`
        );
    }
}