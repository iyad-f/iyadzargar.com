// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package web

import "embed"

// Static is the site's embedded static assets.
//
//go:embed all:static
var Static embed.FS
