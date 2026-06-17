# SPDX-FileCopyrightText: 2026 Iyad
# SPDX-License-Identifier: Apache-2.0

FROM caddy:2.11.4-builder AS build
RUN xcaddy build v2.11.4 --with github.com/mholt/caddy-ratelimit@v0.1.0

FROM caddy:2.11.4
COPY --from=build /usr/bin/caddy /usr/bin/caddy
