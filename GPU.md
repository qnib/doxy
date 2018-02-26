# GPU Support

As a Proof-of-Concept how to implement GPU (and InfiniBand) support, `doxy` was extended to allow for injection of payload to the `docker create` call.

## CentOS7 box

### Check for CUDA devices

The BOX has 10 Tesla K80 devices (`/dev/nvidia[0-9]`).

```
[root@odin001 ~]# nvidia-smi -L
GPU 0: Tesla K80 (UUID: GPU-4095713a-1f9b-791d-841d-8b35143127d4)
GPU 1: Tesla K80 (UUID: GPU-ab541226-7c4f-ab59-3927-1535f68a3a8f)
GPU 2: Tesla K80 (UUID: GPU-8310202d-1d32-bac5-cc36-2add9e21d9d6)
GPU 3: Tesla K80 (UUID: GPU-cb3d675d-ba3b-5cdb-2331-9534bfd20679)
GPU 4: Tesla K80 (UUID: GPU-8f1511d6-5326-e718-c682-cd2377bbf7cf)
GPU 5: Tesla K80 (UUID: GPU-997c3b02-765c-7cde-daff-f19dadeb6894)
GPU 6: Tesla K80 (UUID: GPU-bb1a7162-859a-4e9c-eeac-150bd35ff767)
GPU 7: Tesla K80 (UUID: GPU-9698fbda-39fd-5f1f-1691-626c7e780f36)
GPU 8: Tesla K80 (UUID: GPU-243c49d3-5dc3-26cc-2f6e-42e2baf6ac93)
GPU 9: Tesla K80 (UUID: GPU-0fc004de-3f1a-030a-3931-33043b895514)
```

### Starting the proxy
```
[root@odin001 ~]# docker run -v /var/run:/var/run/ -ti --rm qnib/doxy:gpu doxy --pattern-key=hpc --debug --proxy-socket=/var/run/hpc.sock --gpu
> execute CMD 'doxy --pattern-key=hpc --debug --proxy-socket=/var/run/hpc.sock --gpu'
2018/02/25 00:07:17 [II] Start Version: 0.2.4
2018/02/25 00:07:17 Error reading patterns file 'open /etc/doxy.pattern: no such file or directory', using hpc patterns
2018/02/25 00:07:17 [doxy] Listening on /var/run/hpc.sock
2018/02/25 00:07:17 Serving proxy on '/var/run/hpc.sock'
[negroni] 2018-02-25T00:07:52Z | 200 |   843.548µs | docker | GET /_ping
Add GPU stuff
New device: /dev/nvidia0:/dev/nvidia0:rwm
New device: /dev/nvidiactl:/dev/nvidiactl:rwm
```

### Running a CUDA image:

Since passing through the log-API calls gave me some headache, the container is created using the proxy, but started using the proper docker API.

```
$ docker start -ai $(docker -H unix:///var/run/hpc.sock create \
                            -ti -e SKIP_ENTRYPOINTS=true \
                            qnib/cplain-cuda nvidia-smi -L)
[II] qnib/init-plain script v0.4.28
> execute CMD 'nvidia-smi -L'
GPU 0: Tesla K80 (UUID: GPU-4095713a-1f9b-791d-841d-8b35143127d4)
```

## AWS p2.xlarge
On an `p2.xlarge` instance with CUDA support, `doxy:gpu` was started like that.

```bash
$ docker run -v /var/run:/var/run/ -ti --rm qnib/doxy:gpu doxy --pattern-key=hpc --debug --proxy-socket=/var/run/hpc.sock --gpu
```

The pattern `hpc` allows for read/write endpoints to be used and `--gpu` injects bind-mounts, device-mappings and environment variables to make it work.

```bash
root@ip-172-31-28-27:~# docker -H unix:///var/run/hpc.sock create nvidia/cuda nvidia-smi -L
6b708d7cda36b9c37e325893108839f3b02f172e40ab97182fa77d770cc219fb
root@ip-172-31-28-27:~# docker start -a 6b708d7cda36b9c37e325893108839f3b02f172e40ab97182fa77d770cc219fb
GPU 0: Tesla K80 (UUID: GPU-234d5537-ea27-68cd-7337-ee21b2f34bf1)
root@ip-172-31-28-27:~#
```

When `start`ing the container the default socket is used; as I had issue to pass `stdin`/`stdout` through.
The command hangs with the `doxy`-socket.
I suspect issues with `stdin`/`stdout` passthrough when proxying the two unix-sockets. :/
