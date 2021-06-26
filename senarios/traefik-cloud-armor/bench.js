import http from 'k6/http';
import {Counter, Trend} from 'k6/metrics';

const ip = {
    cloudArmor: 'IP',
    customOrigin: 'IP',
    traefik: 'IP',
    serviceLB: 'IP'
}

const cloudArmorDur = new Trend('cloudArmorDur');
const cloudArmorDurCustom = new Trend('cloudArmorDurCustom');
const traefikOnlyDur = new Trend('traefikOnlyDur');
const serviceLBDur = new Trend('serviceLBDur');

const cloudArmorCount = new Counter('cloudArmorCount');
const cloudArmorCountCustom = new Counter('cloudArmorCountCustom');
const traefikOnlyCount = new Counter('traefikOnlyCount');
const serviceLBCount = new Counter('serviceLBCount');

export let options = {
    discardResponseBodies: true,
    scenarios: {
        cloudArmor: {
            exec: 'cloudArmor',
            executor: 'constant-vus',
            vus: 50,
            duration: '1m',
        },
        cloudArmorCustom: {
            exec: 'cloudArmorCustom',
            executor: 'constant-vus',
            vus: 50,
            duration: '1m',
        },
        traefikOnly: {
            exec: 'traefikOnly',
            executor: 'constant-vus',
            vus: 50,
            duration: '1m',
        },
        serviceLB: {
            exec: 'serviceLB',
            executor: 'constant-vus',
            vus: 50,
            duration: '1m',
        },
    },
};

export function cloudArmor() {
    const r = http.get(ip.cloudArmor);
    cloudArmorDur.add(r.timings.duration);
    cloudArmorCount.add(1);
}

export function cloudArmorCustom() {
    const r = http.get(ip.customOrigin);
    cloudArmorDurCustom.add(r.timings.duration);
    cloudArmorCountCustom.add(1);
}

export function traefikOnly() {
    const r = http.get(ip.traefik);
    traefikOnlyDur.add(r.timings.duration);
    traefikOnlyCount.add(1);
}


export function serviceLB() {
    const r = http.get(ip.serviceLB);
    serviceLBDur.add(r.timings.duration);
    serviceLBCount.add(1);
}
