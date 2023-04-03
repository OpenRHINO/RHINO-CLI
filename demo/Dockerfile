FROM openrhino/mpibuilder_base:v0.1.0 as builder

ARG func_name ${func_name}
ARG file ${file}
ARG make_args ${make_args}
ENV FUNC_NAME=${func_name}

COPY src/ /app/src
COPY ldd.sh /app/
RUN cd $(dirname ${file}) && make -B -f $(basename ${file}) ${make_args}

RUN sh ldd.sh

FROM openrhino/mpirun_base:v0.1.0

ARG func_name ${func_name}
COPY --from=builder /app/${func_name}  /app/${func_name}
COPY --from=builder /shared_lib /usr/local/lib

CMD ["/bin/ash"]