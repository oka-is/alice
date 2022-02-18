package engine

import (
	"github.com/gin-gonic/gin"
)

type (
	Engine   = gin.Engine
	ActionFN func(*Context)
)
