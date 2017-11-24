# GPU Support

As a Proof-of-Concept how to implement GPU (and InfiniBand) support, `doxy` was extended to allow for injection of payload to the `docker create` call.

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
