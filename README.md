# Tests

Studing diferent approachs to serialize grpc calls, as well as broadcasting counters.

#### channel as Smaphore and channell as working group

       /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: load.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)


     ✗ status is OK
      ↳  99% — ✓ 302988 / ✗ 12

     checks...............: 99.99% ✓ 302988    ✗ 12
     data_received........: 26 MB  352 kB/s
     data_sent............: 28 MB  377 kB/s
     grpc_req_duration....: avg=242.79ms min=305.81µs med=246.38ms max=315.88ms p(90)=265.97ms p(95)=272.59ms
     iteration_duration...: avg=25.59s   min=16.9s    med=25.79s   max=26.1s    p(90)=26.07s   p(95)=26.08s
     iterations...........: 3000   40.676733/s
     vus..................: 821    min=821     max=1000
     vus_max..............: 1000   min=1000    max=1000


running (1m13.8s), 0000/1000 VUs, 3000 complete and 0 interrupted iterations
default ✓ [======================================] 1000 VUs  1m0s

### Global mutext


          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: load.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)


     ✗ status is OK
      ↳  99% — ✓ 304505 / ✗ 10

     checks...............: 99.99% ✓ 304505    ✗ 10
     data_received........: 26 MB  384 kB/s
     data_sent............: 28 MB  412 kB/s
     grpc_req_duration....: avg=228.83ms min=249.39µs med=235.96ms max=697.76ms p(90)=256.22ms p(95)=263.86ms
     iteration_duration...: avg=24.17s   min=1.07s    med=24.51s   max=25.17s   p(90)=25.01s   p(95)=25.06s
     iterations...........: 3015   44.380321/s
     vus..................: 75     min=75      max=1000
     vus_max..............: 1000   min=1000    max=1000


running (1m07.9s), 0000/1000 VUs, 3015 complete and 0 interrupted iterations
default ✓ [======================================] 1000 VUs  1m0s