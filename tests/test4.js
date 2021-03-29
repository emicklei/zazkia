import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 20 },
    { duration:  '2m', target: 200 },
    { duration: '10s', target: 0 },
  ],
};
/*
export let options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m30s', target: 10 },
    { duration: '20s', target: 0 },
  ],
};
*/
export default function () {
  let res = http.get('http://__IP__:9090');
  check(res, { 
    'status was 200': (r) => r.status === 200,
    'body size is 100 bytes': (r) => r.body.length == 100,
  });
  sleep(1);
}
