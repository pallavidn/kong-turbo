# Set the base image

FROM alpine:3.3

# Set the file maintainer

MAINTAINER Pallavi Debnath <pallavi.debnath@turbonomic.com>


ADD probe/_output/kongturbo.linux /bin/kongturbo


ENTRYPOINT ["/bin/kongturbo"]
