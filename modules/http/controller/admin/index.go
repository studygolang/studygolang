// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	echo "github.com/labstack/echo/v4"
)

func AdminIndex(ctx echo.Context) error {
	return render(ctx, "index.html", nil)
}
