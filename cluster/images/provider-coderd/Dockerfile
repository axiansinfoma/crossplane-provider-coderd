# The provider is a statically linked Go binary that embeds the coderd Terraform
# provider and drives it in-process, so the image needs nothing beyond the
# binary and CA certificates for talking to the Coder deployment. The distroless
# base ships those and a "nonroot" user with UID 65532.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

ARG TARGETOS
ARG TARGETARCH

COPY "bin/${TARGETOS}_${TARGETARCH}/provider" /usr/local/bin/provider

USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/provider"]
