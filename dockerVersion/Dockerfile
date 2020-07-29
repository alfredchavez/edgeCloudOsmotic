FROM rustlang/rust:nightly
COPY sieve /usr/src/sieve/
WORKDIR /usr/src/sieve
RUN cargo build
ENTRYPOINT ["/usr/src/sieve/target/debug/sieve", "-e"]
CMD ["1000"]

