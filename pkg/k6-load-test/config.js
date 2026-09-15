export const BASE_URL = 'http://localhost:8080';

export const USERS = 100;

export function getUsername(vu) {
    return `k6_loadtest_user_${vu}`;
}

export function getPassword() {
    return 'K6_TestPassword_123!';
}